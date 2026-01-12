package domain

import (
	"time"
)

type Point struct {
	Date time.Time `json:"date"`
	MA7  float64   `json:"ma7"`
	MA25 float64   `json:"ma25"`
	MA99 float64   `json:"ma99"`
}

type PriceData struct {
	Date  time.Time
	Price float64
}

func CalculateSeries(dates []time.Time, prices []float64) []Point {
	if len(dates) != len(prices) || len(prices) == 0 {
		return nil
	}

	ma7 := ema(prices, 7)
	ma25 := ema(prices, 25)
	ma99 := ema(prices, 99)

	var points []Point
	for i := range prices {
		if i < 99 {
			continue
		}
		points = append(points, Point{
			Date: dates[i],
			MA7:  ma7[i],
			MA25: ma25[i],
			MA99: ma99[i],
		})
	}
	return points
}

func ema(prices []float64, period int) []float64 {
	alpha := 2.0 / float64(period+1)
	result := make([]float64, len(prices))
	result[0] = prices[0]
	for i := 1; i < len(prices); i++ {
		result[i] = alpha*prices[i] + (1-alpha)*result[i-1]
	}
	return result
}
