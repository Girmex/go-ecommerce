package config

import (
	"log"
	"os"

	"github.com/joho/godotenv"
)

type Config struct {
	HTTPPort     string
	DatabaseURL  string
	JWTSecret    string
}

func Load() *Config {
	_ = godotenv.Load()

	cfg := &Config{
		HTTPPort:    os.Getenv("HTTP_PORT"),
		DatabaseURL: os.Getenv("DATABASE_URL"),
		JWTSecret: os.Getenv("JWT_SECRET"),
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