package kafka

import (
	"context"

	kafkapkg "github.com/Girmex/go-ecommerce/microservices/pkg/kafka"
)

type EventPublisher struct {
	producer *kafkapkg.Producer
}

func NewEventPublisher(producer *kafkapkg.Producer) *EventPublisher {
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
