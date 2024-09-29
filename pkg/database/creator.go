package database

import (
	"github.com/daddydemir/crypto/config/database"
	"github.com/daddydemir/crypto/pkg/model"
	"log/slog"
)

func CreateTables() {

	err := database.D.AutoMigrate(
		model.DailyModel{},
		model.ExchangeModel{},
		model.WeeklyModel{},
		model.SignalModel{},
	)
	if err != nil {
		slog.Error("CreateTables.AutoMigrate", "error", err)
	}
}
