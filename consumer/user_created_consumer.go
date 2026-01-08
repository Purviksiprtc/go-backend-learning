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

func StartUserCreatedConsumer() {
	cfg := appConfig.LoadPaota(
		"user.created.queue",
		"user.created",
	)

	// Register Paota application config (REQUIRED)
	if err := paotaConfig.GetConfigProvider().SetApplicationConfig(cfg); err != nil {
		log.Fatalf("failed to set paota config: %v", err)
	}

	wp, err := workerpool.NewWorkerPool(
		context.Background(),
		10,
		"user-created-consumer",
	)
	if err != nil {
		log.Fatalf("failed to create worker pool: %v", err)
	}

	if err := wp.RegisterTasks(map[string]interface{}{
		"USER_CREATED": handleUserCreated,
	}); err != nil {
		log.Fatalf("failed to register task: %v", err)
	}

	log.Println("📨 USER_CREATED consumer started")
	if err := wp.Start(); err != nil {
		log.Fatalf("failed to start worker pool: %v", err)
	}
}

func handleUserCreated(arg *schema.Signature) error {
	var payload map[string]interface{}
	if err := json.Unmarshal(arg.RawArgs, &payload); err != nil {
		return err
	}

	data := payload["data"].(map[string]interface{})
	email := data["email"].(string)

	fmt.Printf("[USER_CREATED] Welcome email sent to %s\n", email)
	return nil // ✅ ACK
}
