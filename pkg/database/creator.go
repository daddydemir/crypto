package database

import (
	"github.com/daddydemir/crypto/config/database"
	"github.com/daddydemir/crypto/pkg/model"
	"log"
)

func CreateTables() {

	err := database.D.AutoMigrate(
		model.DailyModel{},
		model.ExchangeModel{},
		model.WeeklyModel{},
	)
	if err != nil {
		log.Println("::CreateTables:: AutoMigrate err:{}", err)
	}
}
