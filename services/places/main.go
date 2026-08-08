package main

import (
	"context"
	"log"
	"net/http"
	"os"
	"os/signal"
	"strings"
	"syscall"
	"time"

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
	httpAddr := os.Getenv("PLACES_HTTP_ADDR")
	if httpAddr == "" {
		httpAddr = ":8082"
	}

	sensorsRepo := infrastructure.NewMemoryPlaceSensorsRepository()
	placesRepo := infrastructure.NewMemoryPlaceRepository()

	mux := http.NewServeMux()
	app.NewPlaceEndpoints(placesRepo).Register(mux)
	httpServer := &http.Server{
		Addr:    httpAddr,
		Handler: mux,
	}
	go func() {
		log.Printf("places: http listening on %s", httpAddr)
		if err := httpServer.ListenAndServe(); err != nil && err != http.ErrServerClosed {
			log.Printf("places: http server: %v", err)
		}
	}()

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
	defer func() {
		shutdownCtx, shutdownCancel := context.WithTimeout(context.Background(), 5*time.Second)
		defer shutdownCancel()
		if err := httpServer.Shutdown(shutdownCtx); err != nil {
			log.Printf("places: http shutdown: %v", err)
		}
	}()

	raw := make(chan []byte, 256)
	go func() {
		log.Println("kafka run: starting")
		err:=consumer.Run(ctx, raw)
		log.Println("kafka run: done", err)
		if err != nil{
			log.Printf("kafka run: %v", err)
		}
		log.Println("kafka run: done", err)
	}()

	log.Printf("places: consuming topic %q from %s (default place %q)", topic, summarizeBrokers(brokers), defaultPlace)
	app.RunSensorStream(ctx, raw, sensorsRepo)
	log.Println("places: shutdown complete")
}

func summarizeBrokers(b string) string {
	parts := strings.Split(b, ",")
	if len(parts) > 2 {
		return strings.Join(parts[:2], ",") + ",..."
	}
	return b
}
