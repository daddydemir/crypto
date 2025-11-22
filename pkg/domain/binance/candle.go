package binance

import "time"

type Candle struct {
	ID                       uint
	Symbol                   string    `gorm:"column:symbol"`
	KlineOpenTime            time.Time `gorm:"column:open_time"`
	OpenPrice                float64   `gorm:"column:open_price"`
	HighPrice                float64   `gorm:"column:high_price"`
	LowPrice                 float64   `gorm:"column:low_price"`
	ClosePrice               float64   `gorm:"column:close_price"`
	Volume                   float64   `gorm:"column:volume"`
	KlineCloseTime           time.Time `gorm:"column:close_time"`
	QuoteAssetVolume         float64   `gorm:"column:quote_asset_volume"`
	NumberOfTrades           int64     `gorm:"column:number_of_trades"`
	TakerBuyBaseAssetVolume  float64   `gorm:"column:taker_buy_base_asset_volume"`
	TakerBuyQuoteAssetVolume float64   `gorm:"column:taker_buy_quote_asset_volume"`
}

type CandleRepository interface {
	Save(candle Candle) error
	SaveMany(candles []Candle) error
}

type CandleDataSource interface {
	Fetch(symbol, interval string, start, end int64, limit int) ([]Candle, error)
}
