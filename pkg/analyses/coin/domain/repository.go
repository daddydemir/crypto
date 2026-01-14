package domain

type Repository interface {
	GetCurrentPrices() ([]Coin, error)
	GetPriceChanges() ([]PriceResult, error)
}
