package domain

type Repository interface {
	GetLastNDaysPricesWithDates(coinID string, days int) ([]PriceData, error)
}
