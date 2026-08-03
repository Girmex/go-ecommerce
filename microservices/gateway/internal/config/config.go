package config

import (
	"os"

	"github.com/joho/godotenv"
)

type Config struct {
	AppName         string
	AppEnv          string
	HTTPPort        string
	GRPCAuthHost    string
	GRPCCatalogHost string
	GRPCOrderHost   string
	GRPCPaymentHost string
	JWTSecret       string
}

func Load() *Config {
	_ = godotenv.Load(".env.gateway")

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
