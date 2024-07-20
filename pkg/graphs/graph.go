package graphs

import (
	"time"
)

type Graph interface {
	Calculate()
	Draw()
}

type ChartModel struct {
	Value float32
	Date  time.Time
}
