package sms

import (
	"context"
	"log"
)

type LoggingSMSSender struct{}

func NewLoggingSMSSender() *LoggingSMSSender {
	return &LoggingSMSSender{}
}

func (s *LoggingSMSSender) Send(
	ctx context.Context,
	phone string,
	message string,
) error {
	log.Printf(
		"SMS → %s: %s",
		phone,
		message,
	)

	return nil
}
