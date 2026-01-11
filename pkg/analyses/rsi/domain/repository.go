package domain

type Repository interface {
	GetTopCoins(limit int) ([]Coin, error)
	GetPriceAt(coinID string, date int) (float64, error)
	GetCurrentPrices() ([]Coin, error)
	GetPriceChanges() ([]PriceResult, error)
}

type HistoryRepository interface {
	GetPriceAt(coinID string, date int) (float64, error)
}

type MarketRepository interface {
	GetCurrentPrices() ([]Coin, error)
	GetPriceChanges() ([]PriceResult, error)
}
