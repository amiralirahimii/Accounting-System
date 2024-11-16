package config

import (
	"accountingsystem/internal/constants"
	"log"
	"os"

	"github.com/joho/godotenv"
)

func InitConfig() error {
	err := godotenv.Load()
	if err != nil {
		log.Fatalf("Error loading .env file: %v", err)
		return err
	}
	return nil
}

func GetEnv(key string) (string, error) {
	if value, exists := os.LookupEnv(key); exists {
		return value, nil
	}
	return "", constants.ErrEnvNotFound
}
