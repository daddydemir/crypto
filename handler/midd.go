package handler

import (
	"context"
	"errors"
	"log/slog"
	"net/http"
	"slices"

	httpResp "github.com/daddydemir/crypto/config/http"
)

var ignoredEndpoints = []string{
	"/health",
}

func setJSONContentType(next http.Handler) http.Handler {
	return http.HandlerFunc(func(w http.ResponseWriter, r *http.Request) {
		w.Header().Set("Content-Type", "application/json")
		next.ServeHTTP(w, r)
	})
}

func setLogging(next http.Handler) http.Handler {
	return http.HandlerFunc(func(w http.ResponseWriter, r *http.Request) {
		if slices.Contains(ignoredEndpoints, r.URL.Path) {
			next.ServeHTTP(w, r)
			return
		}
		ip := r.Header.Get("X-Forwarded-For")
		if ip == "" {
			ip = r.Header.Get("X-Real-IP")
		}
		if ip == "" {
			ip = r.RemoteAddr
		}
		slog.Info("endpoint invoked",
			"url", r.URL.RequestURI(),
			"method", r.Method,
			"IP", ip,
		)
		next.ServeHTTP(w, r)
	})
}

func auth(next http.Handler) http.Handler {
	return http.HandlerFunc(func(w http.ResponseWriter, r *http.Request) {
		identity, err := authenticationClient.validate(r.Header.Get("Authorization"))
		if err != nil {
			reason := "auth_validation_failed"
			status := http.StatusUnauthorized
			var validationError *authValidationError
			if errors.As(err, &validationError) {
				reason = validationError.reason
				if validationError.status != 0 {
					status = validationError.status
				}
			}
			w.Header().Set("X-Auth-Reason", reason)
			slog.Warn("authentication rejected", "method", r.Method, "path", r.URL.Path, "reason", reason, "status", status)
			message := "Unauthorized"
			if status >= http.StatusInternalServerError {
				message = "Authentication service temporarily unavailable"
			}
			httpResp.WriteJSONError(w, status, httpResp.NewHttpError(message))
			return
		}
		ctx := context.WithValue(r.Context(), "username", identity.Username)
		ctx = context.WithValue(ctx, "sessionID", identity.SessionID)
		next.ServeHTTP(w, r.WithContext(ctx))
	})
}
