package app

import "github.com/daddydemir/crypto/pkg/channels/donchian/domain"

type App struct {
	repo domain.Repository
}

func NewApp(repo domain.Repository) *App {
	return &App{repo: repo}
}

func (d *App) Series(symbol string) ([]domain.DonchianChannel, error) {

	datas, err := d.repo.GetRawDataWithSymbol(symbol)
	if err != nil {
		return nil, err
	}

	channels, err := domain.CalculateDonchian(datas)
	return channels, err
}
