package graphs

import (
	"github.com/daddydemir/crypto/pkg/remote/coincap"
	"github.com/go-echarts/go-echarts/v2/charts"
	"github.com/go-echarts/go-echarts/v2/opts"
	"github.com/go-echarts/go-echarts/v2/types"
	"math"
)

var GlobalOptions = []charts.GlobalOpts{
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
	}),
}

var SeriesOptions = []charts.SeriesOpts{
	charts.WithLineStyleOpts(opts.LineStyle{
		Color: "orange",
	}),
}

func PrepareData(list []ChartModel) ([]string, []opts.LineData) {
	dates := make([]string, len(list))
	values := make([]opts.LineData, len(list))

	for i := 0; i < len(list); i++ {
		dates[i] = list[i].Date.Format("2006-01-02")
		values[i] = opts.LineData{
			Value: math.Floor(float64(list[i].Value)*10000) / 10000,
		}
	}

	return dates, values
}

func GetTitleGlobalOpts(title string) charts.GlobalOpts {
	return charts.WithTitleOpts(opts.Title{
		Title: title,
	})
}

func PrepareDataWithHistory(list []coincap.History) ([]string, []opts.LineData) {
	dates := make([]string, len(list))
	values := make([]opts.LineData, len(list))

	for i := 0; i < len(list); i++ {
		dates[i] = list[i].Date.Format("2006-01-02")
		values[i] = opts.LineData{
			Value: math.Floor(float64(list[i].PriceUsd)*10000) / 10000,
		}
	}

	return dates, values
}
