package model

import (
	"github.com/daddydemir/crypto/pkg/adapter"
	"github.com/google/uuid"
	"time"
)

type DailyModel struct {
	Id         uuid.UUID
	Min        float32
	Max        float32
	Avg        float32
	FirstPrice float32
	LastPrice  float32
	Date       time.Time
	Rate       float32
	Modulus    float32
	ExchangeId string
}

func (d DailyModel) ToMap() map[string]DailyModel {
	m := make(map[string]DailyModel)
	m[d.ExchangeId] = d
	return m
}

func AdapterToDaily(a adapter.Adapter, morning bool) DailyModel {
	daily := DailyModel{
		Id:         uuid.New(),
		Min:        a.Low24H,
		Max:        a.High24H,
		Avg:        (a.Low24H + a.High24H) / 2,
		Date:       time.Now(),
		Rate:       ((a.High24H - a.Low24H) * 100) / ((a.Low24H + a.High24H) / 2),
		Modulus:    a.High24H - a.Low24H,
		ExchangeId: a.Symbol,
	}

	if morning {
		daily.FirstPrice = a.CurrentPrice
	} else {
		daily.LastPrice = a.CurrentPrice
	}

	return daily
}
