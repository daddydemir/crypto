package handler

import (
	"encoding/json"
	httpResp "github.com/daddydemir/crypto/config/http"
	"github.com/google/uuid"
	"net/http"
	"time"

	"github.com/daddydemir/crypto/pkg/token/jwt"
)

func createDeviceLink(w http.ResponseWriter, r *http.Request) {
	link, err := authenticationClient.createDeviceLink(r.Header.Get("Authorization"))
	if err != nil {
		httpResp.WriteJSONError(w, http.StatusBadGateway, httpResp.NewHttpError(err.Error()))
		return
	}
	httpResp.WriteJSONError(w, http.StatusCreated, link)
}
func claimDeviceLink(w http.ResponseWriter, r *http.Request) {
	var request struct {
		Code string `json:"code"`
	}
	if json.NewDecoder(r.Body).Decode(&request) != nil {
		httpResp.WriteJSONError(w, http.StatusBadRequest, httpResp.NewHttpError("invalid request"))
		return
	}
	username, err := authenticationClient.claimDeviceLink(request.Code)
	if err != nil {
		httpResp.WriteJSONError(w, http.StatusUnauthorized, httpResp.NewHttpError(err.Error()))
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
