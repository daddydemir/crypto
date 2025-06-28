package model

type CoinViewModel struct {
	Index     int
	Name      string
	Symbol    string
	PriceUsd  float32
	Rsi       float32
	RsiClass  string
	Id        string
	GraphUrls map[string]string
}
