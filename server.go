package main

// @title           Inventory Management API
// @version         1.0
// @description     API for inventory management
// @host            localhost:8080
// @BasePath        /
// @securityDefinitions.apikey BearerAuth
// @in header
// @name Authorization

import (
	"log"
	"os"

	_ "investory-management-backend/docs"
	"investory-management-backend/internal/database"
	"investory-management-backend/internal/middleware"
	"investory-management-backend/internal/router"

	"github.com/gofiber/fiber/v3"
	"github.com/gofiber/fiber/v3/middleware/cors"
	"github.com/joho/godotenv"
)

func main() {
	godotenv.Load()
	database.Connect()
	database.Migrate()
	if os.Getenv("SEED") == "true" {
		database.Seed()
	}

	app := fiber.New()

	app.Use(middleware.Language)

	app.Use(cors.New(cors.Config{
		AllowOrigins: []string{
			"http://localhost:5173",
			"https://investory-management-frontend-chi.vercel.app",
		},
		AllowHeaders: []string{"Origin", "Content-Type", "Accept", "Authorization", "X-API-Key"},
		AllowMethods: []string{"GET", "POST", "PUT", "PATCH", "DELETE", "OPTIONS"},
	}))

	router.SetupRoutes(app)

	port := os.Getenv("APP_PORT")
	if port == "" {
		port = "8080"
	}

	certFile := os.Getenv("TLS_CERT")
	keyFile := os.Getenv("TLS_KEY")

	if certFile != "" && keyFile != "" {
		log.Printf("listening on https://localhost:%s", port)
		// log.Fatal(app.Listen(":"+port, fiber.ListenConfig{CertFile: certFile, CertKeyFile: keyFile}))
	} else {
		log.Printf("listening on http://localhost:%s", port)
		log.Fatal(app.Listen(":" + port))
	}
}
