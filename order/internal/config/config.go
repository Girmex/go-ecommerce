package config

import (
	"log"
	"os"

	"github.com/joho/godotenv"
)

type Config struct {
	DatabaseURL string
	HTTPPort    string
	JWTSecret   string
	UserServiceURL   string
	ProductServiceURL string
	PaymentServiceURL string

}

func Load() *Config {
	_ = godotenv.Load()

	cfg := &Config{
		DatabaseURL: os.Getenv("DATABASE_URL"),
		HTTPPort:    os.Getenv("HTTP_PORT"),
		JWTSecret:   os.Getenv("JWT_SECRET"),
		UserServiceURL:   os.Getenv("USER_SERVICE_URL"),
		ProductServiceURL: os.Getenv("PRODUCT_SERVICE_URL"),
		PaymentServiceURL: os.Getenv("PAYMENT_SERVICE_URL"),
	}

	if cfg.DatabaseURL == "" {
		log.Fatal("DATABASE_URL is required")
	}

	if cfg.HTTPPort == "" {
		cfg.HTTPPort = "8082"
	}

	if cfg.UserServiceURL == "" {
		log.Fatal("USER_SERVICE_URL is required")
	}

	if cfg.ProductServiceURL == "" {
		log.Fatal("PRODUCT_SERVICE_URL is required")
	}

	if cfg.PaymentServiceURL == "" {
		log.Fatal("PAYMENT_SERVICE_URL is required")
	}

	return cfg
}