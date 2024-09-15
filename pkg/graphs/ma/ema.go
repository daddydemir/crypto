package ma

import (
	"fmt"
	"github.com/daddydemir/crypto/pkg/cache"
	localCharts "github.com/daddydemir/crypto/pkg/charts"
	"github.com/daddydemir/crypto/pkg/graphs"
	"github.com/daddydemir/crypto/pkg/remote/coincap"
	"log/slog"
	"net/http"
)

type ema struct {
	name         string
	period       int
	smoothing    float32
	cacheService cache.Cache
}

func NewEma(name string, period int) *ema {
	return &ema{
		name:         name,
		period:       period,
		smoothing:    2 / float32(1+period),
		cacheService: cache.GetCacheService(),
	}
}

func (e *ema) Index() float32 {
	panic("invalid method")
}

func (e *ema) Calculate() []graphs.ChartModel {
	list := make([]coincap.History, 0)
	err := e.cacheService.GetList(e.name, &list, 0, -1)
	if err != nil {
		slog.Error("Calculate:cacheService.GetLis", "coin", e.name, "error", err)
		return nil
	}

	response := make([]graphs.ChartModel, len(list), len(list))

	s := NewSma(e.name, e.period)
	prev := s.createAvg(list[:e.period])
	response[e.period] = graphs.ChartModel{
		Value: prev,
		Date:  list[e.period].Date,
	}

	for i := e.period + 1; i < len(list); i++ {
		response[i] = graphs.ChartModel{
			Date:  list[i].Date,
			Value: (list[i].PriceUsd-prev)*e.smoothing + prev,
		}
		prev = response[i].Value
	}

	return response
}

func (e *ema) Draw(list []graphs.ChartModel) func(w http.ResponseWriter, r *http.Request) {
	return func(w http.ResponseWriter, r *http.Request) {
		chart := localCharts.CreateLineChart("EMA (Exponential Moving Average)")
		dates, data := localCharts.ChartModel2lineData(list)

		chart.SetXAxis(dates).AddSeries(fmt.Sprintf("%v", e.period), data)

		err := chart.Render(w)
		if err != nil {
			slog.Error("Draw:chart.Render", "error", err)
		}
	}
}
