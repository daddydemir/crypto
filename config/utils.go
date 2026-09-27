package config

import (
	"context"
	"log"
	"os"
	"time"

	lockgate "github.com/daddydemir/crypto/config/lockgateclient"
	"github.com/joho/godotenv"
)

var path string

var secretValues map[string]string

func init() {
	path = ".env"
	err := godotenv.Load(path)
	if err != nil {
		println("Error loading .env file", err)
	}

	client := lockgate.New(os.Getenv("LOCKGATE_URL"), os.Getenv("LOCKGATE_TOKEN"))
	ctx, cancel := context.WithTimeout(context.Background(), 10*time.Minute)
	defer cancel()

	values, err := client.GetConfig(ctx)
	if err != nil {
		log.Fatal(err)
	}
	secretValues = values
}

func Get(key string) string {
	return secretValues[key]
}
