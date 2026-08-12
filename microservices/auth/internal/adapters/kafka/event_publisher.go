package kafka

import (
	"context"

	"github.com/Girmex/go-ecommerce/microservices/auth/internal/ports"
	kafkapkg "github.com/Girmex/go-ecommerce/microservices/pkg/kafka"
)

var _ ports.EventPublisher = (*EventPublisher)(nil)

type EventPublisher struct {
	producer *kafkapkg.Producer
}

func NewEventPublisher(
	producer *kafkapkg.Producer,
) *EventPublisher {
	return &EventPublisher{
		producer: producer,
	}
}

func (p *EventPublisher) Publish(
	ctx context.Context,
	topic string,
	key string,
	event any,
) error {
	return p.producer.Publish(
		ctx,
		topic,
		key,
		event,
	)
}
