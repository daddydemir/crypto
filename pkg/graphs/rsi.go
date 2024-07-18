package graphs

import (
	"github.com/daddydemir/crypto/pkg/coincap"
	"github.com/daddydemir/crypto/pkg/model"
	"github.com/go-echarts/go-echarts/v2/charts"
	"github.com/go-echarts/go-echarts/v2/opts"
	"github.com/go-echarts/go-echarts/v2/types"
	"math"
	"net/http"
)

type RSI struct {
	model.RsiModel
}

func (r RSI) Calculate(s string) []model.RsiModel {
	history := coincap.HistoryWithId(s)
	response := make([]model.RsiModel, 0)

	start, end := 0, 14

	for end <= len(history) {
		rsiModel := calculateIndex(history[start:end])
		response = append(response, rsiModel)
		start++
		end++
	}

	return response
}

func calculateIndex(list []coincap.History) model.RsiModel {

	response := new(model.RsiModel)
	period := 14
	gain := make([]model.RsiDate, 0)
	loss := make([]model.RsiDate, 0)
	data := make([]float32, period)

	for i, d := range list {

		data = append(data, d.PriceUsd)

		if i == 0 {
			continue
		}

		tmp := d.PriceUsd - list[i-1].PriceUsd

		if tmp > 0 {
			gain = append(gain, model.RsiDate{
				Date: d.Date,
				Val:  tmp,
			})
		} else {
			loss = append(loss, model.RsiDate{
				Date: d.Date,
				Val:  tmp,
			})
		}
	}

	gainAverage := sum(gain) / 14
	lossAverage := sum(loss) / 14 * -1

	index := 100 - (100 / (1 + (gainAverage / lossAverage)))

	response.Gain = gain
	response.Loss = loss
	response.Data = data
	response.Index = index
	response.Date = list[len(list)-1].Date

	return *response
}

func sum(arr []model.RsiDate) float32 {
	var t float32
	for _, d := range arr {
		t += d.Val
	}
	return t
}

func (r RSI) Draw(list []model.RsiModel) func(w http.ResponseWriter, r *http.Request) {
	return useCharts(list)
}

var List []model.RsiModel

func useCharts(list []model.RsiModel) func(w http.ResponseWriter, r *http.Request) {
	List = list
	return drawChart
}

func drawChart(w http.ResponseWriter, _ *http.Request) {
	line := charts.NewLine()

	line.SetGlobalOptions(
		charts.WithTitleOpts(opts.Title{
			Title: "RSI",
		}),
		charts.WithXAxisOpts(opts.XAxis{
			Type: "Time",
		}),
		charts.WithYAxisOpts(opts.YAxis{
			Name: "Index",
		}),
	)

	dates, values := getDataForCharts(List)

	line.SetGlobalOptions(
		charts.WithTitleOpts(opts.Title{
			Title: "RSI",
		}),
		charts.WithXAxisOpts(opts.XAxis{
			SplitNumber: 20,
		}),
		charts.WithYAxisOpts(opts.YAxis{
			Scale: opts.Bool(true),
		}),
		charts.WithDataZoomOpts(opts.DataZoom{
			Type:       "inside",
			Start:      50,
			End:        100,
			XAxisIndex: []int{0},
		}),
		charts.WithDataZoomOpts(opts.DataZoom{
			Type:       "slider",
			Start:      50,
			End:        100,
			XAxisIndex: []int{0},
		}),
		charts.WithInitializationOpts(opts.Initialization{
			Theme: types.ThemeWesteros,
			//Theme: types.ThemeRoma,
		}),
	)

	line.SetXAxis(dates).AddSeries("tron", generateLineItems(values),
		charts.WithMarkLineNameYAxisItemOpts(
			opts.MarkLineNameYAxisItem{Name: "70", ValueDim: "70", YAxis: 70},
			opts.MarkLineNameYAxisItem{Name: "30", ValueDim: "30", YAxis: 30},
		),
		charts.WithLineStyleOpts(opts.LineStyle{
			Color: "orange",
		}),
	)

	err := line.Render(w)
	if err != nil {
		// todo
	}
}

func getDataForCharts(list []model.RsiModel) ([]string, []float64) {
	dates := make([]string, len(list))
	values := make([]float64, len(list))

	for i := 0; i < len(list); i++ {
		dates[i] = list[i].Date.Format("2006-01-02")
		values[i] = math.Floor(float64(list[i].Index)*100) / 100
	}

	return dates, values
}

func generateLineItems(values []float64) []opts.LineData {
	items := make([]opts.LineData, len(values))
	for i, v := range values {
		items[i] = opts.LineData{Value: v}
	}
	return items
}
