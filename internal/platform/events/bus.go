package events

import (
	"context"
	"encoding/json"

	"github.com/IBM/sarama"
	"go.uber.org/zap"
)

type KafkaBus interface {
	Publisher
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

func (k *kafkaBus) Close() error {
	return k.producer.Close()
}
