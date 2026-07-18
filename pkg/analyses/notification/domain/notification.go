package domain

type Notification struct {
	Type       string
	Coin       string
	CreateTime int64
	Image      string
}

type Notifications interface {
	GetAll() []Notification
}
