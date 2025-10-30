package coin

type StatsDTO struct {
	ID        string  `json:"id"`
	Name      string  `json:"name"`
	Symbol    string  `json:"symbol"`
	Price     float64 `json:"price"`
	Change24h float64 `json:"change24h"`
	Change7d  float64 `json:"change7d"`
}

type RSIDTO struct {
	Name   string  `json:"name"`
	CoinID string  `json:"coin_id"`
	RSI    float64 `json:"rsi"`
}

type RSIHistoryDTO struct {
	RSI  float64 `json:"rsi"`
	Date string  `json:"date"`
}
