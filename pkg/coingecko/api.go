package coingecko

import (
	"encoding/json"
	"github.com/daddydemir/crypto/pkg/adapter"
	"io/ioutil"
	"log/slog"
	"net/http"
)

const (
	topHundred string = "https://api.coingecko.com/api/v3/coins/markets?vs_currency=usd&order=market_cap_desc&per_page=100&page=1&sparkline=false"
)

func GetTopHundred() []adapter.Adapter {

	req, err := http.NewRequest(http.MethodGet, topHundred, nil)
	if err != nil {
		slog.Error("GetTopHundred:http.NewRequest", "url", topHundred, "error", err)
		return nil
	}

	client := &http.Client{}
	resp, err := client.Do(req)
	if err != nil {
		slog.Error("GetTopHundred:client.Do", "url", topHundred, "error", err)
		return nil
	}
	defer resp.Body.Close()

	body, err := ioutil.ReadAll(resp.Body)
	if err != nil {
		slog.Error("GetTopHundred:ioutil.ReadAll", "error", err)
		return nil
	}
	var adapts []adapter.Adapter
	err = json.Unmarshal(body, &adapts)
	if err != nil {
		slog.Error("GetTopHundred:json.Unmarshal", "error", err)
	}
	return adapts
}
