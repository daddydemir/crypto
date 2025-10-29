package indicator

type PriceRepository interface {
	GetLastNDaysPrices(ids []string, days int) (map[string][]float64, error)
	GetTopCoinIDs() ([]string, error)
}
