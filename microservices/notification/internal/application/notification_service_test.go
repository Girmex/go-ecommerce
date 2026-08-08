package application_test

import (
	"context"
	"fmt"
	"testing"

	authproto "github.com/Girmex/go-ecommerce/microservices/auth/proto"
	"github.com/Girmex/go-ecommerce/microservices/notification/internal/adapters/email"
	"github.com/Girmex/go-ecommerce/microservices/notification/internal/adapters/sms"
	"github.com/Girmex/go-ecommerce/microservices/notification/internal/application"
	"github.com/Girmex/go-ecommerce/microservices/pkg/events"
)

type mockAuthClient struct {
	users map[uint]*authproto.User
}

func (m *mockAuthClient) GetUser(ctx context.Context, userID uint) (*authproto.User, error) {
	if u, ok := m.users[userID]; ok {
		return u, nil
	}
	return nil, fmt.Errorf("user not found")
}

func TestHandlePaymentCompleted(t *testing.T) {
	emailAdapter := email.NewLoggerEmailSender()
	smsAdapter := sms.NewLoggerSMSSender()
	authClient := &mockAuthClient{
		users: map[uint]*authproto.User{
			12: {Id: 12, Name: "John Doe", Email: "john@example.com"},
		},
	}

	svc := application.NewNotificationService(emailAdapter, smsAdapter, authClient)

	// Test with explicit email
	eventWithEmail := events.PaymentCompleted{
		PaymentID: 101,
		OrderID:   55,
		UserID:    12,
		Email:     "explicit@example.com",
		Phone:     "+1234567890",
		Amount:    99.99,
	}

	err := svc.HandlePaymentCompleted(context.Background(), eventWithEmail)
	if err != nil {
		t.Fatalf("expected no error, got: %v", err)
	}

	// Test with empty email (should resolve via authClient)
	eventWithoutEmail := events.PaymentCompleted{
		PaymentID: 102,
		OrderID:   56,
		UserID:    12,
		Amount:    150.00,
	}

	err = svc.HandlePaymentCompleted(context.Background(), eventWithoutEmail)
	if err != nil {
		t.Fatalf("expected no error resolving user email, got: %v", err)
	}
}

func TestHandleOrderCreated(t *testing.T) {
	emailAdapter := email.NewLoggerEmailSender()
	smsAdapter := sms.NewLoggerSMSSender()
	authClient := &mockAuthClient{
		users: map[uint]*authproto.User{
			12: {Id: 12, Name: "Jane Doe", Email: "jane@example.com"},
		},
	}

	svc := application.NewNotificationService(emailAdapter, smsAdapter, authClient)

	event := events.OrderCreated{
		OrderID: 200,
		UserID:  12,
		Amount:  299.50,
	}

	err := svc.HandleOrderCreated(context.Background(), event)
	if err != nil {
		t.Fatalf("expected no error, got: %v", err)
	}
}

func TestHandleUserRegistered(t *testing.T) {
	emailAdapter := email.NewLoggerEmailSender()
	smsAdapter := sms.NewLoggerSMSSender()

	svc := application.NewNotificationService(emailAdapter, smsAdapter, nil)

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

	svc := application.NewNotificationService(emailAdapter, smsAdapter, nil)

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
