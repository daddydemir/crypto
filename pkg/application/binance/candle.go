package binance

import "github.com/daddydemir/crypto/pkg/domain/binance"

type CandleService struct {
	repo   binance.CandleRepository
	source binance.CandleDataSource
}

func NewCandleService(repo binance.CandleRepository, source binance.CandleDataSource) *CandleService {
	return &CandleService{
		repo:   repo,
		source: source,
	}
}

func (cs *CandleService) FetchAndStore(symbol, interval string, start, end int64, limit int) error {
	candles, err := cs.source.Fetch(symbol, interval, start, end, limit)
	if err != nil {
		return err
	}
	return cs.repo.SaveMany(candles)
}

func (cs *CandleService) Fetch(symbol, interval string, start, end int64, limit int) ([]binance.Candle, error) {
	return cs.source.Fetch(symbol, interval, start, end, limit)
}
