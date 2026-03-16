package events

import (
	"context"

	"github.com/IBM/sarama"
)

// NoopBus is a no-op implementation of KafkaBus for when Kafka is disabled.
type NoopBus struct{}

func (n *NoopBus) Publish(ctx context.Context, topic string, event Event) error {
	return nil
}

func (n *NoopBus) Subscribe(ctx context.Context, topics []string, group sarama.ConsumerGroup, handler HandlerFunc) error {
	return nil
}

func (n *NoopBus) Close() error {
	return nil
}
