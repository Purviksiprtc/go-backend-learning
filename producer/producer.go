package producer

import (
	"context"
	"encoding/json"
	"os"
	"time"

	appConfig "user-crud-go/config"

	"github.com/surendratiwari3/paota/schema"
	"github.com/surendratiwari3/paota/workerpool"
)

type UserProducer struct {
	pool *workerpool.Pool
}

/* ---------- USER CREATED PRODUCER ---------- */

func NewUserCreatedProducer() (*UserProducer, error) {
	queue := os.Getenv("USER_CREATED_QUEUE")

	cfg := appConfig.LoadPaotaConfig(queue, "user.created")

	wp, err := workerpool.NewWorkerPoolWithConfig(
		context.Background(),
		5,
		queue,
		cfg,
	)
	if err != nil {
		return nil, err
	}

	return &UserProducer{pool: &wp}, nil
}

/* ---------- USER UPDATED PRODUCER ---------- */

func NewUserUpdatedProducer() (*UserProducer, error) {
	queue := os.Getenv("USER_UPDATED_QUEUE")

	cfg := appConfig.LoadPaotaConfig(queue, "user.updated")

	wp, err := workerpool.NewWorkerPoolWithConfig(
		context.Background(),
		5,
		queue,
		cfg,
	)
	if err != nil {
		return nil, err
	}

	return &UserProducer{pool: &wp}, nil
}

/* ---------- PUBLISH ---------- */

func (p *UserProducer) Publish(
	event string,
	userID uint,
	name, email, routingKey string,
) error {

	payload := map[string]interface{}{
		"event":     event,
		"version":   "1.0",
		"timestamp": time.Now().UTC(),
		"data": map[string]interface{}{
			"user_id": userID,
			"name":    name,
			"email":   email,
		},
	}

	body, err := json.Marshal(payload)
	if err != nil {
		return err
	}

	task := &schema.Signature{
		Name:       event,
		RawArgs:    body,
		RoutingKey: routingKey,
		RetryCount: 5,
	}

	_, err = (*p.pool).SendTaskWithContext(context.Background(), task)
	return err
}
