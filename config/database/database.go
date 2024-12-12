package database

import "gorm.io/gorm"

type Database interface {
	Connect()
	Close()
}

var D *gorm.DB

func GetDatabaseService() *gorm.DB {
	var db = PostgresDB{}
	db.Connect()
	return D
}
