package broker

type Broker interface {
	SendMessage(message string) error
}
