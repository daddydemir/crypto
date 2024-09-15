package ma

import (
	"github.com/daddydemir/crypto/pkg/cache"
	localCharts "github.com/daddydemir/crypto/pkg/charts"
	"github.com/daddydemir/crypto/pkg/graphs"
	"github.com/daddydemir/crypto/pkg/remote/coincap"
	"log/slog"
	"net/http"
)

type sma struct {
	name         string
	period       int
	cacheService cache.Cache
}

func NewSma(name string) *sma {
	return &sma{
		name:         name,
		period:       10,
		cacheService: cache.GetCacheService(),
	}
}

func (s *sma) Index() float32 {
	panic("invalid method")
	return 0
}

func (s *sma) Calculate() []graphs.ChartModel {
	list := make([]coincap.History, 0)
	err := s.cacheService.GetList(s.name, &list, 0, -1)

	if err != nil {
		slog.Error("Calculate:cacheService.GetList", "coin", s.name, "error", err)
		return nil
	}

	response := make([]graphs.ChartModel, 0, len(list)-s.period)

	for z := 0; z <= len(list)-s.period; z++ {
		response = append(response, graphs.ChartModel{
			Date:  list[z+s.period-1].Date,
			Value: s.createAvg(list[z : z+s.period]),
		})
	}

	return response
}

func (s *sma) Draw(list []graphs.ChartModel) func(w http.ResponseWriter, r *http.Request) {
	return func(w http.ResponseWriter, r *http.Request) {
		chart := localCharts.CreateLineChart("SMA (Simple Moving Average)")
		dates, data := localCharts.ChartModel2lineData(list)

		chart.SetXAxis(dates).AddSeries(s.name, data)

		err := chart.Render(w)
		if err != nil {
			slog.Error("Draw:chart.Render", "error", err)
		}
	}
}

func (s *sma) createAvg(list []coincap.History) float32 {

	var totalPrice float32
	for _, h := range list {
		totalPrice += h.PriceUsd
	}
	return totalPrice / float32(s.period)
}
