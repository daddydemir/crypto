package domain

type Repository interface {
	GetPrices(coinID string) ([]PriceData, error)
}
