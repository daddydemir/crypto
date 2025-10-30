package indicator

import (
	"github.com/daddydemir/crypto/pkg/remote/coincap"
)

type PriceRepository interface {
	GetLastNDaysPrices(ids []string, days int) (map[string][]float64, error)
	GetTopCoinIDs() ([]coincap.Coin, error)
	GetHistoricalPrices(coinID string, days int) ([]PriceData, error)
}
