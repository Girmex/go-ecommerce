package kafka

import (
	"context"
	"encoding/json"
	"fmt"

	"github.com/segmentio/kafka-go"
)

type Producer struct {
	writer *kafka.Writer
}

func NewProducer(brokers []string) *Producer {
	return &Producer{
		writer: &kafka.Writer{
			Addr:     kafka.TCP(brokers...),
			Balancer: &kafka.LeastBytes{},
		},
	}
}

func (p *Producer) Publish(
	ctx context.Context,
	topic string,
	key string,
	event any,
) error {
	payload, err := json.Marshal(event)
	if err != nil {
		return fmt.Errorf("marshal kafka event: %w", err)
	}

	err = p.writer.WriteMessages(
		ctx,
		kafka.Message{
			Topic: topic,
			Key:   []byte(key),
			Value: payload,
		},
	)
	if err != nil {
		return fmt.Errorf("publish kafka event: %w", err)
	}

	return nil
}

func (p *Producer) Close() error {
	return p.writer.Close()
}
