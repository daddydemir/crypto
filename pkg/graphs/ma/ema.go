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

	list7 := e.calculate(coin, 7)
	list25 := e.calculate(coin, 25)
	list99 := e.calculate(coin, 99)

	_, values7 := graphs.PrepareData(list7)
	_, values25 := graphs.PrepareData(list25)
	_, values99 := graphs.PrepareData(list99)

	histories := coincap.HistoryWithId(coin)
	dates, datas := graphs.PrepareDataWithHistory(histories)

	line := charts.NewLine()
	line.SetGlobalOptions(graphs.GlobalOptions...)
	line.SetGlobalOptions(graphs.GetTitleGlobalOpts("EMA (Exponential Moving Average)"))

	line.SetXAxis(dates).
		AddSeries("7", values7).
		AddSeries("25", values25).
		AddSeries("99", values99).
		AddSeries("original", datas)

	err := line.Render(w)
	if err != nil {
		fmt.Printf("render err: %v\n", err)
	}
}
