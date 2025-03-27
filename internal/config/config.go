package config

import (
	"errors"
	"os"
)

type Config struct {
	DatabaseURL string
	JWTSecret   string
}

func LoadConfig() (*Config, error) {
	dbURL := os.Getenv("DATABASE_URL")
	jwtSecret := os.Getenv("JWT_SECRET")
	if dbURL == "" || jwtSecret == "" {
		return nil, errors.New("DATABASE_URL or JWT_SECRET not set")
	}
	return &Config{
		DatabaseURL: dbURL,
		JWTSecret:   jwtSecret,
	}, nil
}
