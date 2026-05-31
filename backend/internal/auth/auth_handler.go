package auth

import (
	"log/slog"
	"net/http"
	"time"

	"github.com/aniruddha-jafa/go-auth-v1/internal/apperrors"
	"github.com/aniruddha-jafa/go-auth-v1/internal/config"
	"github.com/aniruddha-jafa/go-auth-v1/internal/constants"
	"github.com/aniruddha-jafa/go-auth-v1/internal/logger"
	"github.com/aniruddha-jafa/go-auth-v1/internal/users"
	"github.com/gofiber/fiber/v2"
)

type AuthHandler interface {
	Login(ctx *fiber.Ctx) error
	SignUp(ctx *fiber.Ctx) error
	RefreshToken(ctx *fiber.Ctx) error
	GetOrResetCSRFToken(ctx *fiber.Ctx) error
	Logout(ctx *fiber.Ctx) error
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

func (h *AuthHandlerImpl) Login(c *fiber.Ctx) error {
	ctx := c.UserContext()
	log := logger.WithContext(h.baseLogger, ctx)

	loginRequest := new(LoginRequest)
	if err := c.BodyParser(loginRequest); err != nil {
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

	refreshTokenCookie := h.getRefreshTokenCookie(loginRes.RefreshToken, loginRes.RefreshTokenExpiresAt)
	csrfTokenCookie := h.getCSRFTokenCookie(loginRes.CSRFToken)
	c.Cookie(refreshTokenCookie)
	c.Cookie(csrfTokenCookie)

	log.Info("Login successful for user", "userId", loginRes.ID, "email", loginRes.Email)
	c.Status(http.StatusOK).JSON(loginRes)
	return nil
}

// Returns a fiber.Cookie for the refresh token
func (h *AuthHandlerImpl) getRefreshTokenCookie(refreshTokenValue string, expiresAt time.Time) *fiber.Cookie {
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

func (h *AuthHandlerImpl) getCSRFTokenCookie(csrfTokenValue string) *fiber.Cookie {
	appConfig := config.InitAppConfig()
	expiresAt := time.Now().Add(appConfig.CSRFTokenValidity)
	return &fiber.Cookie{
		Name:     constants.CSRF_TOKEN_COOKIE_NAME,
		Value:    csrfTokenValue,
		Expires:  expiresAt,
		HTTPOnly: false,
		Secure:   true,
		Path:     constants.API + constants.V1 + constants.AUTH,
		SameSite: fiber.CookieSameSiteLaxMode,
	}
}

// Gets the CSRF token from the cookie if it exists, otherwise generates a new one and sets it in the cookie.
func (h *AuthHandlerImpl) GetOrResetCSRFToken(c *fiber.Ctx) error {
	log := logger.WithContext(h.baseLogger, c.UserContext())
	if c.Cookies(constants.CSRF_TOKEN_COOKIE_NAME) != "" {
		log.Info("CSRF token found in cookie")
		return c.Status(http.StatusOK).JSON(NewCsrfTokenResponse(c.Cookies(constants.CSRF_TOKEN_COOKIE_NAME)))
	}
	csrfTokenValue, err := generateCSRFToken()
	if err != nil {
		log.Error("Error generating CSRF token", "error", err)
		return apperrors.ErrInternalServerError
	}
	c.Cookie(h.getCSRFTokenCookie(csrfTokenValue))
	c.Status(http.StatusOK).JSON(NewCsrfTokenResponse(csrfTokenValue))
	log.Info("New CSRF token generated successfully")
	return nil
}

func (h *AuthHandlerImpl) SignUp(c *fiber.Ctx) error {
	ctx := c.UserContext()
	log := logger.WithContext(h.baseLogger, ctx)

	log.Info("SignUp request received")
	signupRequest := new(users.UserCreationRequest)
	if err := c.BodyParser(signupRequest); err != nil {
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
func (h *AuthHandlerImpl) RefreshToken(c *fiber.Ctx) error {
	ctx := c.UserContext()
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
func (h *AuthHandlerImpl) Logout(c *fiber.Ctx) error {
	ctx := c.UserContext()
	refreshToken := c.Cookies(constants.REFRESH_TOKEN_COOKIE_NAME)
	if refreshToken == "" {
		return apperrors.ErrRefreshTokenCookieNotFound
	}
	err := h.AuthService.Logout(&ctx, refreshToken)
	if err != nil {
		return err
	}
	c.ClearCookie(constants.REFRESH_TOKEN_COOKIE_NAME)
	c.ClearCookie(constants.CSRF_TOKEN_COOKIE_NAME)
	c.Status(http.StatusNoContent)
	return nil
}
