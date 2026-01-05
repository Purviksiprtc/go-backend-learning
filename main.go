package main

import (
	"log"
	"os"

	"user-crud-go/config"
	"user-crud-go/models"
	"user-crud-go/routes"

	"github.com/go-playground/validator/v10"
	"github.com/joho/godotenv"
	"github.com/labstack/echo/v4"
)

func main() {
	// ===============================
	// Load Environment Variables
	// ===============================
	if err := godotenv.Load(); err != nil {
		log.Println("⚠️ No .env file found, using system environment variables")
	}

	// ===============================
	// Database
	// ===============================
	config.ConnectDB()
	config.DB.AutoMigrate(&models.User{})

	// ===============================
	// Echo Server
	// ===============================
	e := echo.New()

	// Attach request validator
	e.Validator = &config.CustomValidator{
		Validator: validator.New(),
	}

	// Register routes
	routes.RegisterRoutes(e)

	// ===============================
	// Server Port
	// ===============================
	port := os.Getenv("APP_PORT")
	if port == "" {
		port = "8080"
	}

	log.Printf("🚀 API server started on port %s\n", port)

	if err := e.Start(":" + port); err != nil {
		log.Fatalf("❌ Failed to start server: %v", err)
	}
}
