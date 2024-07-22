package ma

import (
	"fmt"
	"github.com/daddydemir/crypto/pkg/coincap"
	"github.com/daddydemir/crypto/pkg/graphs"
	"github.com/go-echarts/go-echarts/v2/charts"
	"github.com/gorilla/mux"
	"net/http"
)

type Ema struct{}

func (e Ema) calculate(coin string, period int) []graphs.ChartModel {

	histories := coincap.HistoryWithId(coin)
	response := make([]graphs.ChartModel, len(histories), len(histories))

	sma := Sma{}
	firstPoint := sma.averagePerDay(histories[:period])
	factor := e.smoothingFactor(period)
	prev := firstPoint

	firstDay := graphs.ChartModel{
		Date:  histories[period-1].Date,
		Value: firstPoint,
	}

	response[period-1] = firstDay

	for i := period; i < len(histories); i++ {
		val := (histories[i].PriceUsd-prev)*factor + prev
		prev = val

		response[i] = graphs.ChartModel{
			Date:  histories[i].Date,
			Value: val,
		}
	}

	return response
}

func (e Ema) smoothingFactor(period int) float32 {
	return 2 / float32(period+1)
}

func (e Ema) Draw() func(w http.ResponseWriter, r *http.Request) {
	return e.draw
}

func (e Ema) draw(w http.ResponseWriter, r *http.Request) {

	vars := mux.Vars(r)
	coin := vars["coin"]

	list := e.calculate(coin, 3)
	_, values := graphs.PrepareData(list)
	histories := coincap.HistoryWithId(coin)
	dates, datas := graphs.PrepareDataWithHistory(histories)

	line := charts.NewLine()
	line.SetGlobalOptions(graphs.GlobalOptions...)
	line.SetGlobalOptions(graphs.GetTitleGlobalOpts("EMA (Exponential Moving Average)"))

	line.SetXAxis(dates).
		AddSeries(coin, values, graphs.SeriesOptions...).
		AddSeries("original", datas)

	err := line.Render(w)
	if err != nil {
		fmt.Printf("render err: %v\n", err)
	}
}
