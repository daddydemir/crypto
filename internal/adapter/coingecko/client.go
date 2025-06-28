package coingecko

import (
	"encoding/json"
	"errors"
	"github.com/daddydemir/crypto/internal/domain/model"
	"log/slog"
	"net/http"
)

const topHundredURL = "https://api.coingecko.com/api/v3/coins/markets?vs_currency=usd&order=market_cap_desc&per_page=100&page=1&sparkline=false"

type HttpGeckoClient struct {
	httpClient *http.Client
}

func NewHttpGeckoClient(client *http.Client) *HttpGeckoClient {
	if client == nil {
		client = &http.Client{}
	}
	return &HttpGeckoClient{httpClient: client}
}

func (g *HttpGeckoClient) GetTopHundred() ([]model.ExchangeModel, error) {
	req, err := http.NewRequest(http.MethodGet, topHundredURL, nil)
	if err != nil {
		slog.Error("coingecko.NewRequest", "url", topHundredURL, "error", err)
		return nil, err
	}

	resp, err := g.httpClient.Do(req)
	if err != nil {
		slog.Error("coingecko.client.Do", "url", topHundredURL, "error", err)
		return nil, err
	}
	defer resp.Body.Close()

	if resp.StatusCode != http.StatusOK {
		return nil, errors.New("coingecko: unexpected status code")
	}

	var exchanges []model.ExchangeModel
	if err = json.NewDecoder(resp.Body).Decode(&exchanges); err != nil {
		slog.Error("coingecko.json.Unmarshal", "error", err)
		return nil, err
	}

	return exchanges, nil
}

func (g *HttpGeckoClient) GetTopHundredDaily() ([]model.DailyModel, error) {
	req, err := http.NewRequest(http.MethodGet, topHundredURL, nil)
	if err != nil {
		slog.Error("coingecko.NewRequest", "url", topHundredURL, "error", err)
		return nil, err
	}

	resp, err := g.httpClient.Do(req)
	if err != nil {
		slog.Error("coingecko.client.Do", "url", topHundredURL, "error", err)
		return nil, err
	}
	defer resp.Body.Close()

	if resp.StatusCode != http.StatusOK {
		return nil, errors.New("coingecko: unexpected status code")
	}

	var dailies []model.DailyModel
	if err = json.NewDecoder(resp.Body).Decode(&dailies); err != nil {
		slog.Error("coingecko.json.Unmarshal", "error", err)
		return nil, err
	}

	return dailies, nil
}
