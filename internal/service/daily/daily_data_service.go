package daily

import (
	"log/slog"

	"github.com/daddydemir/crypto/internal/domain/model"
	"github.com/daddydemir/crypto/internal/port"
)

type DailyDataService struct {
	client port.CoingeckoClient
}

func NewDailyDataService(client port.CoingeckoClient) *DailyDataService {
	return &DailyDataService{client: client}
}

func (d *DailyDataService) GetDaily(morning bool) []model.DailyModel {
	dailies, err := d.client.GetTopHundredDaily()
	if err != nil {
		slog.Error("DailyDataService:GetTopHundred", "error", err)
		return nil
	}

	//	var dailies []model.DailyModel
	//	for _, a := range dailies {
	//		dailies = append(dailies, model.AdapterToDaily(a, morning))
	//	}

	return dailies
}
