package config

import (
	"log"
	"os"

	"github.com/streadway/amqp"
)

func SetupRabbitMQQueues() {
	conn, err := amqp.Dial(os.Getenv("RABBITMQ_URL"))
	if err != nil {
		log.Fatalf("RabbitMQ connection failed: %v", err)
	}
	defer conn.Close()

	ch, err := conn.Channel()
	if err != nil {
		log.Fatalf("RabbitMQ channel failed: %v", err)
	}
	defer ch.Close()

	// USER CREATED queues
	declareQueue(ch, "user.created.queue")
	declareQueue(ch, "user.created.delay.queue")
	declareQueue(ch, "user.created.timeout.queue")
	declareQueue(ch, "user.created.failed.queue")

	// USER UPDATED queues
	declareQueue(ch, "user.updated.queue")
	declareQueue(ch, "user.updated.delay.queue")
	declareQueue(ch, "user.updated.timeout.queue")
	declareQueue(ch, "user.updated.failed.queue")

	log.Println("RabbitMQ queues created successfully")
}

func declareQueue(ch *amqp.Channel, name string) {
	_, err := ch.QueueDeclare(
		name,
		true,  // durable
		false, // auto delete
		false, // exclusive
		false, // no wait
		nil,
	)
	if err != nil {
		log.Fatalf("Failed to declare queue %s: %v", name, err)
	}
}
