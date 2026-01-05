package paota

import (
	"context"
	"encoding/json"
	"time"

	"github.com/surendratiwari3/paota/schema"
	"github.com/surendratiwari3/paota/workerpool"
)

type EventProducer struct {
	pool *workerpool.Pool
}

func NewEventProducer(queue, routingKey string) (*EventProducer, error) {
	cfg := LoadPaotaConfig(queue, routingKey)

	wp, err := workerpool.NewWorkerPoolWithConfig(
		context.Background(),
		5,
		queue,
		cfg,
	)
	if err != nil {
		return nil, err
	}

	return &EventProducer{pool: &wp}, nil
}

func (p *EventProducer) PublishUserEvent(
	event string,
	userID uint,
	name, email, routingKey string,
) error {

	payload := UserEventPayload{
		Event:     event,
		Version:   "1.0",
		Timestamp: time.Now().UTC(),
		Data: UserData{
			UserID: userID,
			Name:   name,
			Email:  email,
		},
	}

	body, _ := json.Marshal(payload)

	task := &schema.Signature{
		Name:       event,
		RawArgs:    body,
		RoutingKey: routingKey,
		RetryCount: 5,
	}

	_, err := (*p.pool).SendTaskWithContext(context.Background(), task)
	return err
}
