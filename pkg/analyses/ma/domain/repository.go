package domain

type PriceHistoryRepository interface {
	GetPrices(coinID string) ([]PriceData, error)
}
