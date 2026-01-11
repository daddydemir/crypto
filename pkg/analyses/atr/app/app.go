package app

import (
	"github.com/daddydemir/crypto/pkg/analyses/atr/domain"
)

type App struct {
	repository domain.Repository
}

func NewApp(repository domain.Repository) *App {
	return &App{
		repository: repository,
	}
}

func (a *App) GetPoints(symbol string) ([]domain.Point, error) {
	atrPoints, err := a.repository.GetPointsBySymbol(symbol)
	if err != nil {
		return nil, err
	}
	points, err := domain.CalculateATR(atrPoints)
	return points, err
}
