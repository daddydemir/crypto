package charts

import (
	"github.com/daddydemir/crypto/internal/port/chart"
	"github.com/daddydemir/crypto/pkg/graphs"
	"github.com/go-echarts/go-echarts/v2/charts"
	"github.com/go-echarts/go-echarts/v2/opts"
)

func ChartModel2lineData(list []graphs.ChartModel) ([]string, []opts.LineData) {

	dates := make([]string, 0, len(list))
	data := make([]opts.LineData, 0, len(list))

	for _, chart := range list {
		dates = append(dates, chart.Date.Format("2006-01-02"))
		data = append(data, opts.LineData{Value: chart.Value})
	}

	return dates, data
}

func getTitleGlobalOpts(title string) charts.GlobalOpts {
	return charts.WithTitleOpts(opts.Title{
		Title: title,
	})
}

func GetLineStyle(color string) charts.SeriesOpts {
	return charts.WithLineStyleOpts(opts.LineStyle{
		Color: color,
	})
}

func LineDataFromList(list []chart.LineConvertible) ([]string, []opts.LineData) {
	dates := make([]string, 0, len(list))
	data := make([]opts.LineData, 0, len(list))

	for _, point := range list {
		dates = append(dates, point.GetDate().Format("2006-01-02"))
		data = append(data, opts.LineData{Value: point.GetValue()})
	}

	return dates, data
}
