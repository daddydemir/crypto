package broker

import "github.com/daddydemir/crypto/pkg/broker/rabbit"

type Broker interface {
	SendMessage(message string) error
}

func GetBrokerService() Broker {
	return &rabbit.Publisher{}
}
