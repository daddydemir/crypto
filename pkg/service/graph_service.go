package service

import (
	"fmt"
	"github.com/daddydemir/crypto/pkg/broker"
	"github.com/daddydemir/crypto/pkg/database/service"
	"github.com/daddydemir/crypto/pkg/model"
	"log/slog"
)

var myMap map[string][]model.DailyModel

func RSIGraph(broker broker.Broker) {
	models := service.GetDailyForGraph()
	prepareData(models)
	calculateRsiIndex(broker)

}

func prepareData(models []model.DailyModel) {

	myMap = make(map[string][]model.DailyModel)

	for _, item := range models {

		if myMap[item.ExchangeId] != nil {
			currentList := myMap[item.ExchangeId]
			currentList = append(currentList, item)
			myMap[item.ExchangeId] = currentList
		} else {
			currentList := make([]model.DailyModel, 0)
			currentList = append(currentList, item)
			myMap[item.ExchangeId] = currentList
		}
	}
}

func calculateRsiIndex(broker broker.Broker) {

	for key, value := range myMap {

		var (
			positiveSum,
			negativeSum,
			gain,
			loss float32
		)

		for i, item := range value {

			if i == 0 {
				continue
			}

			data := item.Avg - value[i-1].Avg

			if data > 0 {
				gain += data
			} else {
				loss += data
			}
		}
		positiveSum = gain / 14
		negativeSum = loss / 14 * -1
		rsiCount := 100 - (100 / (1 + (positiveSum / negativeSum)))
		slog.Info("service.calculateRsiIndex", "positiveSum", positiveSum, "negativeSum", negativeSum)
		slog.Info("service.calculateRsiIndex", "Coin", key, "rsiCount", rsiCount)

		if rsiCount <= 30 || rsiCount >= 70 {
			slog.Info("service.calculateRsiIndex -> important", "Coin", key, "rsiCount", rsiCount)
			err := broker.SendMessage(fmt.Sprintf("(14) coin: %v rsi: %v ", key, rsiCount))
			if err != nil {
				slog.Error("calculateRsiIndex:broker.SendMessage", "err", err)
			}
		}
	}
}
