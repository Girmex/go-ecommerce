package application_test

import (
	"context"
	"testing"

	"github.com/Girmex/go-ecommerce/microservices/notification/internal/adapters/email"
	"github.com/Girmex/go-ecommerce/microservices/notification/internal/adapters/sms"
	"github.com/Girmex/go-ecommerce/microservices/notification/internal/application"
	"github.com/Girmex/go-ecommerce/microservices/pkg/events"
)

func TestHandlePaymentCompleted(t *testing.T) {
	emailAdapter := email.NewLoggerEmailSender()
	smsAdapter := sms.NewLoggerSMSSender()

	svc := application.NewNotificationService(emailAdapter, smsAdapter)

	event := events.PaymentCompleted{
		PaymentID: 101,
		OrderID:   55,
		UserID:    12,
		Email:     "user@example.com",
		Phone:     "+1234567890",
		Amount:    99.99,
	}

	err := svc.HandlePaymentCompleted(context.Background(), event)
	if err != nil {
		t.Fatalf("expected no error, got: %v", err)
	}
}

func TestHandleUserRegistered(t *testing.T) {
	emailAdapter := email.NewLoggerEmailSender()
	smsAdapter := sms.NewLoggerSMSSender()

	svc := application.NewNotificationService(emailAdapter, smsAdapter)

	event := events.UserRegistered{
		UserID: 12,
		Name:   "John Doe",
		Email:  "john@example.com",
	}

	err := svc.HandleUserRegistered(context.Background(), event)
	if err != nil {
		t.Fatalf("expected no error, got: %v", err)
	}
}

func TestHandleUserVerification(t *testing.T) {
	emailAdapter := email.NewLoggerEmailSender()
	smsAdapter := sms.NewLoggerSMSSender()

	svc := application.NewNotificationService(emailAdapter, smsAdapter)

	event := events.UserVerification{
		UserID: 12,
		Email:  "john@example.com",
		Token:  "ABC-123-XYZ",
	}

	err := svc.HandleUserVerification(context.Background(), event)
	if err != nil {
		t.Fatalf("expected no error, got: %v", err)
	}
}
