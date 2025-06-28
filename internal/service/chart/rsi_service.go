package chart

import (
	"github.com/daddydemir/crypto/internal/domain/model"
	"github.com/daddydemir/crypto/internal/port"
	"github.com/daddydemir/crypto/pkg/charts"
	glb "github.com/go-echarts/go-echarts/v2/charts"
	"github.com/go-echarts/go-echarts/v2/opts"
	"log/slog"
	"math"
	"net/http"
	"time"
)

type RsiService struct {
	client port.CoinCapAPI
	cache  port.Cache
	period int
	coin   string
	list   []model.RsiModel
}

func NewRsiService(coin string, client port.CoinCapAPI, cache port.Cache) *RsiService {
	return &RsiService{
		client: client,
		cache:  cache,
		period: 14,
		coin:   coin,
	}
}

func (r *RsiService) Calculate() []model.RsiModel {
	err, history := r.client.HistoryWithId(r.coin)
	if err != nil || len(history) < r.period {
		slog.Error("Rsi:Calculate", "coin", r.coin, "error", err)
		return nil
	}

	r.list = make([]model.RsiModel, 0)
	for i := 0; i+r.period <= len(history); i++ {
		rsi := calculateIndex(history[i : i+r.period])
		r.list = append(r.list, rsi)
	}

	return r.list
}

func (r *RsiService) Index() float32 {
	today := time.Now().Truncate(24 * time.Hour)
	err, history := r.client.HistoryWithTime(r.coin, today.AddDate(0, 0, -15).UnixNano(), today.UnixNano())
	if err != nil || len(history) < r.period {
		slog.Error("Rsi:Index", "coin", r.coin, "error", err)
		return -1
	}

	return calculateIndex(history[len(history)-r.period:]).Index
}

func (r *RsiService) Draw(w http.ResponseWriter, _ *http.Request) {
	if len(r.list) == 0 {
		r.Calculate()
	}
	dates, values := getRsiChartData(r.list)

	line := charts.CreateLineChart("RSI (Relative Strength Index)")
	line.SetXAxis(dates).AddSeries(r.coin, values,
		glb.WithMarkLineNameYAxisItemOpts(
			opts.MarkLineNameYAxisItem{Name: "70", YAxis: 70},
			opts.MarkLineNameYAxisItem{Name: "30", YAxis: 30},
		),
	)

	if err := line.Render(w); err != nil {
		slog.Error("Rsi:Draw", "error", err)
	}
}

func calculateIndex(list []model.History) model.RsiModel {
	var gainSum, lossSum float32
	for i := 1; i < len(list); i++ {
		delta := list[i].PriceUsd - list[i-1].PriceUsd
		if delta > 0 {
			gainSum += delta
		} else {
			lossSum += -delta
		}
	}

	avgGain := gainSum / float32(len(list))
	avgLoss := lossSum / float32(len(list))
	index := 100 - (100 / (1 + (avgGain / avgLoss)))

	return model.RsiModel{
		Index: index,
		Date:  list[len(list)-1].Date,
	}
}

func getRsiChartData(list []model.RsiModel) ([]string, []opts.LineData) {
	dates := make([]string, len(list))
	values := make([]opts.LineData, len(list))
	for i, item := range list {
		dates[i] = item.Date.Format("2006-01-02")
		values[i] = opts.LineData{Value: math.Floor(float64(item.Index)*100) / 100}
	}
	return dates, values
}
