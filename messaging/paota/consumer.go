package paota

import (
	"context"
	"fmt"

	"github.com/surendratiwari3/paota/config"
	"github.com/surendratiwari3/paota/schema"
	"github.com/surendratiwari3/paota/workerpool"
)

/*
Consumer wraps Paota worker pool for reusable consumption logic.
This is NOT an executable.
*/
type Consumer struct {
	workerPool *workerpool.Pool
}

// NewConsumer initializes a Paota consumer for a given queue and routing key
func NewConsumer(queueName, bindingKey string) (*Consumer, error) {
	cfg := LoadPaotaConfig(queueName, bindingKey)

	// Set application-level config (required by Paota)
	if err := config.GetConfigProvider().SetApplicationConfig(cfg); err != nil {
		return nil, fmt.Errorf("failed to set paota config: %w", err)
	}

	wp, err := workerpool.NewWorkerPool(
		context.Background(),
		10, // concurrency
		queueName,
	)
	if err != nil {
		return nil, err
	}

	return &Consumer{
		workerPool: &wp,
	}, nil
}

// RegisterHandlers registers task handlers
func (c *Consumer) RegisterHandlers(handlers map[string]interface{}) error {
	return (*c.workerPool).RegisterTasks(handlers)
}

// Start starts consuming messages
func (c *Consumer) Start() error {
	return (*c.workerPool).Start()
}

// Ack helper (Paota auto-acks on nil error)
func Ack(_ *schema.Signature) error {
	return nil
}
