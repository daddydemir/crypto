package database

import (
	"github.com/daddydemir/crypto/config"
	"gorm.io/driver/mysql"
	"gorm.io/gorm"
	"log/slog"
)

type MySQLDB struct {
}

func (d *MySQLDB) Connect() {
	dsn := config.Get("DSN")
	db, err := gorm.Open(mysql.Open(dsn), &gorm.Config{})
	if err != nil {
		slog.Error("Connect:gorm.Open", "error", err)
	}
	D = db
	slog.Info("Connect:gorm.Open", "message", "connection was successful")
}

func (d *MySQLDB) Close() {
	// not implemented ...
}
