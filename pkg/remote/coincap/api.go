package coincap

import (
	"encoding/json"
	"fmt"
	"github.com/daddydemir/crypto/config"
	"github.com/daddydemir/crypto/pkg/remote"
	"log/slog"
)

var (
	baseUrl              = "https://rest.coincap.io"
	allCoins             = "/v3/assets?apiKey=%v"
	priceHistoryWithId   = "/v3/assets/%v/history?interval=d1"
	priceHistoryWithTime = "/v3/assets/%v/history?interval=d1&start=%d&end=%d"
	token                = config.Get("API_TOKEN")
)

type Client struct {
}

func NewClient() *Client {
	return &Client{}
}

// todo refactor error
func (c *Client) ListCoins() (error, []Coin) {

	var data Data[Coin]
	endpoint := fmt.Sprintf(allCoins, token)
	client := remote.NewApiClient(baseUrl, token)

	resp, err := client.DoRequest("GET", endpoint, nil)
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
	client := remote.NewApiClient(baseUrl, token)

	resp, err := client.DoRequest("GET", endpoint, nil)
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
	client := remote.NewApiClient(baseUrl, token)
	resp, err := client.DoRequest("GET", endpoint, nil)
	//resp, err := http.Get(url)
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
