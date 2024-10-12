package bollingerBands

import (
	"github.com/daddydemir/crypto/pkg/cache"
	localCharts "github.com/daddydemir/crypto/pkg/charts"
	"github.com/daddydemir/crypto/pkg/graphs"
	"github.com/daddydemir/crypto/pkg/remote/coincap"
	"log/slog"
	"math"
	"net/http"
)

var upperBand []graphs.ChartModel
var lowerBand []graphs.ChartModel
var original []graphs.ChartModel

type bollingerBands struct {
	name         string
	period       int
	cacheService cache.Cache
}

func NewBollingerBands(name string, period int) *bollingerBands {
	return &bollingerBands{name: name,
		period:       period,
		cacheService: cache.GetCacheService(),
	}
}

func (b *bollingerBands) Index() float32 {
	return 0
}

func (b *bollingerBands) Calculate() []graphs.ChartModel {

	list := make([]coincap.History, 0)
	err := b.cacheService.GetList(b.name, &list, 0, -1)
	if err != nil {
		slog.Error("Calculate:CacheService.GetList", "coin", b.name, "error", err)
		return nil
	}
	var response []graphs.ChartModel

	for i := 0; i < len(list)-b.period+1; i++ {

		mean := b.arithmeticMean(list[i : i+b.period])
		deviation := b.standardDeviation(list[i:i+b.period], mean)

		chartModel := graphs.ChartModel{
			Date:  list[i+b.period-1].Date,
			Value: mean,
		}

		upper := graphs.ChartModel{
			Date:  list[i+b.period-1].Date,
			Value: mean + (deviation * 2),
		}

		lower := graphs.ChartModel{
			Date:  list[i+b.period-1].Date,
			Value: mean - (deviation * 2),
		}

		o := graphs.ChartModel{
			Date:  list[i+b.period-1].Date,
			Value: list[i+b.period-1].PriceUsd,
		}

		response = append(response, chartModel)
		upperBand = append(upperBand, upper)
		lowerBand = append(lowerBand, lower)
		original = append(original, o)
	}

	return response
}

func (b *bollingerBands) Draw(list []graphs.ChartModel) func(w http.ResponseWriter, r *http.Request) {
	return func(w http.ResponseWriter, r *http.Request) {
		chart := localCharts.CreateLineChart("Bollinger Bands")
		dates, data := localCharts.ChartModel2lineData(list)
		_, uppers := localCharts.ChartModel2lineData(upperBand)
		_, lowers := localCharts.ChartModel2lineData(lowerBand)
		_, origin := localCharts.ChartModel2lineData(original)
		chart.SetXAxis(dates).
			AddSeries("low", lowers).
			AddSeries("middle", data).
			AddSeries("up", uppers).
			AddSeries(b.name, origin)

		err := chart.Render(w)
		upperBand = nil
		lowerBand = nil
		original = nil
		if err != nil {
			slog.Error("Draw:Render", "chart", b.name, "error", err)
		}
	}
}

func (b *bollingerBands) arithmeticMean(list []coincap.History) float32 {
	var sum float32

	for _, history := range list {
		sum += history.PriceUsd
	}

	return sum / float32(b.period)
}

func (b *bollingerBands) standardDeviation(list []coincap.History, arithmeticMean float32) float32 {
	var sum float64

	for _, history := range list {
		pow := math.Pow(float64(arithmeticMean-history.PriceUsd), 2.0)
		sum += pow
	}

	return float32(math.Sqrt(sum / float64(b.period-1)))
}
