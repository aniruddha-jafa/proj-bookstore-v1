package middleware

import (
	"crypto/subtle"

	"github.com/aniruddha-jafa/go-auth-v1/internal/apperrors"
	"github.com/aniruddha-jafa/go-auth-v1/internal/constants"
	"github.com/gofiber/fiber/v3"
)

func CSRF(c fiber.Ctx) error {
	csrfToken := c.Get(constants.CSRF_TOKEN_HEADER, "")
	if csrfToken == "" {
		return apperrors.ErrUnauthorized
	}
	cookieCsrfToken := c.Cookies(constants.CSRF_TOKEN_COOKIE_NAME, "")
	if cookieCsrfToken == "" {
		return apperrors.ErrUnauthorized
	}
	if subtle.ConstantTimeCompare([]byte(csrfToken), []byte(cookieCsrfToken)) != 1 {
		return apperrors.ErrUnauthorized
	}
	return c.Next()
}
