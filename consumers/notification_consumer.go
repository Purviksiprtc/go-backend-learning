package main

import (
	"context"
	"encoding/json"
	"fmt"
	"log"

	"user-crud-go/messaging/paota"

	"github.com/joho/godotenv"
	"github.com/surendratiwari3/paota/config"
	"github.com/surendratiwari3/paota/schema"
	"github.com/surendratiwari3/paota/workerpool"
)

func main() {
	// 🔴 VERY IMPORTANT: Load .env for RabbitMQ config
	if err := godotenv.Load(); err != nil {
		log.Println("⚠️ No .env file found, using system environment variables")
	}

	// Load Paota config (now env vars WILL be available)
	cfg := paota.LoadPaotaConfig(
		"user.created.queue",
		"user.created",
	)

	// 🔴 REQUIRED: set application config BEFORE creating worker pool
	if err := config.GetConfigProvider().SetApplicationConfig(cfg); err != nil {
		log.Fatalf("failed to set paota config: %v", err)
	}

	// Create worker pool
	wp, err := workerpool.NewWorkerPool(
		context.Background(),
		10,
		"notification-consumer",
	)
	if err != nil {
		log.Fatalf("failed to create worker pool: %v", err)
	}

	// Register task handlers
	if err := wp.RegisterTasks(map[string]interface{}{
		"USER_CREATED": handleUserCreated,
		"USER_UPDATED": handleUserUpdated,
	}); err != nil {
		log.Fatalf("failed to register tasks: %v", err)
	}

	log.Println("📨 Notification consumer started...")
	if err := wp.Start(); err != nil {
		log.Fatalf("failed to start worker pool: %v", err)
	}
}

func handleUserCreated(arg *schema.Signature) error {
	var payload paota.UserEventPayload
	if err := json.Unmarshal(arg.RawArgs, &payload); err != nil {
		return err // retry will happen
	}

	fmt.Printf(
		"[USER_CREATED] Welcome email sent to %s\n",
		payload.Data.Email,
	)
	return nil // ✅ ACK
}

func handleUserUpdated(arg *schema.Signature) error {
	var payload paota.UserEventPayload
	if err := json.Unmarshal(arg.RawArgs, &payload); err != nil {
		return err // retry will happen
	}

	fmt.Printf(
		"[USER_UPDATED] User %d profile updated\n",
		payload.Data.UserID,
	)
	return nil // ✅ ACK
}
