package client

import (
	"fmt"
	"github.com/daddydemir/crypto/pkg/token/provider"
	"github.com/daddydemir/crypto/pkg/token/strategy"
	"net/http"
	"time"
)

type TokenAwareClient struct {
	BaseURL       string
	HttpClient    *http.Client
	TokenProvider provider.TokenProvider
	Strategy      strategy.TokenStrategy
}

func NewTokenAwareClient(url string, tokenProvider provider.TokenProvider, strategy strategy.TokenStrategy) *TokenAwareClient {
	return &TokenAwareClient{
		TokenProvider: tokenProvider,
		HttpClient:    &http.Client{},
		Strategy:      strategy,
		BaseURL:       url,
	}
}

func (tac *TokenAwareClient) DoGet(endpoint string) (*http.Response, string, error) {
	token, err := tac.TokenProvider.GetValidToken()
	if err != nil {
		return nil, "", err
	}

	req, err := http.NewRequest("GET", tac.BaseURL+endpoint, nil)
	if err != nil {
		return nil, token, err
	}

	if err := tac.Strategy.Apply(req, token); err != nil {
		return nil, token, err
	}

	resp, err := tac.HttpClient.Do(req)
	if err != nil {
		return nil, token, err
	}

	if resp.StatusCode == http.StatusForbidden {
		_ = tac.TokenProvider.MarkTokenAsExpired(token, time.Until(endOfMonth()))
		return nil, token, fmt.Errorf("token expired: %s", token)
	}

	return resp, token, nil
}

func endOfMonth() time.Time {
	now := time.Now()
	year, month := now.Year(), now.Month()
	firstOfNext := time.Date(year, month+1, 1, 0, 0, 0, 0, now.Location())
	return firstOfNext
}
