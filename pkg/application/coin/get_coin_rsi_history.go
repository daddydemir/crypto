package coin

import (
	"github.com/daddydemir/crypto/pkg/domain/indicator"
)

type GetCoinRSIHistory struct {
	priceRepo indicator.PriceRepository
	rsiCalc   *indicator.RSI
}

func NewGetCoinRSIHistory(priceRepo indicator.PriceRepository) *GetCoinRSIHistory {
	return &GetCoinRSIHistory{
		priceRepo: priceRepo,
		rsiCalc:   indicator.NewRSI(),
	}
}

func (u *GetCoinRSIHistory) Execute(coinID string, days int) ([]RSIHistoryDTO, error) {
	priceData, err := u.priceRepo.GetHistoricalPrices(coinID, days)
	if err != nil {
		return nil, err
	}

	if len(priceData) < 15 {
		return nil, nil
	}

	var results []RSIHistoryDTO
	for i := 14; i < len(priceData); i++ {
		var window []float64
		for j := i - 14; j <= i; j++ {
			window = append(window, priceData[j].Price)
		}

		rsi := u.rsiCalc.Calculate(window)
		results = append(results, RSIHistoryDTO{
			Date: priceData[i].Date.Format("2006-01-02"),
			RSI:  rsi,
		})
	}

	return results, nil
}
