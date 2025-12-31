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
	// 🔹 DB setup
	config.ConnectDB()
	config.DB.AutoMigrate(&models.User{})

	// 🔹 RabbitMQ connection
	if err := rabbitmq.ConnectRabbitMQ(); err != nil {
		log.Fatalf("RabbitMQ connection failed: %v", err)
	}

	if err := rabbitmq.SetupRabbitMQ(rabbitmq.Channel); err != nil {
		log.Fatalf("RabbitMQ setup failed: %v", err)
	}

	if err := consumer.StartConsumer(rabbitmq.Channel); err != nil {
		log.Fatalf("Consumer failed to start: %v", err)
	}

	// 🔹 Echo setup
	e := echo.New()

	// ✅ Structured request validation (PR requirement)
	e.Validator = &config.CustomValidator{
		Validator: validator.New(),
	}

	// 🔹 Routes
	routes.RegisterRoutes(e)

	// 🔹 Server start
	port := os.Getenv("APP_PORT")
	if port == "" {
		port = "8080"
	}

	if err := e.Start(":" + port); err != nil {
		log.Fatalf("failed to start server: %v", err)
	}
}
