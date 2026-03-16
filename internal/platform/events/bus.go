package events

import (
	"context"
	"encoding/json"
	"errors"

	"github.com/IBM/sarama"
	"go.uber.org/zap"
)

type HandlerFunc func(ctx context.Context, event Event) error

type KafkaBus interface {
	Publisher
	Subscribe(ctx context.Context, topics []string, group sarama.ConsumerGroup, handler HandlerFunc) error
	Close() error
}

type kafkaBus struct {
	producer sarama.SyncProducer
	log      *zap.Logger
}

func NewKafkaBus(producer sarama.SyncProducer, log *zap.Logger) KafkaBus {
	return &kafkaBus{producer: producer, log: log}
}

func (k *kafkaBus) Publish(ctx context.Context, topic string, event Event) error {
	b, err := json.Marshal(event)
	if err != nil {
		return err
	}

	msg := &sarama.ProducerMessage{
		Topic: topic,
		Key:   sarama.StringEncoder(event.AggregateID),
		Value: sarama.ByteEncoder(b),
	}

	_, _, err = k.producer.SendMessage(msg)
	return err
}

func (k *kafkaBus) Subscribe(ctx context.Context, topics []string, group sarama.ConsumerGroup, handler HandlerFunc) error {
	if len(topics) == 0 {
		return errors.New("no topics provided")
	}

	consumer := &consumerGroupHandler{
		log:     k.log,
		handler: handler,
	}

	for {
		if ctx.Err() != nil {
			return nil
		}

		if err := group.Consume(ctx, topics, consumer); err != nil {
			k.log.Error("kafka consume error", zap.Error(err), zap.Strings("topics", topics))
			return err
		}
	}
}

func (k *kafkaBus) Close() error {
	if k.producer != nil {
		return k.producer.Close()
	}
	return nil
}

type consumerGroupHandler struct {
	log     *zap.Logger
	handler HandlerFunc
}

func (h *consumerGroupHandler) Setup(_ sarama.ConsumerGroupSession) error   { return nil }
func (h *consumerGroupHandler) Cleanup(_ sarama.ConsumerGroupSession) error { return nil }

func (h *consumerGroupHandler) ConsumeClaim(session sarama.ConsumerGroupSession, claim sarama.ConsumerGroupClaim) error {
	for msg := range claim.Messages() {
		var evt Event
		if err := json.Unmarshal(msg.Value, &evt); err != nil {
			h.log.Error("failed to unmarshal kafka event", zap.Error(err))
			continue
		}

		if err := h.handler(session.Context(), evt); err != nil {
			h.log.Error("event handler failed", zap.Error(err), zap.String("topic", msg.Topic))
			continue
		}

		session.MarkMessage(msg, "")
	}
	return nil
}
