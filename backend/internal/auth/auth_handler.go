package auth

import (
	"log"
	"net/http"
	"time"

	"github.com/aniruddha-jafa/go-auth-v1/internal/apperrors"
	"github.com/aniruddha-jafa/go-auth-v1/internal/users"
	"github.com/gofiber/fiber/v2"
)

const REFRESH_TOKEN_COOKIE_NAME = "refreshToken"

type AuthHandler interface {
	Login(ctx *fiber.Ctx) error
	SignUp(ctx *fiber.Ctx) error
	RefreshToken(ctx *fiber.Ctx) error
	Logout(ctx *fiber.Ctx) error
}

type AuthHandlerImpl struct {
	AuthService AuthService
}

func (h *AuthHandlerImpl) Login(c *fiber.Ctx) error {
	ctx := c.UserContext()
	loginRequest := new(LoginRequest)
	if err := c.BodyParser(loginRequest); err != nil {
		return apperrors.NewHttpError(http.StatusBadRequest, "unable to parse to login request")
	}
	log.Printf("Login request: %s", loginRequest.String())
	err := loginRequest.Validate()
	if err != nil {
		return apperrors.NewHttpError(http.StatusBadRequest, err.Error())
	}
	loginRes, err := h.AuthService.Login(&ctx, *loginRequest)
	if err != nil {
		log.Printf("Login error: %s", err)
		return err
	}
	refreshTokenCookie := h.getRefreshTokenCookie(loginRes.RefreshToken, loginRes.RefreshTokenExpiresAt)
	c.Cookie(refreshTokenCookie)

	log.Printf("Login successful for user: %v", loginRes)
	c.Status(http.StatusOK).JSON(loginRes)
	return nil
}

// Returns a fiber.Cookie for the refresh token
func (h *AuthHandlerImpl) getRefreshTokenCookie(refreshTokenValue string, expiresAt time.Time) *fiber.Cookie {
	return &fiber.Cookie{
		Name:     REFRESH_TOKEN_COOKIE_NAME,
		Value:    refreshTokenValue,
		Expires:  expiresAt,
		HTTPOnly: true,
		Secure:   true,
		// Allow cross-site requests from the frontend
		SameSite: fiber.CookieSameSiteNoneMode,
	}
}

func (h *AuthHandlerImpl) SignUp(c *fiber.Ctx) error {
	log.Println("SignUp request received")
	ctx := c.UserContext()
	signupRequest := new(users.UserCreationRequest)
	if err := c.BodyParser(signupRequest); err != nil {
		return apperrors.NewHttpError(http.StatusBadRequest, "unable to parse to signup request")
	}
	log.Printf("SignUp request: %s", signupRequest)
	err := signupRequest.Validate()
	if err != nil {
		return apperrors.NewHttpError(http.StatusBadRequest, err.Error())
	}
	userRes, err := h.AuthService.SignUp(&ctx, *signupRequest)
	if err != nil {
		log.Printf("SignUp error: %s", err)
		return err
	}
	log.Printf("SignUp response: %s", userRes)
	c.Status(http.StatusCreated).JSON(userRes)
	return nil
}

// Uses the refresh token to generate a new JWT,
func (h *AuthHandlerImpl) RefreshToken(c *fiber.Ctx) error {
	ctx := c.UserContext()
	refreshToken := c.Cookies(REFRESH_TOKEN_COOKIE_NAME)
	if refreshToken == "" {
		return apperrors.ErrRefreshTokenNotFound
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
	refreshToken := c.Cookies(REFRESH_TOKEN_COOKIE_NAME)
	if refreshToken == "" {
		return apperrors.ErrRefreshTokenNotFound
	}
	err := h.AuthService.Logout(&ctx, refreshToken)
	if err != nil {
		return err
	}
	c.ClearCookie(REFRESH_TOKEN_COOKIE_NAME)
	c.Status(http.StatusNoContent)
	return nil
}
