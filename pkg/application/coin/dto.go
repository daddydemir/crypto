package coin

type StatsDTO struct {
	ID        string  `json:"id"`
	Name      string  `json:"name"`
	Symbol    string  `json:"symbol"`
	Price     float64 `json:"price"`
	Change24h float64 `json:"change24h"`
	Change7d  float64 `json:"change7d"`
}
