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
	smsSender   ports.SMSSender
}

func NewNotificationService(
	userClient ports.UserClient,
	emailSender ports.EmailSender,
	smsSender ports.SMSSender,
) *NotificationService {
	return &NotificationService{
		userClient:  userClient,
		emailSender: emailSender,
		smsSender:   smsSender,
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

func (s *NotificationService) HandlePhoneVerification(
	ctx context.Context,
	event events.UserPhoneVerification,
) error {

	message := fmt.Sprintf(
		"Your verification code is %s. It expires in 5 minutes.",
		event.Code,
	)

	if err := s.smsSender.Send(
		ctx,
		event.Phone,
		message,
	); err != nil {
		return fmt.Errorf(
			"send phone verification SMS: %w",
			err,
		)
	}

	return nil
}
