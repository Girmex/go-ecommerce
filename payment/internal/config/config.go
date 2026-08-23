package config

import (
	"log"
	"os"

	"github.com/joho/godotenv"
)

type Config struct {
	HTTPPort     string
	DatabaseURL  string
}

func Load() *Config {
	_ = godotenv.Load()

	cfg := &Config{
		HTTPPort:    getEnv("HTTP_PORT", "8084"),
		DatabaseURL: getEnv(
			"DATABASE_URL",
			"postgres://postgres:postgres@localhost:5434/payment_db",
		),
	}

	return cfg
}

func getEnv(key, fallback string) string {
	value := os.Getenv(key)

	if value == "" {
		log.Printf(
			"%s not set, using default value",
			key,
		)
		return fallback
	}

	return value
}