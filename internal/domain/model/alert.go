package model

import "time"

type Alert struct {
	ID         uint      `gorm:"primaryKey"`
	Coin       string    `gorm:"type:varchar(10);not null"`
	Price      float32   `gorm:"not null"`
	IsAbove    bool      `gorm:"not null"`
	CreateDate time.Time `gorm:"autoCreateTime"`
	IsActive   bool      `gorm:"default:true"`
}
