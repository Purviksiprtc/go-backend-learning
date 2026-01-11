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
		log.Fatal(err)
	}

	wp, _ := workerpool.NewWorkerPool(context.Background(), 10, "user-updated")

	wp.RegisterTasks(map[string]interface{}{
		"USER_UPDATED": handleUserUpdated,
	})

	log.Println("USER_UPDATED consumer started")
	wp.Start()
}

func handleUserUpdated(sig *schema.Signature) error {
	var payload map[string]interface{}
	json.Unmarshal(sig.RawArgs, &payload)

	id := payload["data"].(map[string]interface{})["user_id"]
	fmt.Println("USER UPDATED EVENT RECEIVED →", id)

	return nil
}
