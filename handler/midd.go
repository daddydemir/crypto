package handler

import (
	"log/slog"
	"net/http"
)

func setJSONContentType(next http.Handler) http.Handler {
	return http.HandlerFunc(func(w http.ResponseWriter, r *http.Request) {
		w.Header().Set("Content-Type", "application/json")
		next.ServeHTTP(w, r)
	})
}

func setLogging(next http.Handler) http.Handler {
	return http.HandlerFunc(func(w http.ResponseWriter, r *http.Request) {
		slog.Info("endpoint hitted", "url", r.URL.RequestURI(), "method", r.Method, "IP", r.RemoteAddr)
		next.ServeHTTP(w, r)
	})
}
