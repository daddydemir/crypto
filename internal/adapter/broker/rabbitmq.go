package brokeradapter

import (
	"context"
	amqp "github.com/rabbitmq/amqp091-go"
)

type RabbitPublisher struct {
	channel *amqp.Channel
	queue   amqp.Queue
}

func NewRabbitPublisher(channel *amqp.Channel, queueName string) (*RabbitPublisher, error) {
	queue, err := channel.QueueDeclare(
		queueName,
		false, false, false, false,
		nil,
	)
	if err != nil {
		return nil, err
	}

	return &RabbitPublisher{
		channel: channel,
		queue:   queue,
	}, nil
}

func (p *RabbitPublisher) SendMessage(msg string) error {
	ctx := context.Background()

	return p.channel.PublishWithContext(ctx,
		"",
		p.queue.Name,
		false,
		false,
		amqp.Publishing{
			ContentType: "text/plain",
			Body:        []byte(msg),
		},
	)
}
