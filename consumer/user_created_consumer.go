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

	if err := paotaConfig.GetConfigProvider().SetApplicationConfig(cfg); err != nil {
		log.Fatal(err)
	}

	wp, _ := workerpool.NewWorkerPool(context.Background(), 10, "user-created")

	wp.RegisterTasks(map[string]interface{}{
		"USER_CREATED": handleUserCreated,
	})

	log.Println("USER_CREATED consumer started")
	wp.Start()
}

func handleUserCreated(sig *schema.Signature) error {
	var payload map[string]interface{}
	json.Unmarshal(sig.RawArgs, &payload)

	email := payload["data"].(map[string]interface{})["email"].(string)
	fmt.Println("USER CREATED EVENT RECEIVED →", email)

	return nil
}
