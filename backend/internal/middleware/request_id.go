package middleware

import (
	"context"

	"github.com/gofiber/fiber/v2"
	"github.com/google/uuid"
)

const RequestIdKey ctxKey = "requestId"

// Adds a request ID to the context
func RequestID(c *fiber.Ctx) error {
	requestID := c.Get("X-Request-ID", "")
	if requestID == "" {
		requestID = uuid.New().String()
	}

	ctx := context.WithValue(c.UserContext(), RequestIdKey, requestID)
	c.SetUserContext(ctx)

	c.Set("X-Request-Id", requestID)

	return c.Next()
}
