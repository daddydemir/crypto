package domain

type Repository interface {
	GetRawDataWithSymbol(symbol string) ([]DonchianData, error)
}
