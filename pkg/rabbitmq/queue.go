package rabbitmq

import (
	"context"
	"github.com/daddydemir/crypto/config"
	amqp "github.com/rabbitmq/amqp091-go"
	"log"
	"time"
)

func getQueue(name string) amqp.Queue {
	q, err := config.Channel.QueueDeclare(
		name,
		false,
		false,
		false,
		false,
		nil,
	)
	if err != nil {
		log.Fatal("", err)
	}
	return q
}

func SendQueue(message string) {
	ctx, cancel := context.WithTimeout(context.Background(), 5*time.Second)
	defer cancel()

	err := config.Channel.PublishWithContext(ctx, "", getQueue(config.Get("QUEUE_NAME")).Name, false, false, amqp.Publishing{
		ContentType: "text/plain",
		Body:        []byte(message),
	})

	if err != nil {
		log.Println(err)
	}
}
