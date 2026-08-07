package config

import (
	"os"

	"github.com/joho/godotenv"
)

type Config struct {
	AppName     string
	AppEnv      string
	KAFKABrokers string
}

func Load() *Config {
	_ = godotenv.Load(".env.notification")

	return &Config{
		AppName:      os.Getenv("APP_NAME"),
		AppEnv:       os.Getenv("APP_ENV"),
		KAFKABrokers: os.Getenv("KAFKA_BROKERS"),
	}
}