package auth

import (
	"log/slog"
	"net/http"
	"time"

	"github.com/aniruddha-jafa/go-auth-v1/internal/apperrors"
	"github.com/aniruddha-jafa/go-auth-v1/internal/constants"
	"github.com/aniruddha-jafa/go-auth-v1/internal/logger"
	"github.com/aniruddha-jafa/go-auth-v1/internal/users"
	"github.com/aniruddha-jafa/go-auth-v1/pkg/util"
	"github.com/gofiber/fiber/v3"
	"github.com/valyala/fasthttp"
)

type AuthHandler interface {
	Login(ctx fiber.Ctx) error
	SignUp(ctx fiber.Ctx) error
	RefreshToken(ctx fiber.Ctx) error
	GetOrResetCSRFToken(ctx fiber.Ctx) error
	Logout(ctx fiber.Ctx) error
}

type AuthHandlerImpl struct {
	AuthService AuthService
	baseLogger  *slog.Logger
}

func NewAuthHandlerImpl(authService AuthService) AuthHandler {
	return &AuthHandlerImpl{
		AuthService: authService,
		baseLogger:  slog.Default().With(logger.LoggerNameKey, "AuthHandler"),
	}
}

func (h *AuthHandlerImpl) Login(c fiber.Ctx) error {
	ctx := c.Context()
	log := logger.WithContext(h.baseLogger, ctx)

	loginRequest := new(LoginRequest)
	if err := c.Bind().Body(loginRequest); err != nil {
		return apperrors.NewHttpError(http.StatusBadRequest, "unable to parse to login request")
	}

	log.Info("Login request", "loginRequest", loginRequest.String())
	err := loginRequest.Validate()
	if err != nil {
		return apperrors.NewHttpError(http.StatusBadRequest, err.Error())
	}

	loginRes, err := h.AuthService.Login(&ctx, *loginRequest)
	if err != nil {
		log.Error("Login error", "error", err)
		return err
	}

	refreshTokenCookie := h.createRefreshTokenCookie(loginRes.RefreshToken, loginRes.RefreshTokenExpiresAt)

	// Reuse refresh token expiry for CSRF token
	csrfTokenCookie := h.createCSRFTokenCookie(loginRes.CSRFToken, loginRes.RefreshTokenExpiresAt)

	c.Cookie(refreshTokenCookie)
	c.Cookie(csrfTokenCookie)

	log.Info("Login successful for user", "userId", loginRes.ID, "email", loginRes.Email)
	c.Status(http.StatusOK).JSON(loginRes)
	return nil
}

// Returns a fiber.Cookie for the refresh token
func (h *AuthHandlerImpl) createRefreshTokenCookie(refreshTokenValue string, expiresAt time.Time) *fiber.Cookie {
	return &fiber.Cookie{
		Name:     constants.REFRESH_TOKEN_COOKIE_NAME,
		Value:    refreshTokenValue,
		Expires:  expiresAt,
		HTTPOnly: true,
		Secure:   true,
		// Allow on all auth routes
		Path: constants.API + constants.V1 + constants.AUTH,
		// Allow cross-site requests from the frontend
		SameSite: fiber.CookieSameSiteNoneMode,
	}
}

func (h *AuthHandlerImpl) createCSRFTokenCookie(csrfTokenValue string, expiresAt time.Time) *fiber.Cookie {
	return &fiber.Cookie{
		Name:     constants.CSRF_TOKEN_COOKIE_NAME,
		Value:    csrfTokenValue,
		Expires:  expiresAt,
		HTTPOnly: true,
		Secure:   true,
		Path:     constants.API + constants.V1 + constants.AUTH,
		SameSite: fiber.CookieSameSiteNoneMode,
	}
}

