package chart

import (
	"github.com/daddydemir/crypto/internal/domain/model"
	"github.com/daddydemir/crypto/internal/port"
	"github.com/daddydemir/crypto/internal/port/chart"
	"github.com/daddydemir/crypto/pkg/charts"
	"log/slog"
	"math"
	"net/http"
)

type BollingerBandsService struct {
	coin   string
	period int
	cache  port.Cache

	middle   []model.ChartModel
	upper    []model.ChartModel
	lower    []model.ChartModel
	original []model.ChartModel
}

func NewBollingerBandsService(coin string, period int, cache port.Cache) *BollingerBandsService {
	return &BollingerBandsService{
		coin:   coin,
		period: period,
		cache:  cache,
	}
}

func (b *BollingerBandsService) CalculateBands() (middle, upper, lower, original []model.ChartModel) {
	var history []model.History
	err := b.cache.GetList(b.coin, &history, 0, -1)
	if err != nil {
		slog.Error("BollingerBands:Calculate:GetList", "coin", b.coin, "error", err)
		return
	}

	for i := 0; i <= len(history)-b.period; i++ {
		window := history[i : i+b.period]
		mean := b.mean(window)
		std := b.stdDev(window, mean)

		date := window[b.period-1].Date

		b.middle = append(b.middle, model.ChartModel{Date: date, Value: mean})
		b.upper = append(b.upper, model.ChartModel{Date: date, Value: mean + 2*std})
		b.lower = append(b.lower, model.ChartModel{Date: date, Value: mean - 2*std})
		b.original = append(b.original, model.ChartModel{Date: date, Value: window[b.period-1].PriceUsd})
	}

	return b.middle, b.upper, b.lower, b.original
}

func (b *BollingerBandsService) DrawBands(w http.ResponseWriter, r *http.Request) {
	chart := charts.CreateLineChart("Bollinger Bands")

	dates, mid := charts.LineDataFromList(toLineConvertible(b.middle))
	_, up := charts.LineDataFromList(toLineConvertible(b.upper))
	_, down := charts.LineDataFromList(toLineConvertible(b.lower))
	_, orig := charts.LineDataFromList(toLineConvertible(b.original))

	chart.SetXAxis(dates).
		AddSeries("Lower", down).
		AddSeries("Middle", mid).
		AddSeries("Upper", up).
		AddSeries(b.coin, orig)

	if err := chart.Render(w); err != nil {
		slog.Error("DrawBands", "coin", b.coin, "error", err)
	}
}

func (b *BollingerBandsService) mean(list []model.History) float32 {
	var sum float32
	for _, h := range list {
		sum += h.PriceUsd
	}
	return sum / float32(b.period)
}

func (b *BollingerBandsService) stdDev(list []model.History, mean float32) float32 {
	var sum float64
	for _, h := range list {
		diff := float64(mean - h.PriceUsd)
		sum += diff * diff
	}
	return float32(math.Sqrt(sum / float64(b.period-1)))
}

func toLineConvertible(list []model.ChartModel) []chart.LineConvertible {
	result := make([]chart.LineConvertible, len(list))
	for i := range list {
		result[i] = list[i]
	}

	return result
}
