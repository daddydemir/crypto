package config

import (
	"github.com/joho/godotenv"
	"log"
	"os"
)

var path string

func init() {
	path = os.Getenv("ENV_PATH")
	if path == "" {
		path = "../config/envs/prod.env"
	}

	println("path: ", path)
}

func Get(key string) string {

	err := godotenv.Load(path)
	if err != nil {
		log.Fatal("Error loading .env file : ", err)
	}
	return os.Getenv(key)
}
