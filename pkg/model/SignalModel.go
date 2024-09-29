package model

import (
	"github.com/google/uuid"
	"time"
)

type SignalModel struct {
	ID         uuid.UUID
	ExchangeId string
	CreateDate time.Time
	Short      float32
	Middle     float32
	Long       float32
	Rsi        float32
	Result     string
}
