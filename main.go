package main

import (
	"log"
	"os"

	"user-crud-go/config"
	"user-crud-go/models"
	"user-crud-go/routes"

	"github.com/go-playground/validator/v10"
	"github.com/labstack/echo/v4"
)

func main() {
	// DB
	config.ConnectDB()
	config.DB.AutoMigrate(&models.User{})

	// Echo
	e := echo.New()

	// ✅ Attach validator
	e.Validator = &config.CustomValidator{
		Validator: validator.New(),
	}

	// Routes
	routes.RegisterRoutes(e)

	// Port
	port := os.Getenv("APP_PORT")
	if port == "" {
		port = "8080"
	}

	if err := e.Start(":" + port); err != nil {
		log.Fatalf("failed to start server: %v", err)
	}
}
