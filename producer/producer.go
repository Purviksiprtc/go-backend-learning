package producer

import (
	"context"
	"encoding/json"
	"time"

	"user-crud-go/config"

	"github.com/surendratiwari3/paota/schema"
	"github.com/surendratiwari3/paota/workerpool"
)

type UserProducer struct {
	pool *workerpool.Pool
}

func NewUserProducer(queue, routingKey string) (*UserProducer, error) {
	cfg := config.LoadPaota(queue, routingKey)

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
