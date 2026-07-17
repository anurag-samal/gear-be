package config

import (
	"os"
	"github.com/joho/godotenv"
)

type Config struct {
	PG_URL       string
	RDS_ADDRESS  string
	RDS_PASSWORD string
	ES_URL       string
	N4J_URL      string
	N4J_USER     string
	N4J_PASSWORD string
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
		N4J_URL:      os.Getenv("N4J_URL"),
		N4J_USER:     os.Getenv("N4J_USER"),
		N4J_PASSWORD: os.Getenv("N4J_PASSWORD"),
	}, nil
}
