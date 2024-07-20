package ma

import (
	"fmt"
	"github.com/daddydemir/crypto/pkg/coincap"
	"github.com/daddydemir/crypto/pkg/graphs"
	"github.com/go-echarts/go-echarts/v2/charts"
	"net/http"
)

var Coin string

func calculate(coin string, period int) []graphs.ChartModel {

	histories := coincap.HistoryWithId(coin)
	response := make([]graphs.ChartModel, 0, len(histories)-period)
	start, end := 0, period

	for end <= len(histories) {
		average := averagePerDay(histories[start:end])

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

func averagePerDay(list []coincap.History) float32 {

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

	list := calculate(Coin, 10)

	line := charts.NewLine()

	line.SetGlobalOptions(graphs.GlobalOptions...)

	dates, values := graphs.PrepareData(list)

	line.SetGlobalOptions(graphs.GetTitleGlobalOpts("SMA (Simple Moving Average)"))

	line.SetXAxis(dates).AddSeries(Coin, values, graphs.SeriesOptions...)

	err := line.Render(w)
	if err != nil {
		fmt.Printf("render error: %v", err)
	}
}
