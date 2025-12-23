package domain

import (
	"fmt"
	"math"
	"time"
)

type AtrPoint struct {
	Symbol         string
	CurrentHigh    float64 `gorm:"column:current_high"`
	CurrentLow     float64 `gorm:"column:current_low"`
	YesterdayClose float64 `gorm:"column:yesterday_close"`
	Time           time.Time
}

type Point struct {
	Time  time.Time
	Point float64
}

func CalculateATR(prices []AtrPoint) ([]Point, error) {
	var result []Point
	period := 14
	first := firstATR(prices, period)
	prev := first.Point
	result = append(result, first)

	for i := period; i < len(prices); i++ {

		if !validateData(prices[i].Time, prices[i-1].Time) {
			return nil, fmt.Errorf("[%s] date order is wrong: %s", prices[i].Symbol, prices[i].Time.Format("2006-01-02"))
		}

		newPoint := Point{
			Time:  prices[i].Time,
			Point: (prev + tr(prices[i])) / float64(period),
		}
		result = append(result, newPoint)
		prev = newPoint.Point
	}

	return result, nil
}

func firstATR(prices []AtrPoint, period int) Point {
	sum := 0.0
	for i := 0; i < period; i++ {
		sum += tr(prices[i])
	}
	return Point{
		Time:  prices[period-1].Time,
		Point: sum / float64(period),
	}
}

func tr(p AtrPoint) float64 {
	return max(p.CurrentHigh-p.CurrentLow,
		math.Abs(p.CurrentHigh-p.YesterdayClose),
		math.Abs(p.CurrentLow-p.YesterdayClose),
	)
}

func validateData(today, yesterday time.Time) bool {
	return today.After(yesterday)
}
