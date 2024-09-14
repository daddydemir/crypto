package main

import (
	"github.com/daddydemir/crypto/config"
	"github.com/daddydemir/crypto/config/database"
	_ "github.com/daddydemir/crypto/config/log/dlog"
	"github.com/daddydemir/crypto/handler"
	"github.com/daddydemir/crypto/pkg/cronjob"
	"log/slog"
	"net/http"
)

func main() {

	db := database.PostgresDB{}
	db.Connect()

	cronjob.StartCronJob()

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
