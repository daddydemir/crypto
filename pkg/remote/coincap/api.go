package coincap

import (
	"encoding/json"
	"fmt"
	"github.com/daddydemir/crypto/pkg/remote/coincap/client"
	"github.com/daddydemir/crypto/pkg/token/provider"
	"github.com/daddydemir/crypto/pkg/token/strategy"
	"log/slog"
)

var (
	baseUrl              = "https://rest.coincap.io"
	allCoins             = "/v3/assets"
	priceHistoryWithId   = "/v3/assets/%v/history?interval=d1"
	priceHistoryWithTime = "/v3/assets/%v/history?interval=d1&start=%d&end=%d"
)

type Client struct {
}

func NewClient() *Client {
	return &Client{}
}

func (c *Client) ListCoins() (error, []Coin) {

	var data Data[Coin]

	strategy := strategy.QueryTokenStrategy{ParamName: "apiKey"}

	client := client.NewTokenAwareClient(baseUrl, provider.NewRedisTokenProvider(), strategy)

	resp, _, err := client.DoGet(allCoins)
	if err != nil {
		slog.Error("ListCoins:http.Get", "url", allCoins, "err", err)
		return err, nil
	} else {
		slog.Info("ListCoins:http.Get", "url", allCoins, "statusCode", resp.StatusCode)
	}

	err = json.NewDecoder(resp.Body).Decode(&data)
	if err != nil {
		slog.Error("ListCoins:json.Decode", "err", err)
		return err, nil
	}

	return nil, data.Data
}

func (c *Client) HistoryWithId(s string) (error, []History) {

	var data Data[History]

	endpoint := fmt.Sprintf(priceHistoryWithId, s)
	url := endpoint
	strategy := strategy.HeaderTokenStrategy{}
	client := client.NewTokenAwareClient(baseUrl, provider.NewRedisTokenProvider(), strategy)

	resp, _, err := client.DoGet(endpoint)
	if err != nil {
		slog.Error("HistoryWithId:http.Get", "url", url, "err", err)
		return err, nil
	} else {
		slog.Info("HistoryWithId:http.Get", "url", url, "statusCode", resp.StatusCode)
	}

	err = json.NewDecoder(resp.Body).Decode(&data)
	if err != nil {
		slog.Error("HistoryWithId:json.Decode", "err", err)
		return err, nil
	}

	return nil, data.Data
}

func (c *Client) HistoryWithTime(s string, start, end int64) (error, []History) {
	var data Data[History]

	endpoint := fmt.Sprintf(priceHistoryWithTime, s, start/1_000_000, end/1_000_000)
	url := endpoint

	strategy := strategy.HeaderTokenStrategy{}
	client := client.NewTokenAwareClient(baseUrl, provider.NewRedisTokenProvider(), strategy)
	resp, _, err := client.DoGet(endpoint)

	if err != nil {
		slog.Error("HistoryWithTime:http.Get", "url", url, "err", err)
		return err, nil
	} else {
		slog.Info("HistoryWithTime:http.Get", "url", url, "statusCode", resp.StatusCode)
	}

	err = json.NewDecoder(resp.Body).Decode(&data)
	if err != nil {
		slog.Error("HistoryWithTime:json.Decode", "error", err, "data", resp.Body)
		return err, nil
	}
	return nil, data.Data
}
