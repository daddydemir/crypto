package database

import (
	"github.com/daddydemir/crypto/config/database"
	"github.com/daddydemir/crypto/config/log"
	"github.com/daddydemir/crypto/pkg/model"
)

func CreateTables() {

	err := database.D.AutoMigrate(
		model.DailyModel{},
		model.ExchangeModel{},
		model.WeeklyModel{},
	)
	if err != nil {
		log.Errorln("::CreateTables:: AutoMigrate err:{}", err)
	}
}
