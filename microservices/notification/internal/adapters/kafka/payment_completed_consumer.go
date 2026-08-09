package kafka

import (
	"context"
	"encoding/json"
	"fmt"

	"github.com/Girmex/go-ecommerce/microservices/notification/internal/application"
	"github.com/Girmex/go-ecommerce/microservices/pkg/events"
	"github.com/Girmex/go-ecommerce/microservices/pkg/kafka"
)

type PaymentCompletedConsumer struct {
	consumer *kafka.Consumer
	service  *application.NotificationService
}

func NewPaymentCompletedConsumer(
	consumer *kafka.Consumer,
	service *application.NotificationService,
) *PaymentCompletedConsumer {
	return &PaymentCompletedConsumer{
		consumer: consumer,
		service:  service,
	}
}

func (c *PaymentCompletedConsumer) Start(
	ctx context.Context,
) error {
	return c.consumer.Start(
		ctx,
		c.handleMessage,
	)
}

func (c *PaymentCompletedConsumer) handleMessage(
	ctx context.Context,
	key []byte,
	value []byte,
) error {

	var event events.PaymentCompleted

	if err := json.Unmarshal(value, &event); err != nil {
		return fmt.Errorf(
			"unmarshal payment completed event: %w",
			err,
		)
	}

	if err := c.service.HandlePaymentCompleted(
		ctx,
		event,
	); err != nil {
		return fmt.Errorf(
			"handle payment completed event: %w",
			err,
		)
	}

	return nil
}

func (c *PaymentCompletedConsumer) Close() error {
	return c.consumer.Close()
}
