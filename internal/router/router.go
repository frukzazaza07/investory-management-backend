package router

import (
	"investory-management-backend/internal/handler"
	"investory-management-backend/internal/middleware"

	"github.com/gofiber/fiber/v3"
	"github.com/swaggo/swag"
)

func SetupRoutes(app *fiber.App) {
	// Swagger JSON
	app.Get("/swagger/doc.json", func(c fiber.Ctx) error {
		doc, err := swag.ReadDoc()
		if err != nil {
			return c.Status(500).SendString("swagger doc not found")
		}
		c.Set("Content-Type", "application/json")
		return c.SendString(doc)
	})

	// Scalar UI
	app.Get("/docs", func(c fiber.Ctx) error {
		html := `<!DOCTYPE html>
<html>
<head><title>API Docs</title></head>
<body>
<script id="api-reference" data-url="/swagger/doc.json"></script>
<script src="https://cdn.jsdelivr.net/npm/@scalar/api-reference"></script>
</body>
</html>`
		c.Set("Content-Type", "text/html")
		return c.SendString(html)
	})

	auth := app.Group("/auth")
	auth.Post("/register", handler.Register)
	auth.Post("/login", handler.Login)

	api := app.Group("/api", middleware.Protected)
	api.Get("/me", func(c fiber.Ctx) error {
		return c.JSON(fiber.Map{
			"user_id": c.Locals("userID"),
			"email":   c.Locals("email"),
		})
	})
}
