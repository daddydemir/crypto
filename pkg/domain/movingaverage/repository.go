package movingaverage

import "context"

type PriceHistoryRepository interface {
	GetLastNDaysPricesWithDates(ctx context.Context, coinID string, days int) ([]PriceData, error)
}
