package ma

import (
	"github.com/daddydemir/crypto/pkg/graphs"
	coincap2 "github.com/daddydemir/crypto/pkg/remote/coincap"
	"github.com/go-echarts/go-echarts/v2/charts"
	"log/slog"
	"net/http"
)

type Sma struct {
}

var Coin string

func (s Sma) calculate(coin string, period int) []graphs.ChartModel {

	histories := coincap2.HistoryWithId(coin)
	response := make([]graphs.ChartModel, 0, len(histories)-period)
	start, end := 0, period

	for end <= len(histories) {
		average := s.averagePerDay(histories[start:end])

		item := graphs.ChartModel{
			Date:  histories[end-1].Date,
			Value: average,
		}

		response = append(response, item)

		start++
		end++
	}

	return response
}

func (s Sma) averagePerDay(list []coincap2.History) float32 {

	var totalPrice float32
	var period int

	for _, history := range list {
		totalPrice += history.PriceUsd
		period++
	}

	return totalPrice / float32(period)
}

func Draw(coin string) func(w http.ResponseWriter, r *http.Request) {

	Coin = coin
	return draw
}

func draw(w http.ResponseWriter, _ *http.Request) {

	sma := Sma{}

	list := sma.calculate(Coin, 10)

	line := charts.NewLine()

	line.SetGlobalOptions(graphs.GlobalOptions...)

	dates, values := graphs.PrepareData(list)

	line.SetGlobalOptions(graphs.GetTitleGlobalOpts("SMA (Simple Moving Average)"))

	line.SetXAxis(dates).AddSeries(Coin, values, graphs.SeriesOptions...)

	err := line.Render(w)
	if err != nil {
		slog.Error("sma.draw:Render", "err", err)
	}
}
