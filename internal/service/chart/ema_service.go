package chart

import (
	"github.com/daddydemir/crypto/internal/domain/model"
	"github.com/daddydemir/crypto/internal/port"
	"github.com/daddydemir/crypto/pkg/charts"
	"log/slog"
	"net/http"
)

type EmaService struct {
	coin      string
	period    int
	cache     port.Cache
	smoothing float32
	list      []model.ChartModel
}

func NewEmaService(coin string, period int, cache port.Cache) *EmaService {
	return &EmaService{
		coin:      coin,
		period:    period,
		cache:     cache,
		smoothing: 2 / float32(1+period),
	}
}

func (e *EmaService) Calculate() []model.ChartModel {
	var history []model.History
	err := e.cache.GetList(e.coin, &history, 0, -1)
	if err != nil {
		slog.Error("Ema:Calculate:GetList", "coin", e.coin, "error", err)
		return nil
	}

	if len(history) < e.period+1 {
		slog.Warn("Ema:Calculate", "coin", e.coin, "message", "not enough data")
		return nil
	}

	e.list = make([]model.ChartModel, len(history))

	prev := e.simpleMovingAverage(history[:e.period])
	e.list[e.period] = model.ChartModel{
		Date:  history[e.period].Date,
		Value: prev,
	}

	for i := e.period + 1; i < len(history); i++ {
		val := (history[i].PriceUsd-prev)*e.smoothing + prev
		e.list[i] = model.ChartModel{
			Date:  history[i].Date,
			Value: val,
		}
		prev = val
	}

	return e.list
}

func (e *EmaService) Draw(w http.ResponseWriter, r *http.Request) {
	if len(e.list) == 0 {
		e.Calculate()
	}

	dates, data := charts.LineDataFromList(toLineConvertible(e.list))

	chart := charts.CreateLineChart("Exponential Moving Average")
	chart.SetXAxis(dates).
		AddSeries(e.coin, data)

	if err := chart.Render(w); err != nil {
		slog.Error("Ema:Draw", "error", err)
	}
}

func (e *EmaService) simpleMovingAverage(data []model.History) float32 {
	var sum float32
	for _, h := range data {
		sum += h.PriceUsd
	}
	return sum / float32(len(data))
}
