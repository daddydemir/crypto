package domain

import "time"

type Alert struct {
	ID         uint
	Coin       string
	Price      float32
	IsAbove    bool
	CreateDate time.Time
	IsActive   bool
}

func NewAlert(coin string, price float32, isAbove bool) Alert {
	return Alert{
		Coin:       coin,
		Price:      price,
		IsAbove:    isAbove,
		CreateDate: time.Now(),
		IsActive:   true,
	}
}

func (a *Alert) Deactivate() {
	a.IsActive = false
}

func (a *Alert) Update(price float32, isAbove bool) {
	a.Price = price
	a.IsAbove = isAbove
}
