package charts

import (
	"github.com/go-echarts/go-echarts/v2/charts"
)

func newCustomLineChart() *charts.Line {
	line := charts.NewLine()
	line.SetGlobalOptions(GlobalOptions...)

	return line
}

func CreateLineChart(title string) *charts.Line {
	lineChart := newCustomLineChart()
	lineChart.SetGlobalOptions(getTitleGlobalOpts(title))

	return lineChart
}
