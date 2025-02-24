package model

import "github.com/google/uuid"

type ExchangeModel struct {
	Id           uuid.UUID
	ExchangeId   string
	Name         string
	CoinImage    string
	InstantPrice float32
}

type ExchangeRepository interface {
	FindAll() ([]ExchangeModel, error)
	SaveAll([]ExchangeModel) error
}
