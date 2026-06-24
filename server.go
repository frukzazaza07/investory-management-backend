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
	"investory-management-backend/internal/router"

	"github.com/gofiber/fiber/v3"
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
	router.SetupRoutes(app)

	port := os.Getenv("APP_PORT")
	if port == "" {
		port = "8080"
	}

	log.Fatal(app.Listen(":" + port))
}
