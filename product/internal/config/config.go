package config

import (
	"log"
	"os"

	"github.com/joho/godotenv"
)

type Config struct {
	DatabaseURL string
	HTTPPort    string
	JwtSecret   string
}

func Load() *Config {
	_ = godotenv.Load()

	cfg := &Config{
		DatabaseURL: os.Getenv("DATABASE_URL"),
		HTTPPort:    os.Getenv("HTTP_PORT"),
		JwtSecret:   os.Getenv("JWT_SECRET"),
	}

	if cfg.DatabaseURL == "" {
		log.Fatal("DATABASE_URL is required")
	}

	if cfg.HTTPPort == "" {
		cfg.HTTPPort = "8081"
	}

	if cfg.JwtSecret == "" {
		log.Fatal("JWT_SECRET is required")
	}

	return cfg
}
