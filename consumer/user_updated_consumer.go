package consumer

import (
	"context"
	"encoding/json"
	"fmt"
	"log"

	appConfig "user-crud-go/config"

	paotaConfig "github.com/surendratiwari3/paota/config"
	"github.com/surendratiwari3/paota/schema"
	"github.com/surendratiwari3/paota/workerpool"
)

func StartUserUpdatedConsumer() {
	cfg := appConfig.LoadPaota(
		"user.updated.queue",
		"user.updated",
	)

	if err := paotaConfig.GetConfigProvider().SetApplicationConfig(cfg); err != nil {
		log.Fatalf("failed to set paota config: %v", err)
	}

	wp, err := workerpool.NewWorkerPool(
		context.Background(),
		10,
		"user-updated-consumer",
	)
	if err != nil {
		log.Fatalf("failed to create worker pool: %v", err)
	}

	if err := wp.RegisterTasks(map[string]interface{}{
		"USER_UPDATED": handleUserUpdated,
	}); err != nil {
		log.Fatalf("failed to register task: %v", err)
	}

	log.Println("📨 USER_UPDATED consumer started")
	if err := wp.Start(); err != nil {
		log.Fatalf("failed to start worker pool: %v", err)
	}
}

func handleUserUpdated(arg *schema.Signature) error {
	var payload map[string]interface{}
	if err := json.Unmarshal(arg.RawArgs, &payload); err != nil {
		return err
	}

	data := payload["data"].(map[string]interface{})
	userID := int(data["user_id"].(float64))

	fmt.Printf("[USER_UPDATED] User %d profile updated\n", userID)
	return nil // ✅ ACK
}
