package auth

import (
	"log"
	"net/http"

	"github.com/aniruddha-jafa/go-auth-v1/internal/apperrors"
	"github.com/aniruddha-jafa/go-auth-v1/internal/refresh_tokens"
	"github.com/aniruddha-jafa/go-auth-v1/internal/users"
	"github.com/gofiber/fiber/v2"
)

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
	log.Println("Login request received")
	ctx := c.UserContext()
	loginRequest := new(LoginRequest)
	if err := c.BodyParser(loginRequest); err != nil {
		return apperrors.NewHttpError(http.StatusBadRequest, "unable to parse to login request")
	}
	log.Printf("Login request: %s", loginRequest)
	err := loginRequest.Validate()
	if err != nil {
		return apperrors.NewHttpError(http.StatusBadRequest, err.Error())
	}
	userRes, err := h.AuthService.Login(&ctx, *loginRequest)
	if err != nil {
		log.Printf("Login error: %s", err)
		return err
	}
	log.Printf("Login successful for user: %v", userRes)
	c.Status(http.StatusOK).JSON(userRes)
	return nil
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

// Use the refresh token to generate a new JWT,
// if the refresh token is still valid.
func (h *AuthHandlerImpl) RefreshToken(c *fiber.Ctx) error {
	ctx := c.UserContext()
	// Get the refresh token from the Authorization header
	bearerToken, err := GetBearerToken(c.GetReqHeaders())
	if err != nil {
		return err
	}
	refreshToken, err := h.AuthService.RefreshToken(&ctx, bearerToken)
	if err != nil {
		return err
	}
	c.Status(http.StatusOK).JSON(refresh_tokens.RefreshTokenResponse{Token: refreshToken.Token})
	return nil
}

// Logs out the user by revoking the refresh token.
func (h *AuthHandlerImpl) Logout(c *fiber.Ctx) error {
	ctx := c.UserContext()
	// Get the refresh token from the Authorization header
	bearerToken, err := GetBearerToken(c.GetReqHeaders())
	if err != nil {
		return err
	}
	err = h.AuthService.Logout(&ctx, bearerToken)
	if err != nil {
		return err
	}
	c.Status(http.StatusNoContent)
	return nil
}
