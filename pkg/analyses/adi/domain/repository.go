package domain

type Repository interface {
	GetRawDataWithSymbol(symbol string) ([]PriceData, error)
}
