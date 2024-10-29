package service

import (
	"fmt"
	"github.com/daddydemir/crypto/pkg/broker"
	"github.com/daddydemir/crypto/pkg/cache"
	"github.com/daddydemir/crypto/pkg/graphs/bollingerBands"
	"github.com/daddydemir/crypto/pkg/remote/coincap"
	"log/slog"
)

type bollingerService struct {
	cache.Cache
	broker.Broker
}

func NewBollingerService() *bollingerService {
	return &bollingerService{
		cache.GetCacheService(),
		broker.GetBrokerService(),
	}
}

func (b *bollingerService) CheckThresholds() {

	coins := coincap.ListCoins()
	if len(coins) == 0 {
		slog.Error("coincap.ListCoins", "error", "List is empty")
		return
	}

	for index, coin := range coins {
		if index > 30 {
			continue
		}
		bands := bollingerBands.NewBollingerBands(coin.Id, 20)
		_ = bands.Calculate()

		if bollingerBands.UpperBand == nil || len(bollingerBands.UpperBand) == 0 {
			slog.Error("BollingerBands:Calculate", "error", "this coin may not be cached either")
			continue
		}

		upperPrice := bollingerBands.UpperBand[len(bollingerBands.UpperBand)-1].Value
		lowerPrice := bollingerBands.LowerBand[len(bollingerBands.LowerBand)-1].Value

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
		bollingerBands.LowerBand = nil
		bollingerBands.UpperBand = nil
	}
}
