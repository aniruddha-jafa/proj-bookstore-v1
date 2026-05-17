package middleware

import (
	"context"
	"log"

	"github.com/aniruddha-jafa/go-auth-v1/internal/apperrors"
	"github.com/aniruddha-jafa/go-auth-v1/internal/auth"
	"github.com/aniruddha-jafa/go-auth-v1/internal/config"
	"github.com/aniruddha-jafa/go-auth-v1/pkg/security"
	"github.com/gofiber/fiber/v2"
)

const UserIdKey ctxKey = "userId"

func RequireAuth(c *fiber.Ctx) error {
	headers := c.GetReqHeaders()
	bearerToken, err := auth.GetBearerToken(headers)
	if err != nil {
		return err
	}
	userId, err := security.ValidateJwt(bearerToken, config.InitAppConfig().JwtSecret)
	if err != nil {
		log.Printf("Error validating JWT: %v", err)
		return apperrors.ErrInvalidCredentials
	}

	ctx := context.WithValue(c.UserContext(), UserIdKey, userId.String())
	c.SetUserContext(ctx)

	c.Set("userId", userId.String())
	return c.Next()
}
