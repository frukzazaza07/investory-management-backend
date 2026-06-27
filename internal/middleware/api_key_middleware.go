package middleware

import (
	"os"

	"github.com/gofiber/fiber/v3"
)

func POSAPIKey(c fiber.Ctx) error {
	key := c.Get("X-API-Key")
	if key == "" {
		key = c.Query("api_key")
	}
	expected := os.Getenv("POS_API_KEY")
	if expected == "" || key != expected {
		return c.Status(fiber.StatusUnauthorized).JSON(fiber.Map{"error": "invalid api key"})
	}
	return c.Next()
}
