package alert

import "time"

type Alert struct {
	ID         uint
	Coin       string
	Price      float32
	IsAbove    bool
	CreateDate time.Time
	IsActive   bool
}

// Factory fonksiyonu
func NewAlert(coin string, price float32, isAbove bool) Alert {
	return Alert{
		Coin:       coin,
		Price:      price,
		IsAbove:    isAbove,
		CreateDate: time.Now(),
		IsActive:   true,
	}
}

// Davranışlar (iş kuralları)
func (a *Alert) Deactivate() {
	a.IsActive = false
}

func (a *Alert) Update(price float32, isAbove bool) {
	a.Price = price
	a.IsAbove = isAbove
}
