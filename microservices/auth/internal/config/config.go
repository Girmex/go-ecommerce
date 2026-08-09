package config

import (
	"os"
	"time"

	"github.com/joho/godotenv"
)

type Config struct {
	AppName string
	AppEnv  string

	GRPCPort string

	PostgresHost     string
	PostgresPort     string
	PostgresUser     string
	PostgresPassword string
	PostgresDatabase string
	JWTSecret        string
	AccessTokenTTL   time.Duration
	RefreshTokenTTL  time.Duration
}

func Load() *Config {

	_ = godotenv.Load("microservices/auth/.env.auth")
	return &Config{
		AppName: os.Getenv("APP_NAME"),
		AppEnv:  os.Getenv("APP_ENV"),

		GRPCPort: os.Getenv("GRPC_PORT"),

		PostgresHost:     os.Getenv("POSTGRES_HOST"),
		PostgresPort:     os.Getenv("POSTGRES_PORT"),
		PostgresUser:     os.Getenv("POSTGRES_USER"),
		PostgresPassword: os.Getenv("POSTGRES_PASSWORD"),
		PostgresDatabase: os.Getenv("POSTGRES_DB"),
		JWTSecret:        os.Getenv("JWT_SECRET"),
		AccessTokenTTL:   time.Hour * 24,
		RefreshTokenTTL:  time.Hour * 24 * 30,
	}
}
