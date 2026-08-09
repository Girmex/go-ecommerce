package application

import (
	"context"
	"fmt"

	"github.com/Girmex/go-ecommerce/microservices/notification/internal/ports"
	"github.com/Girmex/go-ecommerce/microservices/pkg/events"
)

type NotificationService struct {
	userClient  ports.UserClient
	emailSender ports.EmailSender
}

func NewNotificationService(
	userClient ports.UserClient,
	emailSender ports.EmailSender,
) *NotificationService {
	return &NotificationService{
		userClient:  userClient,
		emailSender: emailSender,
	}
}

func (s *NotificationService) HandlePaymentCompleted(
	ctx context.Context,
	event events.PaymentCompleted,
) error {

	user, err := s.userClient.GetUser(
		ctx,
		event.UserID,
	)
	if err != nil {
		return fmt.Errorf(
			"get user for payment notification: %w",
			err,
		)
	}

	subject := "Payment Completed"

	body := fmt.Sprintf(
		"Hello %s,\n\n"+
			"Your payment #%d for order #%d "+
			"has been completed successfully.\n\n"+
			"Amount: %.2f",
		user.Name,
		event.PaymentID,
		event.OrderID,
		event.Amount,
	)

	if err := s.emailSender.Send(
		ctx,
		user.Email,
		subject,
		body,
	); err != nil {
		return fmt.Errorf(
			"send payment notification: %w",
			err,
		)
	}

	return nil
}