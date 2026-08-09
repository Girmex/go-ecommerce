package config

import (
	"os"

	"github.com/joho/godotenv"
)

type Config struct {
	AppName         string
	AppEnv          string
	GRPCPort        string
	AuthGRPCAddress string

	KAFKABrokers string

	SMTPHost     string
	SMTPPort     string
	SMTPUsername string
	SMTPPassword string
	SMTPFrom     string
}

func Load() *Config {
	_ = godotenv.Load("microservices/notification/.env.notification")

	return &Config{
		AppName:         os.Getenv("APP_NAME"),
		AppEnv:          os.Getenv("APP_ENV"),
		GRPCPort:        os.Getenv("GRPC_PORT"),
		AuthGRPCAddress: os.Getenv("AUTH_GRPC_ADDRESS"),

		KAFKABrokers: os.Getenv("KAFKA_BROKERS"),

		SMTPHost:     os.Getenv("SMTP_HOST"),
		SMTPPort:     os.Getenv("SMTP_PORT"),
		SMTPUsername: os.Getenv("SMTP_USERNAME"),
		SMTPPassword: os.Getenv("SMTP_PASSWORD"),
		SMTPFrom:     os.Getenv("SMTP_FROM"),
	}
}
