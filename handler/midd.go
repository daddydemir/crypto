package handler

import (
	"context"
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
			httpResp.WriteJSONError(w, http.StatusUnauthorized, httpResp.NewHttpError("Unauthorized"))
			return
		}
		ctx := context.WithValue(r.Context(), "username", identity.Username)
		ctx = context.WithValue(ctx, "sessionID", identity.SessionID)
		next.ServeHTTP(w, r.WithContext(ctx))
	})
}
