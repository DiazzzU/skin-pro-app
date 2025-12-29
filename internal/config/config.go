package config

import (
	"os"

	"github.com/joho/godotenv"
)

type GlobalConfig struct {
	PgDSN     string
	Port      string
	JWTSecret string
}

func Load() (*GlobalConfig, error) {
	_ = godotenv.Load()
	return &GlobalConfig{
		PgDSN:     os.Getenv("PG_DSN"),
		Port:      os.Getenv("APP_PORT"),
		JWTSecret: os.Getenv("JWT_SECRET"),
	}, nil
}
