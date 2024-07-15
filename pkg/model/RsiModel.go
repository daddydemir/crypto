package model

import "time"

type RsiModel struct {
	Index  float32
	Date   time.Time
	Period int
	Data   []float32
	Gain   []RsiDate
	Loss   []RsiDate
}

type RsiDate struct {
	Val  float32
	Date time.Time
}
