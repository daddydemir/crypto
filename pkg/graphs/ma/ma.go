package ma

import (
	"github.com/daddydemir/crypto/config/log"
	"github.com/daddydemir/crypto/pkg/cache"
	localCharts "github.com/daddydemir/crypto/pkg/charts"
	"github.com/daddydemir/crypto/pkg/graphs"
	"github.com/daddydemir/crypto/pkg/remote/coincap"
	"log/slog"
	"net/http"
)

type ma struct {
	name         string
	period       int
	cacheService cache.Cache
}

func NewMa(name string, period int) *ma {
	return &ma{
		name:         name,
		period:       period,
		cacheService: cache.GetCacheService(),
	}
}

func (m *ma) Index() float32 {
	panic("invalid method")
}

func (m *ma) Calculate() []graphs.ChartModel {
	list := make([]coincap.History, 0)
	err := m.cacheService.GetList(m.name, &list, 0, -1)
	if err != nil {
		slog.Error("Calculate:CacheService.GetList", "coin", m.name, "error", err)
		return nil
	}
	var response []graphs.ChartModel

	for i := 0; i < len(list); i++ {
		response = append(response, graphs.ChartModel{
			Date:  list[i].Date,
			Value: list[i].PriceUsd,
		})
	}

	return response
}

func (m *ma) Draw(list []graphs.ChartModel) func(w http.ResponseWriter, r *http.Request) {
	shortList := NewEma(m.name, 7).Calculate()
	middleList := NewEma(m.name, 25).Calculate()
	longList := NewEma(m.name, 99).Calculate()

	if len(shortList) == 0 || len(middleList) == 0 || len(longList) == 0 {
		log.Errorln("insufficient data for MA calculation")
		return nil
	}

	return func(w http.ResponseWriter, r *http.Request) {
		dates, short := localCharts.ChartModel2lineData(shortList)
		_, middle := localCharts.ChartModel2lineData(middleList)
		_, long := localCharts.ChartModel2lineData(longList)
		_, original := localCharts.ChartModel2lineData(list)
		chart := localCharts.CreateLineChart("MA (Moving Average)")

		chart.SetXAxis(dates).
			AddSeries("short", short).
			AddSeries("middle", middle).
			AddSeries("long", long).
			AddSeries(m.name, original)

		err := chart.Render(w)
		if err != nil {
			slog.Error("Draw:chart.Render", "error", err)
		}
	}
}
