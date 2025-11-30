package binance

import (
	"encoding/json"
	"fmt"
	customHttp "github.com/daddydemir/crypto/config/http"
	"github.com/daddydemir/crypto/pkg/domain/binance"
	"github.com/daddydemir/crypto/pkg/infrastructure/cast"
	"net/http"
	"time"
)

type CandleDataSource struct {
	client *http.Client
}

func NewDataSource() *CandleDataSource {
	return &CandleDataSource{
		client: &http.Client{
			Timeout:   10 * time.Second,
			Transport: customHttp.NewLoggingMiddleware(nil),
		},
	}
}

func (ds *CandleDataSource) Fetch(symbol, interval string, start, end int64, limit int) ([]binance.Candle, error) {
	url := fmt.Sprintf("https://api.binance.com/api/v3/klines?symbol=%s&interval=%s&startTime=%d&endTime=%d&limit=%d", symbol, interval, start, end, limit)
	resp, err := ds.client.Get(url)
	if err != nil {
		return nil, err
	}
	defer resp.Body.Close()

	var raw [][]interface{}
	if err = json.NewDecoder(resp.Body).Decode(&raw); err != nil {
		return nil, err
	}

	candles := make([]binance.Candle, len(raw))
	for i, v := range raw {
		candles[i] = binance.Candle{
			Symbol:                   symbol[:len(symbol)-4],
			KlineOpenTime:            time.UnixMilli(cast.Int64WithoutError(v[0])),
			OpenPrice:                cast.Float64WithoutError(v[1]),
			HighPrice:                cast.Float64WithoutError(v[2]),
			LowPrice:                 cast.Float64WithoutError(v[3]),
			ClosePrice:               cast.Float64WithoutError(v[4]),
			Volume:                   cast.Float64WithoutError(v[5]),
			KlineCloseTime:           time.UnixMilli(cast.Int64WithoutError(v[6])),
			QuoteAssetVolume:         cast.Float64WithoutError(v[7]),
			NumberOfTrades:           cast.Int64WithoutError(v[8]),
			TakerBuyBaseAssetVolume:  cast.Float64WithoutError(v[9]),
			TakerBuyQuoteAssetVolume: cast.Float64WithoutError(v[10]),
		}
	}
	return candles, nil
}
