package main

import (
	"log"
	"os"

	"user-crud-go/config"
	"user-crud-go/consumer"
	"user-crud-go/models"
	"user-crud-go/rabbitmq"
	"user-crud-go/routes"

	"github.com/go-playground/validator/v10"
	"github.com/labstack/echo/v4"
)

func main() {
	// 🔹 DB setup (needed for both API & Consumer)
	config.ConnectDB()
	config.DB.AutoMigrate(&models.User{})

	// 🔹 RabbitMQ connection (needed for both)
	if err := rabbitmq.ConnectRabbitMQ(); err != nil {
		log.Fatalf("RabbitMQ connection failed: %v", err)
	}

	if err := rabbitmq.SetupRabbitMQ(rabbitmq.Channel); err != nil {
		log.Fatalf("RabbitMQ setup failed: %v", err)
	}

	// 🔹 Decide runtime mode
	appMode := os.Getenv("APP_MODE")
	if appMode == "" {
		appMode = "api" // default safety
	}

	// =========================================================
	// 🔹 CONSUMER MODE
	// =========================================================
	if appMode == "consumer" {
		log.Println("Starting application in CONSUMER mode")

		if err := consumer.StartConsumer(rabbitmq.Channel); err != nil {
			log.Fatalf("Consumer failed to start: %v", err)
		}

		// Block forever so container doesn’t exit
		select {}
	}

	// =========================================================
	// 🔹 API MODE
	// =========================================================
	if appMode == "api" {
		log.Println("Starting application in API mode")

		e := echo.New()

		// Structured request validation (PR requirement)
		e.Validator = &config.CustomValidator{
			Validator: validator.New(),
		}

		routes.RegisterRoutes(e)

		port := os.Getenv("APP_PORT")
		if port == "" {
			port = "8080"
		}

		if err := e.Start(":" + port); err != nil {
			log.Fatalf("failed to start server: %v", err)
		}
	}
}
