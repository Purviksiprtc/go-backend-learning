package main

import (
	"user-crud-go/consumer"

	"log"
	"os"
	"user-crud-go/rabbitmq"

	"user-crud-go/config"
	"user-crud-go/models"
	"user-crud-go/routes"

	"github.com/labstack/echo/v4"
)

func main() {
	config.ConnectDB()
	config.DB.AutoMigrate(&models.User{})

	if err := rabbitmq.ConnectRabbitMQ(); err != nil {
		log.Fatalf("RabbitMQ connection failed: %v", err)
	}

	if err := rabbitmq.SetupRabbitMQ(rabbitmq.Channel); err != nil {
		log.Fatalf("RabbitMQ setup failed: %v", err)
	}

	if err := consumer.StartConsumer(rabbitmq.Channel); err != nil {
		log.Fatalf("Consumer failed to start: %v", err)
	}

	e := echo.New()
	routes.RegisterRoutes(e)

	port := os.Getenv("APP_PORT")
	if port == "" {
		port = "8080"
	}

	if err := e.Start(":" + port); err != nil {
		log.Fatalf("failed to start server: %v", err)
	}
}
