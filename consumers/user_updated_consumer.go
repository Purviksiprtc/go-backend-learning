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
	godotenv.Load()

	cfg := paota.LoadPaotaConfig(
		"user.updated.queue",
		"user.updated",
	)

	config.GetConfigProvider().SetApplicationConfig(cfg)

	wp, _ := workerpool.NewWorkerPool(
		context.Background(),
		10,
		"user-updated-consumer",
	)

	wp.RegisterTasks(map[string]interface{}{
		"USER_UPDATED": handleUserUpdated,
	})

	log.Println("📨 USER_UPDATED consumer started")
	wp.Start()
}

func handleUserUpdated(arg *schema.Signature) error {
	var payload paota.UserEventPayload
	json.Unmarshal(arg.RawArgs, &payload)

	fmt.Printf(
		"[USER_UPDATED] User %d profile updated\n",
		payload.Data.UserID,
	)

	return nil
}
