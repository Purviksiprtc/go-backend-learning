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
		"user.created.queue",
		"user.created",
	)

	config.GetConfigProvider().SetApplicationConfig(cfg)

	wp, _ := workerpool.NewWorkerPool(
		context.Background(),
		10,
		"user-created-consumer",
	)

	wp.RegisterTasks(map[string]interface{}{
		"USER_CREATED": handleUserCreated,
	})

	log.Println("📨 USER_CREATED consumer started")
	wp.Start()
}

func handleUserCreated(arg *schema.Signature) error {
	var payload paota.UserEventPayload
	json.Unmarshal(arg.RawArgs, &payload)

	fmt.Printf(
		"[USER_CREATED] Welcome email sent to %s\n",
		payload.Data.Email,
	)

	return nil
}
