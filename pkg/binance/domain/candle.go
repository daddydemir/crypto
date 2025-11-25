package domain

import "time"

type Candle struct {
	Symbol string
	Time   time.Time
	Close  float64
}
