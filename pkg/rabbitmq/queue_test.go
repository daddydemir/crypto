package rabbitmq

import (
	"github.com/daddydemir/crypto/config"
	"testing"
)

func TestSendQueue(t *testing.T) {
	config.NewRabbitMQ()
	SendQueue("test")
}

