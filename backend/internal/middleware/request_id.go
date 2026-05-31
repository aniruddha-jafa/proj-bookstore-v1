package middleware

import (
	"context"

	"github.com/aniruddha-jafa/go-auth-v1/internal/constants"
	"github.com/aniruddha-jafa/go-auth-v1/internal/request_context"
	"github.com/gofiber/fiber/v2"
	"github.com/google/uuid"
)

// Adds a request ID to the context
func RequestID(c *fiber.Ctx) error {
	requestID := c.Get(constants.REQUEST_ID_HEADER, "")
	if requestID == "" {
		requestID = uuid.New().String()
	}

	ctx := context.WithValue(c.UserContext(), request_context.RequestIdKey, requestID)
	c.SetUserContext(ctx)

	c.Set(constants.REQUEST_ID_HEADER, requestID)

	return c.Next()
}
