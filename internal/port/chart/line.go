package chart

import (
	"time"
)

type LineConvertible interface {
	GetDate() time.Time
	GetValue() float32
}
