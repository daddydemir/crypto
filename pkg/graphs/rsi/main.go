package rsi

import (
	"github.com/daddydemir/crypto/pkg/cache"
	localCharts "github.com/daddydemir/crypto/pkg/charts"
	"github.com/daddydemir/crypto/pkg/graphs"
	"github.com/daddydemir/crypto/pkg/remote/coincap"
	"github.com/go-echarts/go-echarts/v2/charts"
	"github.com/go-echarts/go-echarts/v2/opts"
	"log/slog"
	"net/http"
)

type rsi struct {
	name         string
	period       int
	cacheService cache.Cache
}

func NewRsi(name string) *rsi {
	return &rsi{
		name:         name,
		period:       14,
		cacheService: cache.GetCacheService(),
	}
}

func (r *rsi) Index() float32 {
	list := new([]coincap.History)
	err := r.cacheService.GetList(r.name, list, -14, -1)
	if err != nil {
		slog.Error("Index:CacheService.GetList", "coin", r.name, "error", err)
		return -1
	}
	if len(*list) != r.period {
		slog.Error("Index: list size is not valid.", "size", len(*list), "period", r.period)
		return -1
	}
	return r.createIndex(*list)
}

func (r *rsi) Calculate() []graphs.ChartModel {
	list := make([]coincap.History, 0)
	err := r.cacheService.GetList(r.name, &list, 0, -1)
	if err != nil {
		slog.Error("Calculate:CacheService.GetList", "coin", r.name, "error", err)
		return nil
	}

	var response []graphs.ChartModel

	for z := 0; z < len(list)-r.period; z++ {
		response = append(response, graphs.ChartModel{
			Date:  list[z+r.period].Date,
			Value: r.createIndex(list[z+1 : z+r.period+1]),
		})
	}

	return response
}

func (r *rsi) Draw(list []graphs.ChartModel) func(_ http.ResponseWriter, _ *http.Request) {
	return func(w http.ResponseWriter, _ *http.Request) {
		chart := localCharts.CreateLineChart("RSI (Relative Strength Index)")
		dates, data := localCharts.ChartModel2lineData(list)
		chart.SetXAxis(dates).AddSeries(
			r.name,
			data,
			charts.WithMarkLineNameYAxisItemOpts(
				opts.MarkLineNameYAxisItem{Name: "70", ValueDim: "70", YAxis: 70},
				opts.MarkLineNameYAxisItem{Name: "30", ValueDim: "30", YAxis: 30},
			),
		)

		err := chart.Render(w)
		if err != nil {
			slog.Error("Draw:charts.Render", "error", err)
		}
	}
}

func (r *rsi) createIndex(list []coincap.History) float32 {

	var (
		gainSum float32
		lossSum float32
		period  = float32(r.period)
	)

	for i := 1; i < len(list); i++ {
		tmp := list[i].PriceUsd - list[i-1].PriceUsd
		if tmp > 0 {
			gainSum += tmp
		} else {
			lossSum += tmp
		}
	}

	gainSum /= period
	lossSum /= period * -1

	index := 100 - (100 / (1 + (gainSum / lossSum)))
	return index
}
