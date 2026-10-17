package handler

import (
	"bytes"
	"encoding/json"
	"errors"
	"net/http"
	"os"
	"strings"
	"time"

	"github.com/daddydemir/crypto/config"
)

type loginCredentials struct {
	Username string `json:"username"`
	Password string `json:"password"`
}
type sessionRegistration struct {
	ID        string    `json:"id"`
	Username  string    `json:"username"`
	Token     string    `json:"token"`
	IPAddress string    `json:"ipAddress"`
	UserAgent string    `json:"userAgent"`
	ExpiresAt time.Time `json:"expiresAt"`
}
type deviceLinkResponse struct {
	Code      string    `json:"code"`
	ExpiresAt time.Time `json:"expiresAt"`
}

type authIdentity struct {
	Username  string `json:"username"`
	SessionID string `json:"sessionId"`
}

type authClient struct {
	baseURL string
	client  *http.Client
}

type authValidationError struct {
	reason string
	status int
}

func (e *authValidationError) Error() string { return e.reason }

func newAuthClient() *authClient {
	baseURL := config.Get("AUTH_SERVICE_URL")
	if baseURL == "" {
		baseURL = os.Getenv("AUTH_SERVICE_URL")
	}
	return &authClient{baseURL: strings.TrimRight(baseURL, "/"), client: &http.Client{Timeout: 5 * time.Second}}
}

func (c *authClient) authenticate(credentials loginCredentials) (string, error) {
	body, _ := json.Marshal(credentials)
	response, err := c.client.Post(c.baseURL+"/api/v1/credentials/validate", "application/json", bytes.NewReader(body))
	if err != nil {
		return "", err
	}
	defer response.Body.Close()
	if response.StatusCode != http.StatusOK {
		return "", errors.New("invalid username or password")
	}
	var identity authIdentity
	if err = json.NewDecoder(response.Body).Decode(&identity); err != nil {
		return "", err
	}
	return identity.Username, nil
}

func (c *authClient) registerSession(registration sessionRegistration) error {
	body, _ := json.Marshal(registration)
	response, err := c.client.Post(c.baseURL+"/api/v1/sessions", "application/json", bytes.NewReader(body))
	if err != nil {
		return err
	}
	defer response.Body.Close()
	if response.StatusCode != http.StatusCreated {
		return errors.New("could not register session")
	}
	return nil
}

func (c *authClient) createDeviceLink(authorization string) (*deviceLinkResponse, error) {
	request, _ := http.NewRequest(http.MethodPost, c.baseURL+"/api/v1/device-links", nil)
	request.Header.Set("Authorization", authorization)
	response, err := c.client.Do(request)
	if err != nil {
		return nil, err
	}
	defer response.Body.Close()
	if response.StatusCode != http.StatusCreated {
		return nil, errors.New("device link could not be created")
	}
	var link deviceLinkResponse
	if err = json.NewDecoder(response.Body).Decode(&link); err != nil {
		return nil, err
	}
	return &link, nil
}

func (c *authClient) claimDeviceLink(code string) (string, error) {
	body, _ := json.Marshal(map[string]string{"code": code})
	response, err := c.client.Post(c.baseURL+"/api/v1/device-links/claim", "application/json", bytes.NewReader(body))
	if err != nil {
		return "", err
	}
	defer response.Body.Close()
	if response.StatusCode != http.StatusOK {
		return "", errors.New("invalid or expired device link")
	}
	var identity authIdentity
	if err = json.NewDecoder(response.Body).Decode(&identity); err != nil {
		return "", err
	}
	return identity.Username, nil
}

func (c *authClient) validate(authorization string) (*authIdentity, error) {
	if c.baseURL == "" {
		return nil, &authValidationError{reason: "auth_service_not_configured", status: http.StatusServiceUnavailable}
	}
	request, err := http.NewRequest(http.MethodGet, c.baseURL+"/api/v1/sessions/validate", nil)
	if err != nil {
		return nil, err
	}
	request.Header.Set("Authorization", authorization)
	response, err := c.client.Do(request)
	if err != nil {
		return nil, &authValidationError{reason: "auth_service_unreachable", status: http.StatusServiceUnavailable}
	}
	defer response.Body.Close()
	if response.StatusCode != http.StatusOK {
		var payload struct {
			Message string `json:"message"`
		}
		_ = json.NewDecoder(response.Body).Decode(&payload)
		reason := "auth_service_rejected"
		status := response.StatusCode
		switch payload.Message {
		case "missing_token", "session_not_found", "session_revoked", "session_expired", "session_lookup_failed":
			reason = payload.Message
		}
		if response.StatusCode >= http.StatusInternalServerError || reason == "session_lookup_failed" {
			status = http.StatusServiceUnavailable
		}
		return nil, &authValidationError{reason: reason, status: status}
	}
	var identity authIdentity
	if err = json.NewDecoder(response.Body).Decode(&identity); err != nil {
		return nil, &authValidationError{reason: "auth_service_invalid_response", status: http.StatusBadGateway}
	}
	if identity.Username == "" || identity.SessionID == "" {
		return nil, &authValidationError{reason: "auth_service_invalid_identity", status: http.StatusBadGateway}
	}
	return &identity, nil
}
