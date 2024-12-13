package model

import (
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

type DailyRepository interface {
	FindByDateRange(start, end string) ([]DailyModel, error)
	FindByIdAndDateRange(id, start, end string) ([]DailyModel, error)
	FindTopSmallerByRate(start, end string) ([5]DailyModel, error)
	FindTopBiggerByRate(start, end string) ([5]DailyModel, error)
	SaveAll(dailies []DailyModel) error
}
