package service

import (
	"fmt"
	"github.com/daddydemir/crypto/pkg/broker"
	"github.com/daddydemir/crypto/pkg/cache"
	"github.com/daddydemir/crypto/pkg/graphs/bollingerBands"
	"github.com/daddydemir/crypto/pkg/remote/coincap"
	"log/slog"
)

type BollingerService struct {
	cache.Cache
	broker.Broker
	client coincap.CoinCapClient
}

func NewBollingerService(cache2 cache.Cache, broker2 broker.Broker, client coincap.CoinCapClient) *BollingerService {
	return &BollingerService{
		cache2,
		broker2,
		client,
	}
}

func (b *BollingerService) CheckThresholds() {

	_, coins := b.client.ListCoins()
	if len(coins) == 0 {
		slog.Error("coincap.ListCoins", "error", "List is empty")
		return
	}

	for index, coin := range coins {
		if index > 30 {
			continue
		}
		bands := bollingerBands.NewBollingerBands(coin.Id, 20)
		middle := bands.Calculate()

		if bollingerBands.UpperBand == nil || len(bollingerBands.UpperBand) == 0 {
			slog.Error("BollingerBands:Calculate", "error", "this coin may not be cached either")
			continue
		}

		upperLast := len(bollingerBands.UpperBand) - 1
		lowerLast := len(bollingerBands.LowerBand) - 1

		upperPrice := bollingerBands.UpperBand[upperLast].Value
		lowerPrice := bollingerBands.LowerBand[lowerLast].Value

		if coin.PriceUsd > upperPrice {
			message := fmt.Sprintf("[Bollinger] \ncoin: %v price: %.3f long: %.3f \n", coin.Id, coin.PriceUsd, upperPrice)
			slog.Info("Bollinger", "message", message)
			_ = b.SendMessage(message)
		}

		if lowerPrice > coin.PriceUsd {
			message := fmt.Sprintf("[Bollinger] \ncoin: %v price: %.3f short: %.3f \n", coin.Id, coin.PriceUsd, lowerPrice)
			slog.Info("Bollinger", "message", message)
			_ = b.SendMessage(message)
		}

		if isPriceDownCrossing(coin.PriceUsd, bollingerBands.Original[len(bollingerBands.Original)-1].Value, middle[len(middle)-1].Value, middle[len(middle)-2].Value) {
			message := fmt.Sprintf("[Bollinger] \ncoin: %v fiyat dusuyor.", coin.Id)
			slog.Info("Bollinger", "message", message)
			_ = b.SendMessage(message)
		}

		bollingerBands.LowerBand = nil
		bollingerBands.UpperBand = nil
		bollingerBands.Original = nil
	}
}

func isPriceDownCrossing(todayPrice, yesterdayPrice, todayMiddle, yesterdayMiddle float32) bool {
	if todayMiddle > todayPrice && yesterdayPrice > yesterdayMiddle {
		return true
	}
	return false
}
