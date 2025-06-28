package coincapadapter

import (
	"encoding/json"
	"fmt"
	"github.com/daddydemir/crypto/internal/domain/model"
	"github.com/daddydemir/crypto/internal/port"
	"github.com/daddydemir/crypto/pkg/remote/coincap/client"
	"log/slog"
)

const baseUrl = "https://rest.coincap.io"

type RealCoinCapClient struct {
	baseURL       string
	clientFactory func(TokenStrategy) client.TokenAwareClient
}

func NewRealCoinCapClient(
	baseURL string,
	clientFactory func(TokenStrategy) client.TokenAwareClient,
) port.CoinCapAPI {
	return &RealCoinCapClient{
		baseURL:       baseURL,
		clientFactory: clientFactory,
	}
}

func (r *RealCoinCapClient) ListCoins() (error, []model.Coin) {
	url := "/v3/assets"
	var data struct {
		Data []model.Coin `json:"data"`
	}

	cli := r.clientFactory(QueryStrategy)
	//cli := client.NewTokenAwareClient(baseUrl, provider.NewRedisTokenProvider(), strategy.QueryTokenStrategy{ParamName: "apiKey"})

	resp, _, err := cli.DoGet(url)
	if err != nil {
		slog.Error("ListCoins", "url", url, "err", err)
		return err, nil
	}
	defer resp.Body.Close()

	if err = json.NewDecoder(resp.Body).Decode(&data); err != nil {
		slog.Error("ListCoins:decode", "err", err)
		return err, nil
	}

	return nil, data.Data
}

func (r *RealCoinCapClient) HistoryWithId(id string) (error, []model.History) {
	url := fmt.Sprintf("/v3/assets/%s/history?interval=d1", id)
	var data struct {
		Data []model.History `json:"data"`
	}

	cli := r.clientFactory(HeaderStrategy)
	//cli := client.NewTokenAwareClient(baseUrl, provider.NewRedisTokenProvider(), strategy.HeaderTokenStrategy{})

	resp, _, err := cli.DoGet(url)
	if err != nil {
		slog.Error("HistoryWithId", "url", url, "err", err)
		return err, nil
	}
	defer resp.Body.Close()

	if err = json.NewDecoder(resp.Body).Decode(&data); err != nil {
		slog.Error("HistoryWithId:decode", "err", err)
		return err, nil
	}

	return nil, data.Data
}

func (r *RealCoinCapClient) HistoryWithTime(id string, start, end int64) (error, []model.History) {
	url := fmt.Sprintf("/v3/assets/%s/history?interval=d1&start=%d&end=%d", id, start/1_000_000, end/1_000_000)

	var data struct {
		Data []model.History `json:"data"`
	}

	cli := r.clientFactory(HeaderStrategy)
	//cli := client.NewTokenAwareClient(baseUrl, provider.NewRedisTokenProvider(), strategy.HeaderTokenStrategy{})

	resp, _, err := cli.DoGet(url)
	if err != nil {
		slog.Error("HistoryWithTime", "url", url, "err", err)
		return err, nil
	}
	defer resp.Body.Close()

	if err = json.NewDecoder(resp.Body).Decode(&data); err != nil {
		slog.Error("HistoryWithTime:decode", "err", err)
		return err, nil
	}

	return nil, data.Data
}
