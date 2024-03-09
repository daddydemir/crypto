package database

import (
	"github.com/daddydemir/crypto/config"
	"github.com/daddydemir/crypto/config/log"
	"gorm.io/driver/postgres"
	"gorm.io/gorm"
)

type PostgresDB struct {
}

func (d *PostgresDB) Connect() {
	dsn := config.Get("POSTGRE_DSN") // todo: unresolved!c
	db, err := gorm.Open(postgres.Open(dsn), &gorm.Config{})
	if err != nil {
		log.Fatal("Database connection error: ", err)
	}
	D = db
	log.Infoln("Connected to database")

}

func (d *PostgresDB) Close() {
	// not implemented...
}
