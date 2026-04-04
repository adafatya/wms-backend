package config

import (
	"log"
	"os"

	"github.com/joho/godotenv"
)

type Config struct {
	DBDriver      string
	DBSource      string
	ServerAddress string
}

func LoadConfig() (Config, error) {
	err := godotenv.Load()
	if err != nil {
		log.Println("No .env file found, reading from environment variables")
	}

	config := Config{
		DBDriver:      os.Getenv("DB_DRIVER"),
		DBSource:      os.Getenv("DB_SOURCE"),
		ServerAddress: ":" + os.Getenv("PORT"),
	}

	return config, nil
}
