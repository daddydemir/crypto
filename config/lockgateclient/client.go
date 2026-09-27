package lockgateclient

import (
	"bytes"
	"context"
	"crypto/rand"
	"encoding/hex"
	"encoding/json"
	"errors"
	"fmt"
	"io"
	"net/http"
	"net/url"
	"strings"
	"time"
)

type Client struct {
	url, token string
	http       *http.Client
}

type accessResponse struct {
	Status     string                       `json:"status"`
	RequestID  string                       `json:"request_id"`
	RetryAfter int                          `json:"retry_after"`
	Secrets    map[string]map[string]string `json:"secrets"`
}

func New(url, token string) *Client {
	return &Client{
		url:   strings.TrimRight(url, "/"),
		token: token,
		http: &http.Client{
			Timeout: 15 * time.Second,
			CheckRedirect: func(*http.Request, []*http.Request) error {
				return http.ErrUseLastResponse // Never forward a token to another origin.
			},
		},
	}
}

func (c *Client) GetConfig(ctx context.Context) (map[string]string, error) {
	instance := make([]byte, 24)
	if _, err := rand.Read(instance); err != nil {
		return nil, errors.New("lockgate: cannot create request identity")
	}

	result, err := c.call(ctx, http.MethodPost, "/api/v1/access", []byte(`{}`), hex.EncodeToString(instance))
	for err == nil {
		switch result.Status {
		case "approved":
			if len(result.Secrets) != 1 {
				return nil, errors.New("lockgate: expected one configured bundle")
			}
			for _, values := range result.Secrets {
				return values, nil
			}
		case "denied":
			return nil, errors.New("lockgate: access denied by administrator")
		case "waiting_approval":
			if result.RequestID == "" || strings.ContainsAny(result.RequestID, "/?#") {
				return nil, errors.New("lockgate: response has an invalid request ID")
			}
			delay := max(result.RetryAfter, 3)
			select {
			case <-ctx.Done():
				return nil, ctx.Err()
			case <-time.After(time.Duration(delay) * time.Second):
			}
			result, err = c.call(ctx, http.MethodGet, "/api/v1/access/"+url.PathEscape(result.RequestID), nil, "")
		default:
			return nil, fmt.Errorf("lockgate: unexpected status %q", result.Status)
		}
	}
	return nil, err
}

func (c *Client) call(ctx context.Context, method, path string, body []byte, idempotencyKey string) (accessResponse, error) {
	var result accessResponse
	base, err := url.Parse(c.url)
	if err != nil || base.Host == "" || (base.Scheme != "https" && base.Scheme != "http") || base.User != nil || base.Path != "" {
		return result, errors.New("lockgate: URL must be an HTTP(S) origin")
	}
	request, err := http.NewRequestWithContext(ctx, method, c.url+path, bytes.NewReader(body))
	if err != nil {
		return result, err
	}
	request.Header.Set("Authorization", "Bearer "+c.token)
	if body != nil {
		request.Header.Set("Content-Type", "application/json")
	}
	if idempotencyKey != "" {
		request.Header.Set("Idempotency-Key", idempotencyKey)
	}

	response, err := c.http.Do(request)
	if err != nil {
		return result, errors.New("lockgate: server unreachable or request timed out")
	}
	defer response.Body.Close()
	if response.StatusCode != http.StatusOK && response.StatusCode != http.StatusAccepted {
		return result, fmt.Errorf("lockgate: request rejected (HTTP %d)", response.StatusCode)
	}
	if err := json.NewDecoder(io.LimitReader(response.Body, 16<<20)).Decode(&result); err != nil {
		return result, errors.New("lockgate: invalid server response")
	}
	return result, nil
}
