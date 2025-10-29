package coin

import (
	"github.com/daddydemir/crypto/pkg/domain/indicator"
	"strings"
)

type GetTopCoinsRSI struct {
	priceRepo indicator.PriceRepository
	rsiCalc   *indicator.RSI
}

func NewGetTopCoinsRSI(priceRepo indicator.PriceRepository) *GetTopCoinsRSI {
	return &GetTopCoinsRSI{
		priceRepo: priceRepo,
		rsiCalc:   indicator.NewRSI(),
	}
}

func (u *GetTopCoinsRSI) Execute() ([]RSIDTO, error) {
	coins, err := u.priceRepo.GetTopCoinIDs()
	if err != nil {
		return nil, err
	}
	coinIDs := make([]string, 0, len(coins))
	for _, c := range coins {
		coinIDs = append(coinIDs, strings.ToLower(c.Symbol))
	}
	prices, err := u.priceRepo.GetLastNDaysPrices(coinIDs, 14)
	if err != nil {
		return nil, err
	}

	var dtos []RSIDTO
	for _, c := range coins {
		id := strings.ToLower(c.Symbol)
		rsi := u.rsiCalc.Calculate(prices[id])
		dtos = append(dtos, RSIDTO{
			CoinID: id,
			Name:   c.Name,
			RSI:    rsi,
		})
	}

	return dtos, nil
}
