package domain

type Repository interface {
	GetBySymbol(symbol string) ([]Candle, error)
	GetBySymbolAndYear(symbol, year string) ([]Candle, error)
	GetBySymbolAndYearMonth(symbol, year, month string) ([]Candle, error)
}
