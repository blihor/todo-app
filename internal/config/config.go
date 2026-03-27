package config

import (
	"fmt"
	"os"

	"github.com/joho/godotenv"
)

type Config struct {
	Port      string
	SecretJwt string
	DBConnStr string
}

func LoadConfig() (*Config, error) {
	err := godotenv.Load()
	if err != nil {
		return nil, fmt.Errorf("Failed to load .env file")
	}

	return &Config{
		Port:      os.Getenv("PORT"),
		SecretJwt: os.Getenv("SECRET_JWT"),
		DBConnStr: os.Getenv("DB_CONN_STR"),
	}, nil
}
