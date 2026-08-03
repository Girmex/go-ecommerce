package config

import (
	"log"
	"os"

	"github.com/joho/godotenv"
)

type Config struct {
	AppName        string
	AppEnv         string
	HTTPPort       string
	GRPCAuthHost   string
	GRPCCatalogHost string
	GRPCOrderHost  string
	GRPCPaymentHost string
	JWTSecret      string
}

func Load() *Config {
	if err := godotenv.Load(); err != nil {
		log.Println(".env file not found, using system environment variables")
	}

	return &Config{
		AppName:         os.Getenv("APP_NAME"),
		AppEnv:          os.Getenv("APP_ENV"),
		HTTPPort:        os.Getenv("HTTP_PORT"),
		GRPCAuthHost:    os.Getenv("GRPC_AUTH_HOST"),
		GRPCCatalogHost: os.Getenv("GRPC_CATALOG_HOST"),
		GRPCOrderHost:   os.Getenv("GRPC_ORDER_HOST"),
		GRPCPaymentHost: os.Getenv("GRPC_PAYMENT_HOST"),
		JWTSecret:       os.Getenv("JWT_SECRET"),
	}
}