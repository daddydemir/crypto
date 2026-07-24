package http

import (
	"encoding/json"
	"net/http"
)

func WriteJSONError(w http.ResponseWriter, status int, payload any) {
	w.Header().Set("Content-Type", "application/json")
	w.WriteHeader(status)
	json.NewEncoder(w).Encode(payload)
}

type HttpError struct {
	Message string `json:"message"`
}

func NewHttpError(message string) HttpError {
	return HttpError{Message: message}
}
