package kafka

import (
	"context"
	"github.com/daddydemir/crypto/config/broker"
	"github.com/segmentio/kafka-go"
)

type Publisher struct{}

func (p *Publisher) SendMessage(message string) error {

	writer := broker.GetWriter()
	ctx := context.Background()
	err := writer.WriteMessages(ctx, kafka.Message{
		Value: []byte(message),
	})

	return err
}
