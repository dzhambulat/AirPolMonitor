package kafka

import (
	"AirPolMonitor/core/types"
	"context"
	"encoding/json"
	"log"

	"github.com/confluentinc/confluent-kafka-go/v2/kafka"
)

type KafkaProducer struct {
	producer *kafka.Producer
	topic string
}

func (p *KafkaProducer) Produce(ctx context.Context, data <-chan types.AirData) error {
	for {
		select {
		case <-ctx.Done():
			return ctx.Err()
		case d := <-data:
			payload, err := json.Marshal(d)
			if err != nil {
				log.Printf("kafka: marshal AirData: %v", err)
				continue
			}
			if err := p.ProduceMessage(payload); err != nil {
				log.Printf("kafka: produce: %v", err)
			}
		}
	}
}

func NewProducer(kafkaBrokerURL string, topic string) (*KafkaProducer, error) {
	log.Println("Connecting to KAFKA...")
	p, err := kafka.NewProducer(&kafka.ConfigMap{"bootstrap.servers": kafkaBrokerURL})
	if err != nil {
		return nil, err
	}
	return &KafkaProducer{producer: p, topic: topic}, nil
}

// ProduceMessage sends a single message to the given topic (for MQTT bridge, etc.).
func (p *KafkaProducer) ProduceMessage(message []byte) error {
	deliveryChan := make(chan kafka.Event)
	defer close(deliveryChan)

	err := p.producer.Produce(&kafka.Message{
		TopicPartition: kafka.TopicPartition{Topic: &p.topic, Partition: kafka.PartitionAny},
		Value:          message,
	}, deliveryChan)
	if err != nil {
		return err
	}
	e := <-deliveryChan
	m := e.(*kafka.Message)
	if m.TopicPartition.Error != nil {
		return m.TopicPartition.Error
	}

	return nil
}

// Close shuts down the producer.
func (p *KafkaProducer) Close() {
	p.producer.Close()
}
