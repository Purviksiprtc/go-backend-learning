package consumer

import (
	"encoding/json"
	"log"
	"os"

	"user-crud-go/rabbitmq"

	amqp "github.com/rabbitmq/amqp091-go"
)

func StartConsumer(ch *amqp.Channel) error {
	createdQueue := os.Getenv("RABBITMQ_USER_CREATED_QUEUE")
	updatedQueue := os.Getenv("RABBITMQ_USER_UPDATED_QUEUE")

	// Consume USER_CREATED
	createdMsgs, err := ch.Consume(
		createdQueue,
		"",
		false, // auto-ack = false (IMPORTANT)
		false,
		false,
		false,
		nil,
	)
	if err != nil {
		return err
	}

	// Consume USER_UPDATED
	updatedMsgs, err := ch.Consume(
		updatedQueue,
		"",
		false,
		false,
		false,
		false,
		nil,
	)
	if err != nil {
		return err
	}

	// Handle USER_CREATED messages
	go func() {
		for msg := range createdMsgs {
			var event rabbitmq.UserEvent

			if err := json.Unmarshal(msg.Body, &event); err != nil {
				log.Println("Failed to parse USER_CREATED event:", err)
				msg.Nack(false, false) // send to DLQ
				continue
			}

			log.Printf("[USER_CREATED] Welcome email sent to %s\n", event.Data.Email)
			log.Printf("[USER_CREATED] Audit log created for user %d\n", event.Data.UserID)

			msg.Ack(false)
		}
	}()

	// Handle USER_UPDATED messages
	go func() {
		for msg := range updatedMsgs {
			var event rabbitmq.UserEvent

			if err := json.Unmarshal(msg.Body, &event); err != nil {
				log.Println("Failed to parse USER_UPDATED event:", err)
				msg.Nack(false, false)
				continue
			}

			log.Printf("[USER_UPDATED] User %d profile updated\n", event.Data.UserID)
			log.Printf("[USER_UPDATED] Audit log updated for %s\n", event.Data.Email)

			msg.Ack(false)
		}
	}()

	log.Println("RabbitMQ consumer started successfully")
	return nil
}
