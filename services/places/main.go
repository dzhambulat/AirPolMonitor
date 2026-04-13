package main

import (
	"context"
	"log"
	"os"
	"os/signal"
	"strings"
	"syscall"

	"AirPolMonitor/services/places/app"
	"AirPolMonitor/services/places/infrastructure"
)

func main() {
	brokers := os.Getenv("KAFKA_BOOTSTRAP_SERVERS")
	if brokers == "" {
		brokers = "localhost:9092"
	}
	topic := os.Getenv("KAFKA_TOPIC")
	if topic == "" {
		topic = "air-pollution-data"
	}
	group := os.Getenv("KAFKA_GROUP_ID")
	if group == "" {
		group = "places-service"
	}

	defaultPlace := os.Getenv("PLACES_DEFAULT_PLACE_ID")
	if defaultPlace == "" {
		defaultPlace = "default-place"
	}

	repo := infrastructure.NewMemoryPlaceSensorsRepository()

	consumer, err := infrastructure.NewKafkaSubscriber(brokers, group, topic)
	if err != nil {
		log.Fatalf("kafka: %v", err)
	}
	defer func() {
		if err := consumer.Close(); err != nil {
			log.Printf("kafka close: %v", err)
		}
	}()

	ctx, cancel := signal.NotifyContext(context.Background(), syscall.SIGINT, syscall.SIGTERM)
	defer cancel()

	raw := make(chan []byte, 256)
	go func() {
		if err := consumer.Run(ctx, raw); err != nil && err != context.Canceled {
			log.Printf("kafka run: %v", err)
		}
	}()

	log.Printf("places: consuming topic %q from %s (default place %q)", topic, summarizeBrokers(brokers), defaultPlace)
	app.RunSensorStream(ctx, raw, repo)
	log.Println("places: shutdown complete")
}

func summarizeBrokers(b string) string {
	parts := strings.Split(b, ",")
	if len(parts) > 2 {
		return strings.Join(parts[:2], ",") + ",..."
	}
	return b
}
