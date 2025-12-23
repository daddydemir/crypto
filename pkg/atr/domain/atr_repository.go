package domain

type Repository interface {
	GetPointsBySymbol(symbol string) ([]AtrPoint, error)
}
