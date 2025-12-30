package rabbitmq

import (
	"fmt"
	"log"
	"os"
	"time"

	amqp "github.com/rabbitmq/amqp091-go"
)

var (
	Conn    *amqp.Connection
	Channel *amqp.Channel
)

func ConnectRabbitMQ() error {
	host := os.Getenv("RABBITMQ_HOST")
	port := os.Getenv("RABBITMQ_PORT")
	user := os.Getenv("RABBITMQ_USER")
	pass := os.Getenv("RABBITMQ_PASSWORD")

	url := fmt.Sprintf(
		"amqp://%s:%s@%s:%s/",
		user, pass, host, port,
	)

	var err error

	// 🔁 Retry loop
	for i := 1; i <= 10; i++ {
		log.Printf("Attempting RabbitMQ connection (%d/10)...", i)

		Conn, err = amqp.Dial(url)
		if err == nil {
			Channel, err = Conn.Channel()
			if err != nil {
				return err
			}

			log.Println("✅ RabbitMQ connected successfully")
			return nil
		}

		log.Printf("RabbitMQ not ready yet, retrying in 3 seconds...")
		time.Sleep(3 * time.Second)
	}

	return fmt.Errorf("❌ could not connect to RabbitMQ after multiple attempts")
}
