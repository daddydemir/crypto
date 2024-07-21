package service

import (
	"fmt"
	"github.com/daddydemir/crypto/config/log"
	"github.com/daddydemir/crypto/pkg/database/service"
	"github.com/daddydemir/crypto/pkg/model"
	"github.com/daddydemir/crypto/pkg/rabbitmq"
)

var myMap map[string][]model.DailyModel

func RSIGraph() {
	models := service.GetDailyForGraph()
	prepareData(models)
	calculateRsiIndex()

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

func calculateRsiIndex() {

	for key, value := range myMap {

		var (
			positiveSum,
			negativeSum,
			gain,
			loss float32
		)

		log.Infoln("starting for : ", key)
		for i, item := range value {

			if i == 0 {
				continue
			}

			data := item.Avg - value[i-1].Avg

			//fmt.Printf("index: %v minus: %v \n", i, data)

			if data > 0 {
				gain += data
			} else {
				loss += data
			}
		}
		positiveSum = gain / 14
		negativeSum = loss / 14 * -1
		rsiCount := 100 - (100 / (1 + (positiveSum / negativeSum)))
		log.Infoln("positive avg : %v \nnegative avg : %v \n", positiveSum, negativeSum)
		log.Infoln("RSI : %.4f", rsiCount)

		if rsiCount <= 30 || rsiCount >= 70 {
			log.Infoln("RSI Index: %v", rsiCount)
			rabbitmq.SendQueue(fmt.Sprintf("(14) coin : %v rsi: %v", key, rsiCount))
		}
	}
}
