package graphs

import (
	"net/http"
	"time"
)

type Graph interface {
	Calculate() []ChartModel
	Draw(list []ChartModel) func(w http.ResponseWriter, r *http.Request)
	Index() float32
}

type ChartModel struct {
	Value float32
	Date  time.Time
}
