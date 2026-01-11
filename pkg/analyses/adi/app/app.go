package app

import (
	"github.com/daddydemir/crypto/pkg/analyses/adi/domain"
)

type App struct {
	repo domain.Repository
}

func NewApp(repo domain.Repository) *App {
	return &App{
		repo: repo,
	}
}

func (a *App) GetADISeries(coinID string) ([]domain.Point, error) {
	// Get price data
	priceData, err := a.repo.GetRawDataWithSymbol(coinID)
	if err != nil {
		return nil, err
	}
	//
	//if len(priceData) < 14+1 {
	//	return &domain.Series{
	//		CoinID: coinID,
	//		Points: []domain.Point{},
	//	}, nil
	//}

	// Calculate ADI points
	points := domain.CalculateADI(priceData, 14)
	return points, nil
}
