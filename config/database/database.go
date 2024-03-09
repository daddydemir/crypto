package database

import "gorm.io/gorm"

type Database interface {
	Connect()
	Close()
}

var D *gorm.DB
