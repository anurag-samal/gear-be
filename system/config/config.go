package config

import (
	"github.com/joho/godotenv"
	"os"
)

type Config struct {
	PG_URL        string
	RDS_ADDRESS   string
	RDS_PASSWORD  string
	ES_URL        string
	N4J_URL       string
	N4J_USER      string
	N4J_PASSWORD  string
	COOKIE_SECURE bool
	FRONTEND_URL  string
	JWT_SECRET    string
	GOOGLE_CLIENT_ID string
	GOOGLE_CLIENT_SECRET string
	GOOGLE_REDIRECT_URL string
}

func Load() (*Config, error) {
	err := godotenv.Load()
	if err != nil {
		return nil, err
	}
	return &Config{
		PG_URL:        os.Getenv("PG_USER"),
		RDS_ADDRESS:   os.Getenv("RDS_ADDRESS"),
		RDS_PASSWORD:  os.Getenv("RDS_PASSWORD"),
		ES_URL:        os.Getenv("ES_URL"),
		N4J_URL:       os.Getenv("N4J_URL"),
		N4J_USER:      os.Getenv("N4J_USER"),
		N4J_PASSWORD:  os.Getenv("N4J_PASSWORD"),
		COOKIE_SECURE: os.Getenv("COOKIE_SECURE") == "true",
		FRONTEND_URL:  os.Getenv("FRONTEND_URL"),
		JWT_SECRET:    os.Getenv("JWT_SECRET"),
		GOOGLE_CLIENT_ID: os.Getenv("GOOGLE_CLIENT_ID"),
		GOOGLE_CLIENT_SECRET: os.Getenv("GOOGLE_CLIENT_SECRET"),
		GOOGLE_REDIRECT_URL: os.Getenv("GOOGLE_REDIRECT_URL"),
	}, nil
}
