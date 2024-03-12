package main

import (
	"github.com/daddydemir/crypto/config"
	"github.com/daddydemir/crypto/config/database"
	"github.com/daddydemir/crypto/config/log"
	"github.com/daddydemir/crypto/handler"
	"github.com/daddydemir/crypto/pkg/cronjob"
	"net/http"
)

func main() {

	log.InitLogger()
	db := database.PostgresDB{}
	db.Connect()

	config.NewRabbitMQ()

	cronjob.StartCronJob()

	server := &http.Server{
		Addr:    config.Get("PORT"),
		Handler: handler.Route(),
	}
	if config.Get("ENV") == "PROD" {
		err := server.ListenAndServeTLS(config.Get("CERT_PATH"), config.Get("KEY_PATH"))
		if err != nil {
			log.Fatal("::ListenAndServeTLS:: err:{}", err)
		}
	} else {
		err := server.ListenAndServe()
		if err != nil {
			log.Fatal("::ListenAndServe:: err:{}", err)
			fileErr := log.LogFile.Close()
			if fileErr != nil {
				log.Errorln(err)
			}
		}
	}

}
