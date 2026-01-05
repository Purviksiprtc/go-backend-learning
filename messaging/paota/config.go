package paota

import (
	"os"

	"github.com/surendratiwari3/paota/config"
)

func LoadPaotaConfig(queueName, bindingKey string) config.Config {
	return config.Config{
		Broker:        "amqp",
		TaskQueueName: queueName,
		AMQP: &config.AMQPConfig{
			Url:                os.Getenv("RABBITMQ_URL"),
			Exchange:           os.Getenv("RABBITMQ_EXCHANGE"),
			ExchangeType:       "topic",
			BindingKey:         bindingKey,
			PrefetchCount:      10,
			ConnectionPoolSize: 5,

			// ✅ DLQ ONLY FOR MAIN QUEUES
			FailedQueue: os.Getenv("RABBITMQ_DLQ"),
		},
	}
}
