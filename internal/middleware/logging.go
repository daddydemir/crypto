package middleware

import (
	"fmt"
	"log/slog"
	"net/http"
	"time"
)

func Logging(next http.Handler) http.Handler {
	return http.HandlerFunc(func(w http.ResponseWriter, r *http.Request) {
		start := time.Now()

		slog.Info("HTTP request started",
			"method", r.Method,
			"url", r.URL.RequestURI(),
			"remote", r.RemoteAddr,
			"start", start.Format(time.RFC3339),
		)

		next.ServeHTTP(w, r)

		duration := time.Since(start)
		ms := float64(duration.Microseconds()) / 1000.0

		slog.Info("HTTP request completed",
			"method", r.Method,
			"url", r.URL.RequestURI(),
			"status", http.StatusOK, // opsiyonel: ResponseWriter wrapper'ı ile gerçek status code alınabilir
			"duration_ms", fmt.Sprintf("%.2fms", ms),
		)
	})
}
