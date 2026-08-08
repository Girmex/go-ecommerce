package application

import (
	"context"
	"fmt"
	"log"

	"github.com/Girmex/go-ecommerce/microservices/notification/internal/ports"
	"github.com/Girmex/go-ecommerce/microservices/pkg/events"
)

type NotificationService struct {
	emailSender ports.EmailSender
	smsSender   ports.SMSSender
}

func NewNotificationService(
	emailSender ports.EmailSender,
	smsSender ports.SMSSender,
) *NotificationService {
	return &NotificationService{
		emailSender: emailSender,
		smsSender:   smsSender,
	}
}

func (s *NotificationService) HandlePaymentCompleted(
	ctx context.Context,
	event events.PaymentCompleted,
) error {
	var errs []error

	if event.Email != "" && s.emailSender != nil {
		subject := "Payment Completed"
		body := fmt.Sprintf(
			"Your payment #%d for order #%d has been completed successfully. Amount: $%.2f",
			event.PaymentID,
			event.OrderID,
			event.Amount,
		)

		if err := s.emailSender.Send(ctx, event.Email, subject, body); err != nil {
			log.Printf("failed to send payment completed email to %s: %v\n", event.Email, err)
			errs = append(errs, err)
		}
	} else if event.Email == "" {
		log.Printf("[WARN] PaymentCompleted event #%d for user #%d has no recipient email specified\n", event.PaymentID, event.UserID)
	}

	if event.Phone != "" && s.smsSender != nil {
		message := fmt.Sprintf(
			"Payment #%d for order #%d of $%.2f is complete.",
			event.PaymentID,
			event.OrderID,
			event.Amount,
		)

		if err := s.smsSender.Send(ctx, event.Phone, message); err != nil {
			log.Printf("failed to send payment completed SMS to %s: %v\n", event.Phone, err)
			errs = append(errs, err)
		}
	}

	if len(errs) > 0 {
		return fmt.Errorf("errors occurred during PaymentCompleted notification handling: %v", errs)
	}

	return nil
}

func (s *NotificationService) HandleUserRegistered(
	ctx context.Context,
	event events.UserRegistered,
) error {
	if event.Email == "" || s.emailSender == nil {
		return nil
	}

	subject := "Welcome to Go-Ecommerce!"
	body := fmt.Sprintf("Hello %s,\n\nThank you for registering with us! We are thrilled to have you.", event.Name)

	return s.emailSender.Send(ctx, event.Email, subject, body)
}

func (s *NotificationService) HandleUserVerification(
	ctx context.Context,
	event events.UserVerification,
) error {
	if event.Email == "" || s.emailSender == nil {
		return nil
	}

	subject := "Verify Your Email"
	body := fmt.Sprintf("Please use the following verification code to confirm your email: %s", event.Token)

	return s.emailSender.Send(ctx, event.Email, subject, body)
}