package model

import "time"

type Coin struct {
	Id       string
	Symbol   string
	Name     string
	PriceUsd float32 `json:"priceUsd,string"`
}

type History struct {
	PriceUsd float32 `json:"priceUsd,string"`
	Date     time.Time
}
