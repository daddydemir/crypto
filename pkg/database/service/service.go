package service

import (
	"fmt"
	"time"

	"github.com/daddydemir/crypto/config/database"
	"github.com/daddydemir/crypto/pkg/broker"
	"github.com/daddydemir/crypto/pkg/model"
)

func CreateMessage(broker broker.Broker) {
	var smaller []model.DailyModel
	var bigger []model.DailyModel
	var m1, m2, rate, mod string

	startDate := time.Now()
	endDate := time.Now().AddDate(0, 0, -1)
	database.D.Where("date between ? and ? and avg > 1 order by rate desc limit 5", endDate, startDate).Find(&bigger)
	database.D.Where("date between ? and ? and avg < 1 order by rate desc limit 5", endDate, startDate).Find(&smaller)
	for i := 0; i < 5; i++ {
		rate = fmt.Sprintf("%.2f", bigger[i].Rate)
		mod = fmt.Sprintf("%.1f", bigger[i].Modulus)
		m1 += "(" + bigger[i].ExchangeId + ")\t %" + rate + "\t | \t" + mod + "$ \n"

		rate = fmt.Sprintf("%.2f", smaller[i].Rate)
		mod = fmt.Sprintf("%v", smaller[i].Modulus)
		m2 += "(" + smaller[i].ExchangeId + ")\t %" + rate + "\t | \t" + mod + "$ \n"
	}
	broker.SendMessage(m1)
	broker.SendMessage(m2)
}
