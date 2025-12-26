package main

import (
	"log"
	"os"

	"user-crud-go/config"
	"user-crud-go/models"
	"user-crud-go/routes"

	"github.com/labstack/echo/v4"
)

func main() {
	config.ConnectDB()
	config.DB.AutoMigrate(&models.User{})

	e := echo.New()
	routes.RegisterRoutes(e)

	port := os.Getenv("APP_PORT")
	if port == "" {
		port = "8080" // fallback safety
	}

	if err := e.Start(":" + port); err != nil {
		log.Fatalf("failed to start server: %v", err)
	}
}
