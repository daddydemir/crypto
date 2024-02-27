package config

import (
	amqp "github.com/rabbitmq/amqp091-go"
	"log"
)

var Channel *amqp.Channel

func NewRabbitMQ() {

	conn, err := amqp.Dial(Get("RABBIT_MQ_URL"))

	if err != nil {
		log.Println("RabbitMQ connection was failed : ", err)
	}

	log.Println("Connecting to RabbitMQ")

	ch, err := conn.Channel()
	if err != nil {
		log.Println(err)
	}

	Channel = ch
}
