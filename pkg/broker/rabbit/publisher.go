package rabbit

import (
	"context"
	"fmt"
	"github.com/daddydemir/crypto/config"
	"github.com/daddydemir/crypto/config/broker"
	amqp "github.com/rabbitmq/amqp091-go"
)

type Publisher struct{}

func (r *Publisher) SendMessage(message string) error {

	channel := broker.GetChannel()
	ctx := context.Background()

	err := channel.PublishWithContext(ctx,
		"",
		getQueue(config.Get("QUEUE_NAME")).Name,
		false,
		false,
		amqp.Publishing{
			ContentType: "text/plain",
			Body:        []byte(message),
		},
	)

	return err
}

// todo burada yanlislik olabilir | her defasinda bunu declare etmeye gerek olmayabilir...
func getQueue(queueName string) amqp.Queue {
	channel := broker.GetChannel()

	queue, err := channel.QueueDeclare(
		queueName,
		false,
		false,
		false,
		false,
		nil,
	)

	if err != nil {
		fmt.Println("Failed to declare a RabbitMQ queue", err)
	}

	return queue
}
