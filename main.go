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

func main() {
	godotenv.Load()

	// MUST run once
	config.SetupQueues()

	mode := os.Getenv("APP_MODE")
	if mode == "" {
		mode = "api"
	}

	switch mode {
	case "api":
		startAPI()
	case "consumer_created":
		consumer.StartUserCreatedConsumer()
	case "consumer_updated":
		consumer.StartUserUpdatedConsumer()
	default:
		log.Fatal("Invalid APP_MODE")
	}
}

func startAPI() {
	// Database
	config.ConnectDB()
	if err := config.DB.AutoMigrate(&models.User{}); err != nil {
		log.Fatalf("❌ Migration failed: %v", err)
	}

	// Echo
	e := echo.New()
	e.Validator = &config.CustomValidator{
		Validator: validator.New(),
	}

	routes.RegisterRoutes(e)

	port := os.Getenv("APP_PORT")
	if port == "" {
		port = "8080"
	}

	log.Printf("🚀 API server started on port %s\n", port)
	log.Fatal(e.Start(":" + port))
}
