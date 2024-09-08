package broker

import (
	"github.com/daddydemir/crypto/config"
	"github.com/segmentio/kafka-go"
)

var writer *kafka.Writer

func StartKafkaConnection() {
	GetWriter()
}

func GetWriter() *kafka.Writer {

	if writer == nil {
		writer = &kafka.Writer{
			Addr:  kafka.TCP(config.Get("KAFKA_URL")),
			Topic: config.Get("KAFKA_TOPIC"),
		}
	}

	return writer
}
