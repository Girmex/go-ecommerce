package sms

import (
	"context"
	"log"

	"github.com/Girmex/go-ecommerce/microservices/notification/internal/ports"
)

type LoggerSMSSender struct{}

func NewLoggerSMSSender() ports.SMSSender {
	return &LoggerSMSSender{}
}

func (s *LoggerSMSSender) Send(ctx context.Context, phone string, message string) error {
	log.Printf("[SMS] To: %s | Message: %s\n", phone, message)
	return nil
}
