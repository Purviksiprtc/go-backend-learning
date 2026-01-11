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

	// 🔥 CREATE QUEUES FIRST (manager requirement)
	config.SetupRabbitMQQueues()

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
	config.ConnectDB()
	config.DB.AutoMigrate(&models.User{})

	e := echo.New()
	e.Validator = &config.CustomValidator{Validator: validator.New()}
	routes.RegisterRoutes(e)

	log.Println("API server started on port 8080")
	e.Start(":8080")
}
