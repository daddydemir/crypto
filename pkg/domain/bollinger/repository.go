package bollinger

type PriceHistoryRepository interface {
	GetLastNDaysPricesWithDates(coinID string, days int) ([]PriceData, error)
}
