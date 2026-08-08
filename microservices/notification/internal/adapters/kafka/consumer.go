package kafka

import (
	"context"
	"encoding/json"
	"fmt"
	"log"

	"github.com/Girmex/go-ecommerce/microservices/notification/internal/application"
	"github.com/Girmex/go-ecommerce/microservices/pkg/events"
	pkgKafka "github.com/Girmex/go-ecommerce/microservices/pkg/kafka"
)

type EventListener struct {
	brokers []string
	groupID string
	service *application.NotificationService
}

func NewEventListener(
	brokers []string,
	groupID string,
	service *application.NotificationService,
) *EventListener {
	return &EventListener{
		brokers: brokers,
		groupID: groupID,
		service: service,
	}
}

func (l *EventListener) Start(ctx context.Context) error {
	errChan := make(chan error, 4)

	// Payment completed consumer
	go func() {
		consumer := pkgKafka.NewConsumer(l.brokers, pkgKafka.TopicPaymentCompleted, l.groupID)
		defer consumer.Close()

		log.Printf("Started Kafka listener for topic: %s\n", pkgKafka.TopicPaymentCompleted)

		err := consumer.Start(ctx, func(ctx context.Context, key []byte, value []byte) error {
			var event events.PaymentCompleted
			if err := json.Unmarshal(value, &event); err != nil {
				return fmt.Errorf("failed to unmarshal PaymentCompleted event: %w", err)
			}
			return l.service.HandlePaymentCompleted(ctx, event)
		})
		if err != nil {
			errChan <- fmt.Errorf("payment completed consumer error: %w", err)
		}
	}()

	// User registered consumer
	go func() {
		consumer := pkgKafka.NewConsumer(l.brokers, pkgKafka.TopicUserRegistered, l.groupID)
		defer consumer.Close()

		log.Printf("Started Kafka listener for topic: %s\n", pkgKafka.TopicUserRegistered)

		err := consumer.Start(ctx, func(ctx context.Context, key []byte, value []byte) error {
			var event events.UserRegistered
			if err := json.Unmarshal(value, &event); err != nil {
				return fmt.Errorf("failed to unmarshal UserRegistered event: %w", err)
			}
			return l.service.HandleUserRegistered(ctx, event)
		})
		if err != nil {
			errChan <- fmt.Errorf("user registered consumer error: %w", err)
		}
	}()

	// User verification consumer
	go func() {
		consumer := pkgKafka.NewConsumer(l.brokers, pkgKafka.TopicUserVerification, l.groupID)
		defer consumer.Close()

		log.Printf("Started Kafka listener for topic: %s\n", pkgKafka.TopicUserVerification)

		err := consumer.Start(ctx, func(ctx context.Context, key []byte, value []byte) error {
			var event events.UserVerification
			if err := json.Unmarshal(value, &event); err != nil {
				return fmt.Errorf("failed to unmarshal UserVerification event: %w", err)
			}
			return l.service.HandleUserVerification(ctx, event)
		})
		if err != nil {
			errChan <- fmt.Errorf("user verification consumer error: %w", err)
		}
	}()

	// Order created consumer
	go func() {
		consumer := pkgKafka.NewConsumer(l.brokers, pkgKafka.TopicOrderCreated, l.groupID)
		defer consumer.Close()

		log.Printf("Started Kafka listener for topic: %s\n", pkgKafka.TopicOrderCreated)

		err := consumer.Start(ctx, func(ctx context.Context, key []byte, value []byte) error {
			var event events.OrderCreated
			if err := json.Unmarshal(value, &event); err != nil {
				return fmt.Errorf("failed to unmarshal OrderCreated event: %w", err)
			}
			return l.service.HandleOrderCreated(ctx, event)
		})
		if err != nil {
			errChan <- fmt.Errorf("order created consumer error: %w", err)
		}
	}()

	select {
	case <-ctx.Done():
		return ctx.Err()
	case err := <-errChan:
		return err
	}
}
