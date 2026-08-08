package config

import (
	"os"
	"strings"

	"github.com/joho/godotenv"
)

type Config struct {
	AppName          string
	AppEnv           string
	KAFKABrokers     []string
	ConsumerGroupID  string
	EmailProvider    string
	SMTPHost         string
	SMTPPort         string
	SMTPUsername     string
	SMTPPassword     string
	SMTPFrom         string
	SMSProvider      string
	TwilioAccountSID string
	TwilioAuthToken  string
	TwilioFromPhone  string
}

func Load() *Config {
	_ = godotenv.Load(".env.notification")

	brokersStr := getEnv("KAFKA_BROKERS", "localhost:29092")
	brokers := strings.Split(brokersStr, ",")
	for i := range brokers {
		brokers[i] = strings.TrimSpace(brokers[i])
	}

	return &Config{
		AppName:          getEnv("APP_NAME", "Notification Service"),
		AppEnv:           getEnv("APP_ENV", "development"),
		KAFKABrokers:     brokers,
		ConsumerGroupID:  getEnv("CONSUMER_GROUP_ID", "notification-group"),
		EmailProvider:    getEnv("EMAIL_PROVIDER", "logger"),
		SMTPHost:         getEnv("SMTP_HOST", "localhost"),
		SMTPPort:         getEnv("SMTP_PORT", "587"),
		SMTPUsername:     os.Getenv("SMTP_USERNAME"),
		SMTPPassword:     os.Getenv("SMTP_PASSWORD"),
		SMTPFrom:         getEnv("SMTP_FROM", "no-reply@ecommerce.com"),
		SMSProvider:      getEnv("SMS_PROVIDER", "logger"),
		TwilioAccountSID: os.Getenv("TWILIO_ACCOUNT_SID"),
		TwilioAuthToken:  os.Getenv("TWILIO_AUTH_TOKEN"),
		TwilioFromPhone:  os.Getenv("TWILIO_FROM_PHONE"),
	}
}

func getEnv(key, defaultVal string) string {
	if val := os.Getenv(key); val != "" {
		return val
	}
	return defaultVal
}