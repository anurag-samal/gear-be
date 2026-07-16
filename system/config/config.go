package config

import (
	"github.com/joho/godotenv"
	"os"
)

type Config struct {
	PG_URL       string
	RDS_ADDRESS  string
	RDS_PASSWORD string
	ES_URL       string
}

func Load() (*Config, error) {
	err := godotenv.Load()
	if err != nil {
		return nil, err
	}
	return &Config{
		PG_URL:       os.Getenv("PG_USER"),
		RDS_ADDRESS:  os.Getenv("RDS_ADDRESS"),
		RDS_PASSWORD: os.Getenv("RDS_PASSWORD"),
		ES_URL:       os.Getenv("ES_URL"),
	}, nil
}
