package domain

import (
	"math"
	"time"
)

type Point struct {
	Date      time.Time
	MA20      float64
	UpperBand float64
	LowerBand float64
}

type PriceData struct {
	Date  time.Time
	Price float64
}

type Signal struct {
	Id     string  `json:"id"`
	Name   string  `json:"name"`
	Symbol string  `json:"symbol"`
	Price  float64 `json:"price"`
	Point  Point   `json:"point"`
}

func CalculateBollinger(prices []float64, dates []time.Time) []Point {
	var result []Point
	for i := 19; i < len(prices); i++ {
		window := prices[i-19 : i+1]
		ma := mean(window)
		std := stdDev(window)
		result = append(result, Point{
			Date:      dates[i],
			MA20:      ma,
			UpperBand: ma + 2*std,
			LowerBand: ma - 2*std,
		})
	}
	return result
}

func mean(values []float64) float64 {
	sum := 0.0
	for _, v := range values {
		sum += v
	}
	return sum / float64(len(values))
}

func stdDev(values []float64) float64 {
	m := mean(values)
	var sum float64
	for _, v := range values {
		sum += (v - m) * (v - m)
	}
	return sqrt(sum / float64(len(values)))
}

func sqrt(x float64) float64 {
	// simple Newton-Raphson veya math.Sqrt
	return math.Sqrt(x)
}
