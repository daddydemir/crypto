package config

import (
	"github.com/joho/godotenv"
	"log"
	"os"
)

func Get(key string) string {
	err := godotenv.Load("../config/envs/prod.env")
	if err != nil {
		log.Fatal("Error loading .env file : ", err)
	}
	return os.Getenv(key)
}
