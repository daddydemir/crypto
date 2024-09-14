package charts

import (
	"github.com/go-echarts/go-echarts/v2/charts"
	"github.com/go-echarts/go-echarts/v2/opts"
	"github.com/go-echarts/go-echarts/v2/types"
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
