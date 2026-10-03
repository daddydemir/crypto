package broker

import (
	"fmt"
	"github.com/daddydemir/crypto/config"
	amqp "github.com/rabbitmq/amqp091-go"
	"strings"
)

var channel *amqp.Channel

func StartRabbitmqConnection() {
	GetChannel()
}

func GetChannel() *amqp.Channel {

	if channel == nil || channel.IsClosed() {
		rabbitURL := strings.TrimSpace(config.Get("RABBIT_MQ_URL"))
		if rabbitURL != "" && !strings.HasPrefix(rabbitURL, "amqp://") && !strings.HasPrefix(rabbitURL, "amqps://") {
			rabbitURL = "amqp://" + rabbitURL
		}
		dial, err := amqp.Dial(rabbitURL)
		if err != nil {
			fmt.Println("Failed to connect to RabbitMQ", err)
			panic(err)
		}

		channel, err = dial.Channel()
		if err != nil {
			fmt.Println("Failed to open a RabbitMQ channel", err)
		}
	}

	return channel
}
