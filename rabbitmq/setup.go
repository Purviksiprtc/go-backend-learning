package rabbitmq

import (
	"fmt"
	"os"

	amqp "github.com/rabbitmq/amqp091-go"
)

func SetupRabbitMQ(ch *amqp.Channel) error {

	exchange := os.Getenv("RABBITMQ_EXCHANGE")
	userCreatedQueue := os.Getenv("RABBITMQ_USER_CREATED_QUEUE")
	userUpdatedQueue := os.Getenv("RABBITMQ_USER_UPDATED_QUEUE")
	dlq := os.Getenv("RABBITMQ_DLQ")

	if exchange == "" || userCreatedQueue == "" || userUpdatedQueue == "" || dlq == "" {
		return fmt.Errorf("rabbitmq env variables not set properly")
	}

	// 1️⃣ Declare Exchange
	err := ch.ExchangeDeclare(
		exchange,
		"topic",
		true,  // durable
		false, // auto-delete
		false,
		false,
		nil,
	)
	if err != nil {
		return fmt.Errorf("failed to declare exchange: %w", err)
	}

	// 2️⃣ Declare DLQ
	_, err = ch.QueueDeclare(
		dlq,
		true,
		false,
		false,
		false,
		nil,
	)
	if err != nil {
		return fmt.Errorf("failed to declare DLQ: %w", err)
	}

	// Common arguments for queues (DLQ attached)
	queueArgs := amqp.Table{
		"x-dead-letter-exchange":    "",
		"x-dead-letter-routing-key": dlq,
	}

	// 3️⃣ Declare USER_CREATED queue
	_, err = ch.QueueDeclare(
		userCreatedQueue,
		true,
		false,
		false,
		false,
		queueArgs,
	)
	if err != nil {
		return fmt.Errorf("failed to declare user created queue: %w", err)
	}

	// 4️⃣ Declare USER_UPDATED queue
	_, err = ch.QueueDeclare(
		userUpdatedQueue,
		true,
		false,
		false,
		false,
		queueArgs,
	)
	if err != nil {
		return fmt.Errorf("failed to declare user updated queue: %w", err)
	}

	// 5️⃣ Bind queues
	err = ch.QueueBind(
		userCreatedQueue,
		"user.created",
		exchange,
		false,
		nil,
	)
	if err != nil {
		return fmt.Errorf("failed to bind user.created queue: %w", err)
	}

	err = ch.QueueBind(
		userUpdatedQueue,
		"user.updated",
		exchange,
		false,
		nil,
	)
	if err != nil {
		return fmt.Errorf("failed to bind user.updated queue: %w", err)
	}

	fmt.Println("RabbitMQ setup completed successfully")
	return nil
}
