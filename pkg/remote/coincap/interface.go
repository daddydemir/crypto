package coincap

type CoinCapClient interface {
	ListCoins() (error, []Coin)
	HistoryWithId(s string) (error, []History)
	HistoryWithTime(s string, start, end int64) (error, []History)
}
