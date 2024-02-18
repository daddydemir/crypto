package service

import (
	"github.com/daddydemir/crypto/pkg/model"
	"testing"
)

func TestPrepareData(t *testing.T) {
	list := make([]model.DailyModel, 0)
	m1 := model.DailyModel{ExchangeId: "btc", Avg: 2006}
	m2 := model.DailyModel{ExchangeId: "btc", Avg: 1990}
	m3 := model.DailyModel{ExchangeId: "btc", Avg: 1942}
	m4 := model.DailyModel{ExchangeId: "btc", Avg: 1949}
	m5 := model.DailyModel{ExchangeId: "btc", Avg: 1953}
	m6 := model.DailyModel{ExchangeId: "btc", Avg: 1942}
	m7 := model.DailyModel{ExchangeId: "btc", Avg: 1973}
	m8 := model.DailyModel{ExchangeId: "btc", Avg: 1996}
	m9 := model.DailyModel{ExchangeId: "btc", Avg: 2011}
	m10 := model.DailyModel{ExchangeId: "btc", Avg: 2065}
	m11 := model.DailyModel{ExchangeId: "btc", Avg: 2049}
	m12 := model.DailyModel{ExchangeId: "btc", Avg: 2093}
	m13 := model.DailyModel{ExchangeId: "btc", Avg: 2088}
	m14 := model.DailyModel{ExchangeId: "btc", Avg: 2085}

	list = append(list, m1)
	list = append(list, m2)
	list = append(list, m3)
	list = append(list, m4)
	list = append(list, m5)
	list = append(list, m6)
	list = append(list, m7)
	list = append(list, m8)
	list = append(list, m9)
	list = append(list, m10)
	list = append(list, m11)
	list = append(list, m12)
	list = append(list, m13)
	list = append(list, m14)

	prepareData(list)

	calculateRsiIndex()
}
