package email

import (
	"context"
	"log"

	"github.com/Girmex/go-ecommerce/microservices/notification/internal/ports"
)

type LoggerEmailSender struct{}

func NewLoggerEmailSender() ports.EmailSender {
	return &LoggerEmailSender{}
}

func (s *LoggerEmailSender) Send(
	ctx context.Context,
	to string,
	subject string,
	body string,
) error {
	log.Printf("[EMAIL] To: %s | Subject: %s | Body: %s\n", to, subject, body)
	return nil
}
