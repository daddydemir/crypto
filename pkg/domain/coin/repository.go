package coin

type Repository interface {
	GetTopCoins(limit int) ([]Coin, error)
}

type HistoryRepository interface {
	GetPriceAt(coinID string, date int) (float64, error)
}

type MarketRepository interface {
	GetCurrentPrices() ([]Coin, error)
}
