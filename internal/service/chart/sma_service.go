package chart

import (
	"github.com/daddydemir/crypto/internal/domain/model"
	"github.com/daddydemir/crypto/internal/port"
	"github.com/daddydemir/crypto/pkg/charts"
	"log/slog"
	"net/http"
)

type SmaService struct {
	coin   string
	period int
	cache  port.Cache
	list   []model.ChartModel
}

func NewSmaService(coin string, period int, cache port.Cache) *SmaService {
	return &SmaService{
		coin:   coin,
		period: period,
		cache:  cache,
	}
}

func (s *SmaService) Calculate() []model.ChartModel {
	var history []model.History
	err := s.cache.GetList(s.coin, &history, 0, -1)
	if err != nil {
		slog.Error("Sma:Calculate:GetList", "coin", s.coin, "error", err)
		return nil
	}

	if len(history) < s.period {
		slog.Warn("Sma:Calculate", "coin", s.coin, "message", "insufficient data")
		return nil
	}

	s.list = make([]model.ChartModel, 0, len(history)-s.period+1)

	for i := 0; i <= len(history)-s.period; i++ {
		avg := s.average(history[i : i+s.period])
		s.list = append(s.list, model.ChartModel{
			Date:  history[i+s.period-1].Date,
			Value: avg,
		})
	}

	return s.list
}

func (s *SmaService) Draw(w http.ResponseWriter, r *http.Request) {
	if len(s.list) == 0 {
		s.Calculate()
	}

	dates, data := charts.LineDataFromList(toLineConvertible(s.list))
	chart := charts.CreateLineChart("Simple Moving Average")

	chart.SetXAxis(dates).
		AddSeries(s.coin, data)

	if err := chart.Render(w); err != nil {
		slog.Error("Sma:Draw", "error", err)
	}
}

func (s *SmaService) average(history []model.History) float32 {
	var sum float32
	for _, h := range history {
		sum += h.PriceUsd
	}
	return sum / float32(len(history))
}
