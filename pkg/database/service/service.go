package service

import (
	"github.com/daddydemir/crypto/config/database"
	"github.com/daddydemir/crypto/config/log"
	"github.com/daddydemir/crypto/pkg/adapter"
	"github.com/daddydemir/crypto/pkg/coingecko"
	"github.com/daddydemir/crypto/pkg/dao"
	"github.com/daddydemir/crypto/pkg/model"
)

func GetDailyFromDatabase() []model.DailyModel {
	var dailies []model.DailyModel
	start, end := getToday()
	database.D.Where("date between ? and ?", start, end).Find(&dailies)
	return dailies
}

func CreateDaily(morning bool) {
	var adapts []adapter.Adapter
	adapts = coingecko.GetTopHundred()
	var dailies []model.DailyModel
	for i := 0; i < len(adapts); i++ {
		dailies = append(dailies, adapts[i].AdapterToDaily(morning))
	}
	if !morning {
		dailyFromDb := GetDailyFromDatabase()
		sortSlice(dailyFromDb)
		sortSlice(dailies)

		for i := 0; i < len(dailies); i++ {
			if dailies[i].ExchangeId == dailyFromDb[i].ExchangeId {
				dailyFromDb[i].LastPrice = dailies[i].LastPrice

				//Min
				if dailyFromDb[i].Min > dailies[i].Min {
					dailyFromDb[i].Min = dailies[i].Min
				}

				//Max
				if dailyFromDb[i].Max < dailies[i].Max {
					dailyFromDb[i].Max = dailies[i].Max
				}

				dailyFromDb[i].Avg = (dailyFromDb[i].Min + dailyFromDb[i].Max) / 2
				dailyFromDb[i].Modulus = dailyFromDb[i].Max - dailyFromDb[i].Min
				dailyFromDb[i].Rate = dailyFromDb[i].Modulus * 100 / dailyFromDb[i].Avg
			} else {
				log.Infoln("::CreateDaily::false err:{}")
			}

		}
		database.D.Save(&dailyFromDb)
		return
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

func GetDailyFromDb(date dao.Date) []model.DailyModel {
	var dailies []model.DailyModel
	database.D.Where("date between ? and ?", date.StartDate, date.EndDate).Find(&dailies)
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

func CreateExchange() {
	var adapts []adapter.Adapter
	adapts = coingecko.GetTopHundred()
	var exchanges []model.ExchangeModel
	for i := 0; i < len(adapts); i++ {
		exchanges = append(exchanges, adapts[i].AdapterToExchange())
	}
	database.D.Save(&exchanges)
}

func GetExchangeFromDb() []model.ExchangeModel {
	var exchanges []model.ExchangeModel
	database.D.Find(&exchanges)
	return exchanges
}

func CreateWeekly() {
	var dailies []model.DailyModel
	weekStart, weekEnd := getWeek()

	database.D.Where("date between ? and ? ", weekStart, weekEnd).Find(&dailies)
	// todo
}
