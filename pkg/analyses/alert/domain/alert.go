package domain

import "time"

type Alert struct {
	ID             uint
	Username       string `gorm:"size:25;index" json:"-"`
	Coin           string
	Price          float32
	IsAbove        bool
	CreateDate     time.Time
	IsActive       bool
	TriggeredAt    *time.Time
	TriggeredPrice float32
}

func NewAlert(username, coin string, price float32, isAbove bool) Alert {
	return Alert{
		Username:   username,
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
	a.IsActive = true
	a.TriggeredAt = nil
	a.TriggeredPrice = 0
}
