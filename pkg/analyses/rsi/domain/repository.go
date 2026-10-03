package domain

type Coin struct {
	Symbol string
	Name   string
}

type Repository interface {
	GetLastNDaysPrices(ids []string, days int) (map[string][]float64, error)
	GetTopCoinIDs() ([]Coin, error)
}
