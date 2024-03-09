package database

import (
	"github.com/daddydemir/crypto/config"
	"github.com/daddydemir/crypto/config/log"
	"gorm.io/driver/mysql"
	"gorm.io/gorm"
)

type MySQLDB struct {
}

func (d *MySQLDB) Connect() {
	dsn := config.Get("DSN")
	db, err := gorm.Open(mysql.Open(dsn), &gorm.Config{})
	if err != nil {
		log.Fatal("Database connection error: ", err)
	}
	D = db
	log.Infoln("Connected to database")
}

func (d *MySQLDB) Close() {
	// not implemented ...
}