// Gets the CSRF token from the cookie if it exists, otherwise generates a new one and sets it in the cookie.
func (h *AuthHandlerImpl) GetOrResetCSRFToken(c fiber.Ctx) error {
	ctx := c.Context()
	log := logger.WithContext(h.baseLogger, ctx)

	refreshTokenValue := c.Cookies(constants.REFRESH_TOKEN_COOKIE_NAME)
	if refreshTokenValue == "" {
		log.Error("Refresh token not found in cookie")
		return apperrors.ErrUnauthorized
	}
	now := util.Now()
	refreshToken, err := h.AuthService.ValidateRefreshToken(&ctx, refreshTokenValue, now)
	if err != nil {
		log.Error("Refresh token is not valid", "error", err)
		return apperrors.ErrUnauthorized
	}

	if c.Cookies(constants.CSRF_TOKEN_COOKIE_NAME) != "" {
		log.Info("CSRF token found in cookie")
		return c.Status(http.StatusOK).JSON(NewCsrfTokenResponse(c.Cookies(constants.CSRF_TOKEN_COOKIE_NAME)))
	}
	csrfTokenValue, err := generateCSRFToken()
	if err != nil {
		log.Error("Error generating CSRF token", "error", err)
		return apperrors.ErrInternalServerError
	}

	c.Cookie(h.createCSRFTokenCookie(csrfTokenValue, refreshToken.ExpiresAt))

	c.Status(http.StatusOK).JSON(NewCsrfTokenResponse(csrfTokenValue))
	log.Info("New CSRF token generated successfully")
	return nil
}

func (h *AuthHandlerImpl) SignUp(c fiber.Ctx) error {
	ctx := c.Context()
	log := logger.WithContext(h.baseLogger, ctx)

	log.Info("SignUp request received")
	signupRequest := new(users.UserCreationRequest)
	if err := c.Bind().Body(signupRequest); err != nil {
		return apperrors.NewHttpError(http.StatusBadRequest, "unable to parse to signup request")
	}

	log.Info("SignUp request", "signupRequest", signupRequest.String())
	err := signupRequest.Validate()
	if err != nil {
		return apperrors.NewHttpError(http.StatusBadRequest, err.Error())
	}

	userRes, err := h.AuthService.SignUp(&ctx, *signupRequest)
	if err != nil {
		log.Error("SignUp error", "error", err)
		return err
	}
	log.Info("SignUp successful", "userId", userRes.ID, "email", userRes.Email)
	c.Status(http.StatusCreated).JSON(userRes)
	return nil
}

// Uses the refresh token to generate a new JWT,
func (h *AuthHandlerImpl) RefreshToken(c fiber.Ctx) error {
	ctx := c.Context()
	refreshToken := c.Cookies(constants.REFRESH_TOKEN_COOKIE_NAME)
	if refreshToken == "" {
		return apperrors.ErrRefreshTokenCookieNotFound
	}
	jwtToken, err := h.AuthService.RefreshToken(&ctx, refreshToken)
	if err != nil {
		return err
	}
	c.Status(http.StatusOK).JSON(jwtToken)
	return nil
}

// Logs out the user by revoking the refresh token.
func (h *AuthHandlerImpl) Logout(c fiber.Ctx) error {
	ctx := c.Context()
	log := logger.WithContext(h.baseLogger, ctx)

	refreshToken := c.Cookies(constants.REFRESH_TOKEN_COOKIE_NAME)
	csrfCookie := c.Cookies(constants.CSRF_TOKEN_COOKIE_NAME)

	if refreshToken == "" {
		return apperrors.ErrRefreshTokenCookieNotFound
	}
	if csrfCookie == "" {
		log.Error("CSRF token not found in cookie")
		return apperrors.ErrUnauthorized
	}

	err := h.AuthService.Logout(&ctx, refreshToken)
	if err != nil {
		return err
	}

	// Recreate the cookies so we can get the same params to clear them
	h.clearCookie(c, h.createRefreshTokenCookie(refreshToken, time.Time{}))
	h.clearCookie(c, h.createCSRFTokenCookie(csrfCookie, time.Time{}))

	c.Status(http.StatusNoContent)
	return nil
}

func (h *AuthHandlerImpl) clearCookie(c fiber.Ctx, cookie *fiber.Cookie) {
	// Set the cookie to expire immediately
	cookie.Expires = fasthttp.CookieExpireDelete
	c.Cookie(cookie)
}

// fiber:context-methods migrated
