package model

import "time"

type ChartModel struct {
	Date  time.Time
	Value float32
}

func (c ChartModel) GetDate() time.Time {
	return c.Date
}

func (c ChartModel) GetValue() float32 {
	return c.Value
}
