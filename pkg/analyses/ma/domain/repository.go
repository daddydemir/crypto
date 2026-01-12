package domain

type PriceHistoryRepository interface {
	GetLastNDaysPricesWithDates(coinID string, days int) ([]PriceData, error)
}
