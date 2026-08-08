package kafka

import (
	"context"
	"fmt"

	"github.com/segmentio/kafka-go"
)

type MessageHandler func(ctx context.Context, key []byte, value []byte) error

type Consumer struct {
	reader *kafka.Reader
}

func NewConsumer(brokers []string, topic string, groupID string) *Consumer {
	return &Consumer{
		reader: kafka.NewReader(kafka.ReaderConfig{
			Brokers:  brokers,
			Topic:    topic,
			GroupID:  groupID,
			MinBytes: 10,
			MaxBytes: 10 * 1024 * 1024, // 10MB
		}),
	}
}

func (c *Consumer) Start(ctx context.Context, handler MessageHandler) error {
	for {
		select {
		case <-ctx.Done():
			return ctx.Err()
		default:
			msg, err := c.reader.ReadMessage(ctx)
			if err != nil {
				if ctx.Err() != nil {
					return ctx.Err()
				}
				return fmt.Errorf("read kafka message: %w", err)
			}

			if err := handler(ctx, msg.Key, msg.Value); err != nil {
				fmt.Printf("Error processing message from topic %s: %v\n", msg.Topic, err)
			}
		}
	}
}

func (c *Consumer) Close() error {
	return c.reader.Close()
}
