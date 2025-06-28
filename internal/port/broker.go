package port

type Broker interface {
	SendMessage(msg string) error
}
