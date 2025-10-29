package coin

import "github.com/daddydemir/crypto/pkg/domain/indicator"

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
	coinIDs, err := u.priceRepo.GetTopCoinIDs()
	if err != nil {
		return nil, err
	}
	prices, err := u.priceRepo.GetLastNDaysPrices(coinIDs, 14)
	if err != nil {
		return nil, err
	}

	var dtos []RSIDTO
	for _, id := range coinIDs {
		rsi := u.rsiCalc.Calculate(prices[id])
		dtos = append(dtos, RSIDTO{
			CoinID: id,
			RSI:    rsi,
		})
	}

	return dtos, nil
}
