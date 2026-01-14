package domain

import "time"

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

type Signal struct {
	Id     string  `json:"id"`
	Name   string  `json:"name"`
	Symbol string  `json:"symbol"`
	Price  float64 `json:"price"`
	Points []Point `json:"points"`
}

func CalculateSeries(dates []time.Time, prices []float64) []Point {
	var result []Point
	for i := range prices {
		if i < 98 {
			continue
		}
		ma7 := mean(prices[i-6 : i+1])
		ma25 := mean(prices[i-24 : i+1])
		ma99 := mean(prices[i-98 : i+1])

		result = append(result, Point{
			Date: dates[i],
			MA7:  ma7,
			MA25: ma25,
			MA99: ma99,
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
