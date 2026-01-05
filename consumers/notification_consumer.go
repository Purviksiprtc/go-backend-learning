package main

import (
	"context"
	"encoding/json"
	"fmt"

	"github.com/surendratiwari3/paota/schema"
	"github.com/surendratiwari3/paota/workerpool"

	"user-crud-go/messaging/paota"
)

func main() {
	cfg := paota.LoadPaotaConfig(
		"user.created.queue",
		"user.created",
	)

	wp, _ := workerpool.NewWorkerPoolWithConfig(
		context.Background(),
		10,
		"notification-service",
		cfg,
	)

	wp.RegisterTasks(map[string]interface{}{
		"USER_CREATED": handleUserCreated,
		"USER_UPDATED": handleUserUpdated,
	})

	wp.Start()
}

func handleUserCreated(arg *schema.Signature) error {
	var payload paota.UserEventPayload
	json.Unmarshal(arg.RawArgs, &payload)

	fmt.Printf("[USER_CREATED] Welcome email sent to %s\n", payload.Data.Email)
	return nil
}

func handleUserUpdated(arg *schema.Signature) error {
	var payload paota.UserEventPayload
	json.Unmarshal(arg.RawArgs, &payload)

	fmt.Printf("[USER_UPDATED] User %d profile updated\n", payload.Data.UserID)
	return nil
}
