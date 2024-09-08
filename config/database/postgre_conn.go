package database

import (
	"github.com/daddydemir/crypto/config"
	"gorm.io/driver/postgres"
	"gorm.io/gorm"
	"log/slog"
)

type PostgresDB struct {
}

func (d *PostgresDB) Connect() {
	dsn := config.Get("POSTGRE_DSN")
	db, err := gorm.Open(postgres.Open(dsn), &gorm.Config{})
	if err != nil {
		slog.Error("Connect:gorm.Open", "error", err)
		panic(err)
	}
	D = db
	slog.Info("Connect:gorm.Open", "message", "connection was successful")
}

func (d *PostgresDB) Close() {
	// not implemented...
}
