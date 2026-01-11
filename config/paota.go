package config

import (
	"os"

	paotaConfig "github.com/surendratiwari3/paota/config"
)

func LoadPaota(queueName, bindingKey string) paotaConfig.Config {
	return paotaConfig.Config{
		Broker:        "amqp",
		TaskQueueName: queueName,
		AMQP: &paotaConfig.AMQPConfig{
			Url:                os.Getenv("RABBITMQ_URL"),
			Exchange:           os.Getenv("RABBITMQ_EXCHANGE"),
			ExchangeType:       "topic",
			BindingKey:         bindingKey,
			PrefetchCount:      10,
			ConnectionPoolSize: 5,
		},
	}
}
