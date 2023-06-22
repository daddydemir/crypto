package main

import (
	"github.com/daddydemir/crypto/config"
	"github.com/daddydemir/crypto/config/database"
	"github.com/daddydemir/crypto/handler"
	"log"
	"net/http"
)

func main() {

	database.InitMySQLConnect()

	server := &http.Server{
		Addr:    config.Get("PORT"),
		Handler: handler.Route(),
	}
	if config.Get("ENV") == "PROD" {
		err := server.ListenAndServeTLS(config.Get("CERT_PATH"), config.Get("KEY_PATH"))
		if err != nil {
			log.Fatal("::ListenAndServe:: err:{}", err)
		}
	} else {
		err := server.ListenAndServe()
		if err != nil {
			log.Fatal("::ListenAndServe:: err:{}", err)
		}
	}

}
