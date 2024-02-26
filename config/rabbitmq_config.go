package config

import (
	"github.com/daddydemir/crypto/config/log"
	amqp "github.com/rabbitmq/amqp091-go"
)

var Channel *amqp.Channel

func NewRabbitMQ() {

	conn, err := amqp.Dial(Get("RABBIT_MQ_URL"))

	if err != nil {
		log.LOG.Fatal("RabbitMQ connection was failed : ", err)
	}

	log.LOG.Println("Connecting to RabbitMQ")

	ch, err := conn.Channel()
	if err != nil {
		log.LOG.Errorln(err)
	}

	Channel = ch
}
