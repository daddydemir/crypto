package model

import (
	"github.com/google/uuid"
	"time"
)

type WeeklyModel struct {
	Id         uuid.UUID
	Min        float32
	Max        float32
	Avg        float32
	MinTime    time.Time
	MaxTime    time.Time
	FirstPrice float32
	LastPrice  float32
	FirstTime  time.Time
	LastTime   time.Time
	TimeHeader string
	Rate       float32
	Modulus    float32
	ExchangeId string
}
