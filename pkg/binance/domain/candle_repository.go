package domain

type Repository interface {
	GetBySymbol(symbol string) ([]Candle, error)
}
