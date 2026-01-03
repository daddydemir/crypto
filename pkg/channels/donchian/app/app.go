package app

import "github.com/daddydemir/crypto/pkg/channels/donchian/domain"

type DonchianApp struct {
	repo domain.Repository
}

func NewDonchianApp(repo domain.Repository) *DonchianApp {
	return &DonchianApp{repo: repo}
}

func (d *DonchianApp) Series(symbol string) ([]domain.DonchianChannel, error) {

	datas, err := d.repo.GetRawDataWithSymbol(symbol)
	if err != nil {
		return nil, err
	}

	channels, err := domain.CalculateDonchian(datas)
	return channels, err
}
