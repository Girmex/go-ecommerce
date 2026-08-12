package kafka

import (
	"context"
	"encoding/json"
	"fmt"

	"github.com/Girmex/go-ecommerce/microservices/notification/internal/application"
	"github.com/Girmex/go-ecommerce/microservices/pkg/events"
	kafkapkg "github.com/Girmex/go-ecommerce/microservices/pkg/kafka"
)

type PhoneVerificationConsumer struct {
	consumer            *kafkapkg.Consumer
	notificationService *application.NotificationService
}

func NewPhoneVerificationConsumer(
	consumer *kafkapkg.Consumer,
	notificationService *application.NotificationService,
) *PhoneVerificationConsumer {
	return &PhoneVerificationConsumer{
		consumer:            consumer,
		notificationService: notificationService,
	}
}

func (c *PhoneVerificationConsumer) Start(ctx context.Context) error {
	return c.consumer.Start(ctx, func(
		ctx context.Context,
		key []byte,
		value []byte,
	) error {

		var event events.UserPhoneVerification

		if err := json.Unmarshal(value, &event); err != nil {
			return fmt.Errorf(
				"unmarshal phone verification event: %w",
				err,
			)
		}

		if err := c.notificationService.HandlePhoneVerification(
			ctx,
			event,
		); err != nil {
			return fmt.Errorf(
				"handle phone verification event: %w",
				err,
			)
		}

		return nil
	})
}
