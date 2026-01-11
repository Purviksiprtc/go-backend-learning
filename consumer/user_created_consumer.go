package consumer

import (
	"context"
	"encoding/json"
	"log"

	appConfig "user-crud-go/config"

	paotaConfig "github.com/surendratiwari3/paota/config"
	"github.com/surendratiwari3/paota/schema"
	"github.com/surendratiwari3/paota/workerpool"
)

func StartUserCreatedConsumer() {

	cfg := appConfig.LoadPaotaConfig(
		// main queue
		appConfig.GetEnv("USER_CREATED_QUEUE"),
		// routing key
		"user.created",
	)

	// 🔥 REQUIRED: register config with Paota
	if err := paotaConfig.GetConfigProvider().SetApplicationConfig(cfg); err != nil {
		log.Fatalf("❌ Failed to set paota config: %v", err)
	}

	wp, err := workerpool.NewWorkerPool(
		context.Background(),
		10,
		"user-created-consumer",
	)
	if err != nil {
		log.Fatalf("❌ Failed to create worker pool: %v", err)
	}

	if err := wp.RegisterTasks(map[string]interface{}{
		"USER_CREATED": handleUserCreated,
	}); err != nil {
		log.Fatalf("❌ Failed to register task: %v", err)
	}

	log.Println("📨 USER_CREATED consumer started")
	if err := wp.Start(); err != nil {
		log.Fatalf("❌ Failed to start worker pool: %v", err)
	}
}

func handleUserCreated(task *schema.Signature) error {
	var payload map[string]interface{}

	if err := json.Unmarshal(task.RawArgs, &payload); err != nil {
		return err
	}

	data := payload["data"].(map[string]interface{})
	email := data["email"].(string)

	log.Printf("✅ USER_CREATED received → welcome email sent to %s\n", email)
	return nil // ACK
}
