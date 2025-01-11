package main

import (
	"github.com/daddydemir/crypto/config"
	"github.com/daddydemir/crypto/handler"
	_ "github.com/daddydemir/crypto/pkg/cronjob"
	_ "github.com/daddydemir/dlog"
	"log/slog"
	"net/http"
)

func main() {

	server := &http.Server{
		Addr:    config.Get("PORT"),
		Handler: handler.Route(),
	}

	if config.Get("ENV") == "PROD" {
		if err := server.ListenAndServeTLS(config.Get("CERT_PATH"), config.Get("KEY_PATH")); err != nil {
			slog.Error("ListenAndServeTLS", "error", err)
			panic(err)
		}
	} else {
		if err := server.ListenAndServe(); err != nil {
			slog.Error("ListenAndServe", "error", err)
			panic(err)
		}
	}

}
