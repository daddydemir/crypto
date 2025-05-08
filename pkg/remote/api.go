package remote

import (
	"fmt"
	"io"
	"net/http"
)

type ApiClient struct {
	BaseURL string
	Token   string
	Client  *http.Client
}

func NewApiClient(baseUrl, token string) *ApiClient {
	return &ApiClient{
		BaseURL: baseUrl,
		Token:   token,
		Client:  &http.Client{},
	}
}

func (api *ApiClient) DoRequest(method, endpoint string, body io.Reader) (*http.Response, error) {

	req, err := http.NewRequest(method, api.BaseURL+endpoint, body)
	if err != nil {
		return nil, fmt.Errorf("request oluşturulamadı: %w", err)
	}

	req.Header.Set("Authorization", "Bearer "+api.Token)
	req.Header.Set("Content-Type", "Application/json")

	resp, err := api.Client.Do(req)
	if err != nil {
		return nil, fmt.Errorf("istek gönderilemedi: %w", err)
	}
	return resp, nil
}
