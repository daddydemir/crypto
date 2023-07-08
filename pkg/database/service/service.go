package service

import (
	"fmt"
	"github.com/daddydemir/crypto/config/database"
	"github.com/daddydemir/crypto/config/log"
	"github.com/daddydemir/crypto/pkg/adapter"
	"github.com/daddydemir/crypto/pkg/coingecko"
	"github.com/daddydemir/crypto/pkg/dao"
	"github.com/daddydemir/crypto/pkg/model"
	"github.com/daddydemir/crypto/pkg/telegram"
)

func GetDailyFromDatabase() []model.DailyModel {
	var dailies []model.DailyModel
	start, end := getToday()
	database.D.Where("date between ? and ?", start, end).Find(&dailies)
	return dailies
}

func Temp() {}

func MergeMap(m1, m2 map[string]model.DailyModel) map[string]model.DailyModel {
	merged := make(map[string]model.DailyModel, len(m1)+len(m2))
	merged = m1
	for key, value := range m2 {
		merged[key] = value
	}

	return merged
}

func CreateDaily(morning bool) {
	var adapts []adapter.Adapter
	adapts = coingecko.GetTopHundred()
	var dailies []model.DailyModel
	for i := 0; i < len(adapts); i++ {
		dailies = append(dailies, adapts[i].AdapterToDaily(morning))
	}
	if !morning {
		db := make(map[string]model.DailyModel)
		gecko := make(map[string]model.DailyModel)
		dailyFromDb := GetDailyFromDatabase()
		save := []model.DailyModel{}

		for i := 0; i < len(dailies); i++ {
			gecko = MergeMap(gecko, dailies[i].ToMap())
			db = MergeMap(db, dailyFromDb[i].ToMap())

		}

		current := model.DailyModel{}

		for key, value := range gecko {
			current = db[key]
			if current.ExchangeId == "" {
				continue
			}
			current.LastPrice = value.LastPrice
			if current.Min > value.Min {
				current.Min = value.Min
			}

			if current.Max < value.Max {
				current.Max = value.Max
			}

			current.Avg = (current.Min + current.Max) / 2
			current.Modulus = current.Max - current.Min
			current.Rate = current.Modulus * 100 / current.Avg

			save = append(save, current)
		}

		result := database.D.Save(&save)
		if result.Error != nil {
			log.Errorln(result.Error)
		}
		CreateMessage()

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

func GetDailyWithId(date dao.Date) []model.DailyModel {
	var dailies []model.DailyModel
	database.D.Where("date between ? and ? and exchange_id = ?", date.StartDate, date.EndDate, date.Id).Find(&dailies)
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

func CreateMessage() (string, string) {
	var smaller []model.DailyModel
	var bigger []model.DailyModel
	var m1, m2, rate, mod string

	start, end := getToday()
	database.D.Where("date between ? and ? and avg > 1 order by rate desc limit 5", start, end).Find(&bigger)
	database.D.Where("date between ? and ? and avg < 1 order by rate desc limit 5", start, end).Find(&smaller)
	for i := 0; i < 5; i++ {
		rate = fmt.Sprintf("%.2f", bigger[i].Rate)
		mod = fmt.Sprintf("%.1f", bigger[i].Modulus)
		m1 += "(" + bigger[i].ExchangeId + ")\t %" + rate + "\t | \t" + mod + "$ \n"

		rate = fmt.Sprintf("%.2f", smaller[i].Rate)
		mod = fmt.Sprintf("%v", smaller[i].Modulus)
		m2 += "(" + smaller[i].ExchangeId + ")\t %" + rate + "\t | \t" + mod + "$ \n"
	}
	telegram.Pre(m1, m2)
	return m1, m2
}
