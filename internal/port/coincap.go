package port

import "github.com/daddydemir/crypto/internal/domain/model"

type CoinCapAPI interface {
	ListCoins() (error, []model.Coin)
	HistoryWithId(id string) (error, []model.History)
	HistoryWithTime(id string, start, end int64) (error, []model.History)
}
