package handler

import (
	"encoding/json"
	httpResp "github.com/daddydemir/crypto/config/http"
	"github.com/daddydemir/crypto/pkg/token/jwt"
	"github.com/google/uuid"
	"net"
	"net/http"
	"strings"
	"time"
)

type loginResponse struct {
	Username string `json:"username"`
	Token    string `json:"token"`
}

func login(w http.ResponseWriter, r *http.Request) {
	var credentials loginCredentials
	if err := json.NewDecoder(r.Body).Decode(&credentials); err != nil {
		httpResp.WriteJSONError(w, http.StatusBadRequest, httpResp.NewHttpError("invalid request"))
		return
	}
	username, err := authenticationClient.authenticate(credentials)
	if err != nil {
		httpResp.WriteJSONError(w, http.StatusUnauthorized, httpResp.NewHttpError("invalid username or password"))
		return
	}
	sessionID := uuid.NewString()
	token, err := tokenService.GenerateToken(username, sessionID)
	if err != nil {
		httpResp.WriteJSONError(w, http.StatusInternalServerError, httpResp.NewHttpError("token could not be created"))
		return
	}
	if err = authenticationClient.registerSession(sessionRegistration{ID: sessionID, Username: username, Token: token, IPAddress: requestIP(r), UserAgent: r.UserAgent(), ExpiresAt: time.Now().Add(jwt.SessionLifetime)}); err != nil {
		httpResp.WriteJSONError(w, http.StatusBadGateway, httpResp.NewHttpError("session could not be created"))
		return
	}
	httpResp.WriteJSONError(w, http.StatusOK, loginResponse{Username: username, Token: token})
}

func requestIP(r *http.Request) string {
	if value := strings.TrimSpace(strings.Split(r.Header.Get("X-Forwarded-For"), ",")[0]); value != "" {
		return value
	}
	if value := r.Header.Get("X-Real-IP"); value != "" {
		return value
	}
	host, _, err := net.SplitHostPort(r.RemoteAddr)
	if err == nil {
		return host
	}
	return r.RemoteAddr
}
