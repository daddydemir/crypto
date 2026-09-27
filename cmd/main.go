package main

import (
	"log/slog"
	"net/http"
	"time"

	"github.com/daddydemir/crypto/config"
	cch "github.com/daddydemir/crypto/config/cache"
	"github.com/daddydemir/crypto/handler"
	bnnc "github.com/daddydemir/crypto/pkg/remote/binance"
	_ "github.com/daddydemir/dlog"
)

func main() {

	go bnnc.NewClient(config.Get("WS_URL"), cch.GetRedisClient()).Fetch()
	go handler.ListenAndBroadcast(cch.GetRedisClient())
	server := &http.Server{
		ReadHeaderTimeout: 3 * time.Second,
		Addr:              config.Get("PORT"),
		Handler:           handler.Route(),
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
