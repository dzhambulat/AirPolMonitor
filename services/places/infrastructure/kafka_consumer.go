package infrastructure

import (
	"context"
	"fmt"

	"github.com/confluentinc/confluent-kafka-go/v2/kafka"
)

// KafkaSubscriber wraps a confluent-kafka consumer and forwards message payloads only.
type KafkaSubscriber struct {
	consumer *kafka.Consumer
	topic    string
}

// NewKafkaSubscriber creates a consumer subscribed to a single topic (consumer group).
func NewKafkaSubscriber(bootstrapServers, groupID, topic string) (*KafkaSubscriber, error) {
	c, err := kafka.NewConsumer(&kafka.ConfigMap{
		"bootstrap.servers":             bootstrapServers,
		"group.id":                      groupID,
		"auto.offset.reset":             "earliest",
		"enable.auto.commit":          true,
		"session.timeout.ms":            45000,
		"partition.assignment.strategy": "cooperative-sticky",
	})
	if err != nil {
		return nil, fmt.Errorf("kafka consumer: %w", err)
	}
	if err := c.Subscribe(topic, nil); err != nil {
		_ = c.Close()
		return nil, fmt.Errorf("kafka subscribe %q: %w", topic, err)
	}
	return &KafkaSubscriber{consumer: c, topic: topic}, nil
}

// Run polls the broker until ctx is cancelled. It sends each message value (copy) to out.
// The caller must drain out or Run will block when the buffer fills.
func (k *KafkaSubscriber) Run(ctx context.Context, out chan<- []byte) error {
	defer close(out)
	for {
		select {
		case <-ctx.Done():
			return ctx.Err()
		default:
		}
		ev := k.consumer.Poll(100)
		if ev == nil {
			continue
		}
		switch e := ev.(type) {
		case *kafka.Message:
			if e.Value == nil {
				continue
			}
			payload := append([]byte(nil), e.Value...)
			select {
			case out <- payload:
			case <-ctx.Done():
				return ctx.Err()
			}
		case kafka.Error:
			if e.IsFatal() {
				return fmt.Errorf("kafka fatal error: %w", e)
			}
			// transient errors: continue polling
		default:
			// ignore other event types
		}
	}
}

// Close releases the consumer.
func (k *KafkaSubscriber) Close() error {
	if k.consumer == nil {
		return nil
	}
	return k.consumer.Close()
}
