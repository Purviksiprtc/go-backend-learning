package main

import (
	"log"
	"os"

	"user-crud-go/config"
	"user-crud-go/consumer"
	"user-crud-go/models"
	"user-crud-go/routes"

	"github.com/go-playground/validator/v10"
	"github.com/joho/godotenv"
	"github.com/labstack/echo/v4"
)

/*
APP_MODE controls what runs:
- api
- consumer_created
- consumer_updated
*/

func main() {
	// ===============================
	// Load Environment Variables
	// ===============================
	if err := godotenv.Load(); err != nil {
		log.Println("⚠️ No .env file found, using system environment variables")
	}

	mode := os.Getenv("APP_MODE")
	if mode == "" {
		mode = "api"
	}

	switch mode {

	case "api":
		startAPI()

	case "consumer_created":
		log.Println("🚀 Starting USER_CREATED consumer")
		consumer.StartUserCreatedConsumer()

	case "consumer_updated":
		log.Println("🚀 Starting USER_UPDATED consumer")
		consumer.StartUserUpdatedConsumer()

	default:
		log.Fatalf("❌ Invalid APP_MODE: %s", mode)
	}
}

/* ===============================
   API STARTUP
================================ */

func startAPI() {
	// ===============================
	// Database
	// ===============================
	config.ConnectDB()

	if err := config.DB.AutoMigrate(&models.User{}); err != nil {
		log.Fatalf("❌ Database migration failed: %v", err)
	}

	// ===============================
	// Echo Server
	// ===============================
	e := echo.New()

	e.Validator = &config.CustomValidator{
		Validator: validator.New(),
	}

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
		log.Fatalf("❌ Failed to start API server: %v", err)
	}
}
