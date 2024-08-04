package graphs

import (
	"net/http"
	"time"
)

type Graph interface {
	Calculate(coin string, period int) []ChartModel
	Draw() func(w http.ResponseWriter, r *http.Request)
	Index(coin string) float32
}

type ChartModel struct {
	Value float32
	Date  time.Time
}
