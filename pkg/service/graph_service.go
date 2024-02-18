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
	println(len(models))
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

	var (
		lastIndex     float32
		positiveSum   float32
		negativeSum   float32
		positiveCount int
		negativeCount int
	)

	for key, value := range myMap {

		log.Infoln("starting for : ", key)
		for i, item := range value {

			if i == 0 {
				positiveSum += item.Avg
				positiveCount++
			} else {
				if lastIndex > item.Avg {
					negativeSum += item.Avg
					negativeCount++
				} else {
					positiveSum += item.Avg
					positiveCount++
				}
			}
			lastIndex = item.Avg
		}
		positiveSum = positiveSum / float32(positiveCount)
		negativeSum = negativeSum / float32(negativeCount)
		rsiCount := 100 - (100 / (1 + (positiveSum / negativeSum)))
		log.Infoln("positive count: %v \nnegative count: %v \npositive avg : %v \nnegative avg : %v \n", positiveCount, negativeCount, positiveSum, negativeSum)
		log.Infoln("RSI : %.4f", rsiCount)

		if rsiCount <= 30 || rsiCount >= 70 {
			log.Infoln("RSI Index: %v", rsiCount)
			rabbitmq.SendQueue(fmt.Sprintf("(14) coin : %v rsi: %v", key, rsiCount))
		}
	}
}
