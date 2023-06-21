package adapter

import (
	"github.com/daddydemir/crypto/pkg/model"
	"github.com/google/uuid"
	"time"
)

type Adapter struct {
	Id                    string  `json:"id"`
	Symbol                string  `json:"symbol"`
	Name                  string  `json:"name"`
	Image                 string  `json:"image"`
	CurrentPrice          float32 `json:"current_price"`
	MarketCap             int64   `json:"market_cap"`
	MarketCapRank         int     `json:"market_cap_rank"`
	FullyDilutedValuation int64   `json:"fully_diluted_valuation"`
	TotalVolume           float64 `json:"total_volume"`
	High24H               float32 `json:"high_24h"`
	Low24H                float32 `json:"low_24h"`
	//PriceChange24H               string    `json:"price_change_24h"`
	//PriceChangePercentage24H     float64   `json:"price_change_percentage_24h"`
	//MarketCapChange24H           int64     `json:"market_cap_change_24h"`
	//MarketCapChangePercentage24H float64   `json:"market_cap_change_percentage_24h"`
	CirculatingSupply   float64   `json:"circulating_supply"`
	TotalSupply         float64   `json:"total_supply"`
	MaxSupply           float64   `json:"max_supply"`
	Ath                 float32   `json:"ath"`
	AthChangePercentage float64   `json:"ath_change_percentage"`
	AthDate             time.Time `json:"ath_date"`
	Atl                 float64   `json:"atl"`
	AtlChangePercentage float64   `json:"atl_change_percentage"`
	AtlDate             time.Time `json:"atl_date"`
	LastUpdated         time.Time `json:"last_updated"`
}

func (a *Adapter) AdapterToDaily(morning bool) model.DailyModel {
	var daily model.DailyModel

	daily.Id = uuid.New()
	daily.Min = a.Low24H
	daily.Max = a.High24H
	daily.Avg = (a.Low24H + a.High24H) / 2
	daily.Date = time.Now()
	daily.Rate = ((a.High24H - a.Low24H) * 100) / ((a.Low24H + a.High24H) / 2)
	daily.Modulus = a.High24H - a.Low24H
	daily.ExchangeId = a.Symbol
	if morning {
		daily.FirstPrice = a.CurrentPrice
	} else {
		daily.LastPrice = a.CurrentPrice
	}

	return daily
}

func (a *Adapter) AdapterToExchange() model.ExchangeModel {
	var exchange model.ExchangeModel

	exchange.Id = uuid.New()
	exchange.ExchangeId = a.Symbol
	exchange.Name = a.Name
	exchange.CoinImage = a.Image
	exchange.InstantPrice = a.CurrentPrice

	return exchange
}
