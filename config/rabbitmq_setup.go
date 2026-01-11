package config

import (
	"log"
	"os"

	"github.com/streadway/amqp"
)

func SetupQueues() {
	conn, err := amqp.Dial(os.Getenv("RABBITMQ_URL"))
	if err != nil {
		log.Fatal(err)
	}
	defer conn.Close()

	ch, err := conn.Channel()
	if err != nil {
		log.Fatal(err)
	}
	defer ch.Close()

	exchange := os.Getenv("RABBITMQ_EXCHANGE")

	// exchange
	err = ch.ExchangeDeclare(
		exchange,
		"topic",
		true,
		false,
		false,
		false,
		nil,
	)
	if err != nil {
		log.Fatal(err)
	}

	// 🔥 ONLY declare delay / timeout / failed queues
	declareDLXQueue(ch, os.Getenv("USER_CREATED_DELAY_QUEUE"), exchange, "user.created")
	declareDLXQueue(ch, os.Getenv("USER_CREATED_TIMEOUT_QUEUE"), exchange, "user.created")
	declareDLXQueue(ch, os.Getenv("USER_CREATED_FAILED_QUEUE"), exchange, "user.created")

	declareDLXQueue(ch, os.Getenv("USER_UPDATED_DELAY_QUEUE"), exchange, "user.updated")
	declareDLXQueue(ch, os.Getenv("USER_UPDATED_TIMEOUT_QUEUE"), exchange, "user.updated")
	declareDLXQueue(ch, os.Getenv("USER_UPDATED_FAILED_QUEUE"), exchange, "user.updated")

	log.Println("✅ RabbitMQ queues initialized (manager-style)")
}

func declareDLXQueue(ch *amqp.Channel, name, exchange, routingKey string) {
	_, err := ch.QueueDeclare(
		name,
		true,
		false,
		false,
		false,
		amqp.Table{
			"x-dead-letter-exchange":    exchange,
			"x-dead-letter-routing-key": routingKey,
		},
	)
	if err != nil {
		log.Fatal(err)
	}
}
