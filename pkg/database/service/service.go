package service

import (
	"github.com/daddydemir/crypto/config/database"
	"github.com/daddydemir/crypto/pkg/adapter"
	"github.com/daddydemir/crypto/pkg/coingecko"
	"github.com/daddydemir/crypto/pkg/model"
)

func CreateDaily(morning bool) {
	var adapts []adapter.Adapter
	adapts = coingecko.GetTopHundred()
	var dailies []model.DailyModel
	for i := 0; i < len(adapts); i++ {
		dailies = append(dailies, adapts[i].AdapterToDaily(morning))
	}
	database.D.Save(&dailies)
}

func GetDaily() []model.DailyModel {
	var adapts []adapter.Adapter
	adapts = coingecko.GetTopHundred()
	var dailies []model.DailyModel
	for i := 0; i < len(adapts); i++ {
		dailies = append(dailies, adapts[i].AdapterToDaily(true))
	}
	return dailies
}

func GetExchange() []model.ExchangeModel {
	var adapts []adapter.Adapter
	adapts = coingecko.GetTopHundred()
	var exchanges []model.ExchangeModel
	for i := 0; i < len(adapts); i++ {
		exchanges = append(exchanges, adapts[i].AdapterToExchange())
	}
	return exchanges
}

func CreateWeekly() {
	var dailies []model.DailyModel
	weekStart, weekEnd := getWeek()

	database.D.Where("date between ? and ? ", weekStart, weekEnd).Find(&dailies)
	// todo
}
