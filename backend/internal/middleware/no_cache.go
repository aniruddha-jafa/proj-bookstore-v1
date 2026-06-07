package middleware

import "github.com/gofiber/fiber/v3"

func NoCache(c fiber.Ctx) error {
	c.Set("Cache-Control", "no-cache, no-store, must-revalidate, max-age=0")
	c.Set("Pragma", "no-cache")
	return c.Next()
}
