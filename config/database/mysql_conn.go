package database

import (
	"github.com/daddydemir/crypto/config"
	"gorm.io/driver/mysql"
	"gorm.io/gorm"
	"log"
)

var D *gorm.DB

func InitMySQLConnect() {
	dsn := config.Get("DSN")
	db, err := gorm.Open(mysql.Open(dsn), &gorm.Config{})
	if err != nil {
		log.Fatal("Database connection error: ", err)
	}
	D = db
}
