package chart

import (
	"github.com/daddydemir/crypto/internal/domain/model"
	"github.com/daddydemir/crypto/internal/port"
	"github.com/daddydemir/crypto/pkg/charts"
	"log/slog"
	"net/http"
)

type TripleEmaService struct {
	coin    string
	cache   port.Cache
	periods []int

	shortList  []model.ChartModel
	middleList []model.ChartModel
	longList   []model.ChartModel
	original   []model.ChartModel
}

func NewTripleEmaService(coin string, cache port.Cache) *TripleEmaService {
	return &TripleEmaService{
		coin:    coin,
		cache:   cache,
		periods: []int{7, 25, 99},
	}
}

func (t *TripleEmaService) Calculate() []model.ChartModel {
	var history []model.History
	err := t.cache.GetList(t.coin, &history, 0, -1)
	if err != nil {
		slog.Error("TripleEma:Calculate:GetList", "coin", t.coin, "error", err)
		return nil
	}

	for _, h := range history {
		t.original = append(t.original, model.ChartModel{
			Date:  h.Date,
			Value: h.PriceUsd,
		})
	}
	t.shortList = computeEma(history, t.periods[0])
	t.middleList = computeEma(history, t.periods[1])
	t.longList = computeEma(history, t.periods[2])

	return t.original
}

func (t *TripleEmaService) Draw(w http.ResponseWriter, r *http.Request) {
	if len(t.original) == 0 {
		t.Calculate()
	}

	dates, short := charts.LineDataFromList(toLineConvertible(t.shortList))
	_, middle := charts.LineDataFromList(toLineConvertible(t.middleList))
	_, long := charts.LineDataFromList(toLineConvertible(t.longList))
	_, origin := charts.LineDataFromList(toLineConvertible(t.original))

	chart := charts.CreateLineChart("Triple EMA Trend")
	chart.SetXAxis(dates).
		AddSeries("EMA 7", short).
		AddSeries("EMA 25", middle).
		AddSeries("EMA 99", long).
		AddSeries(t.coin, origin)

	if err := chart.Render(w); err != nil {
		slog.Error("TripleEma:Draw", "error", err)
	}
}

func computeEma(history []model.History, period int) []model.ChartModel {
	if len(history) < period+1 {
		return nil
	}

	smoothing := 2 / float32(1+period)
	result := make([]model.ChartModel, len(history))
	var sum float32
	for i := 0; i < period; i++ {
		sum += history[i].PriceUsd
	}
	prev := sum / float32(period)
	result[period] = model.ChartModel{Date: history[period].Date, Value: prev}

	for i := period + 1; i < len(history); i++ {
		val := (history[i].PriceUsd-prev)*smoothing + prev
		result[i] = model.ChartModel{
			Date:  history[i].Date,
			Value: val,
		}
		prev = val
	}
	return result
}
