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
	authClient  ports.AuthClient
}

func NewNotificationService(
	emailSender ports.EmailSender,
	smsSender ports.SMSSender,
	authClient ports.AuthClient,
) *NotificationService {
	return &NotificationService{
		emailSender: emailSender,
		smsSender:   smsSender,
		authClient:  authClient,
	}
}

func (s *NotificationService) HandlePaymentCompleted(
	ctx context.Context,
	event events.PaymentCompleted,
) error {
	var errs []error
	userEmail := event.Email
	userPhone := event.Phone

	// Resolve missing user details from Auth service via gRPC
	if userEmail == "" && event.UserID != 0 && s.authClient != nil {
		user, err := s.authClient.GetUser(ctx, event.UserID)
		if err != nil {
			log.Printf("[WARN] Failed to resolve user profile for UserID %d: %v\n", event.UserID, err)
		} else if user != nil {
			userEmail = user.Email
		}
	}

	if userEmail != "" && s.emailSender != nil {
		subject := "Payment Completed"
		body := fmt.Sprintf(
			"Your payment #%d for order #%d has been completed successfully. Amount: $%.2f",
			event.PaymentID,
			event.OrderID,
			event.Amount,
		)

		if err := s.emailSender.Send(ctx, userEmail, subject, body); err != nil {
			log.Printf("failed to send payment completed email to %s: %v\n", userEmail, err)
			errs = append(errs, err)
		}
	} else if userEmail == "" {
		log.Printf("[WARN] PaymentCompleted event #%d for user #%d has no recipient email specified\n", event.PaymentID, event.UserID)
	}

	if userPhone != "" && s.smsSender != nil {
		message := fmt.Sprintf(
			"Payment #%d for order #%d of $%.2f is complete.",
			event.PaymentID,
			event.OrderID,
			event.Amount,
		)

		if err := s.smsSender.Send(ctx, userPhone, message); err != nil {
			log.Printf("failed to send payment completed SMS to %s: %v\n", userPhone, err)
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

func (s *NotificationService) HandleOrderCreated(
	ctx context.Context,
	event events.OrderCreated,
) error {
	userEmail := ""

	// Resolve user email via Auth gRPC client
	if event.UserID != 0 && s.authClient != nil {
		user, err := s.authClient.GetUser(ctx, event.UserID)
		if err != nil {
			log.Printf("[WARN] Failed to resolve user for order #%d: %v\n", event.OrderID, err)
		} else if user != nil {
			userEmail = user.Email
		}
	}

	if userEmail != "" && s.emailSender != nil {
		subject := "Order Confirmation"
		body := fmt.Sprintf("Your order #%d has been placed successfully. Total Amount: $%.2f", event.OrderID, event.Amount)

		return s.emailSender.Send(ctx, userEmail, subject, body)
	}

	return nil
}