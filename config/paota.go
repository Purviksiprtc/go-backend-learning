package config

import (
	"os"

	"github.com/surendratiwari3/paota/config"
)

func LoadPaota(queue, bindingKey string) config.Config {
	return config.Config{
		Broker:        "amqp",
		TaskQueueName: queue,
		AMQP: &config.AMQPConfig{
			Url:                os.Getenv("RABBITMQ_URL"),
			Exchange:           os.Getenv("RABBITMQ_EXCHANGE"),
			ExchangeType:       "topic",
			BindingKey:         bindingKey,
			PrefetchCount:      10,
			ConnectionPoolSize: 5,
			FailedQueue:        os.Getenv("RABBITMQ_DLQ"),
		},
	}
}
