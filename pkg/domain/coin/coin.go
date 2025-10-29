package coin

import "time"

type Coin struct {
	ID       string
	Name     string
	Symbol   string
	PriceUSD float64
}

type ChangeStats struct {
	CoinID    string
	Current   float64
	Change24h float64
	Change7d  float64
	UpdatedAt time.Time
}
