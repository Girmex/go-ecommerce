package application

import (
	"context"
	"fmt"

	"github.com/Girmex/go-ecommerce/microservices/notification/internal/ports"
	"github.com/Girmex/go-ecommerce/microservices/pkg/events"
)

type NotificationService struct {
	emailSender ports.EmailSender
}

func NewNotificationService(
	emailSender ports.EmailSender,
) *NotificationService {
	return &NotificationService{
		emailSender: emailSender,
	}
}

func (s *NotificationService) HandlePaymentCompleted(
	ctx context.Context,
	event events.PaymentCompleted,
) error {
	subject := "Payment Completed"

	body := fmt.Sprintf(
		"Your payment #%d for order #%d has been completed successfully. Amount: %.2f",
		event.PaymentID,
		event.OrderID,
		event.Amount,
	)

	return s.emailSender.Send(
		ctx,
		"",
		subject,
		body,
	)
}
