package rabbitmq

import (
	"context"
	"encoding/json"
	"fmt"
	"os"
	"time"

	amqp "github.com/rabbitmq/amqp091-go"
)

type UserEvent struct {
	Event     string   `json:"event"`
	Version   string   `json:"version"`
	Timestamp string   `json:"timestamp"`
	Data      UserData `json:"data"`
}

type UserData struct {
	UserID uint   `json:"user_id"`
	Name   string `json:"name"`
	Email  string `json:"email"`
}

func PublishEvent(routingKey string, event UserEvent) error {
	// ✅ Improvement: validate exchange early
	exchange := os.Getenv("RABBITMQ_EXCHANGE")
	if exchange == "" {
		return fmt.Errorf("RABBITMQ_EXCHANGE not set")
	}

	ctx, cancel := context.WithTimeout(context.Background(), 5*time.Second)
	defer cancel()

	body, err := json.Marshal(event)
	if err != nil {
		return err
	}

	return Channel.PublishWithContext(
		ctx,
		exchange,
		routingKey,
		false,
		false,
		amqp.Publishing{
			ContentType:  "application/json",
			DeliveryMode: amqp.Persistent,
			Body:         body,
		},
	)
}
