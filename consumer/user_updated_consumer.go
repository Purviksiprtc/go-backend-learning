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

func StartUserUpdatedConsumer() {

	cfg := appConfig.LoadPaotaConfig(
		appConfig.GetEnv("USER_UPDATED_QUEUE"),
		"user.updated",
	)

	// 🔥 REQUIRED
	if err := paotaConfig.GetConfigProvider().SetApplicationConfig(cfg); err != nil {
		log.Fatalf("❌ Failed to set paota config: %v", err)
	}

	wp, err := workerpool.NewWorkerPool(
		context.Background(),
		10,
		"user-updated-consumer",
	)
	if err != nil {
		log.Fatalf("❌ Failed to create worker pool: %v", err)
	}

	if err := wp.RegisterTasks(map[string]interface{}{
		"USER_UPDATED": handleUserUpdated,
	}); err != nil {
		log.Fatalf("❌ Failed to register task: %v", err)
	}

	log.Println("📨 USER_UPDATED consumer started")
	if err := wp.Start(); err != nil {
		log.Fatalf("❌ Failed to start worker pool: %v", err)
	}
}

func handleUserUpdated(task *schema.Signature) error {
	var payload map[string]interface{}

	if err := json.Unmarshal(task.RawArgs, &payload); err != nil {
		return err
	}

	data := payload["data"].(map[string]interface{})
	id := int(data["user_id"].(float64))

	log.Printf("✅ USER_UPDATED received → user %d updated\n", id)
	return nil
}
