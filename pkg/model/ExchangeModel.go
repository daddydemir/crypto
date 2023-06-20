package model

import "github.com/google/uuid"

type ExchangeModel struct {
	Id           uuid.UUID
	ExchangeId   string
	Name         string
	InstantPrice float32
}
