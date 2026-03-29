package main

import (
	"AirPolMonitor/config"
	"AirPolMonitor/core"
	"AirPolMonitor/infrastructure/kafka"
	sensorsource "AirPolMonitor/infrastructure/sensor_source"
	"context"
	"log"
	"net/http"
	"os"

	"github.com/prometheus/client_golang/prometheus/promhttp"
)

func main() {
	cfg, err := config.LoadConfig("config.json")
	if err != nil {
		log.Fatal("cannot load config:", err)
	}

	kafkaProducer, err := kafka.NewProducer(cfg.KafkaBrokerURL, cfg.KafkaTopic)
	if err != nil {
		log.Fatal("cannot create kafka producer:", err)
	}

	// Metrics server for Prometheus (scraped at :9091 when using docker-compose)
	metricsAddr := ":9091"
	if p := os.Getenv("METRICS_PORT"); p != "" {
		metricsAddr = ":" + p
	}
	mux := http.NewServeMux()
	mux.Handle("/metrics", promhttp.Handler())
	go func() {
		log.Printf("metrics listening on %s", metricsAddr)
		if err := http.ListenAndServe(metricsAddr, mux); err != nil && err != http.ErrServerClosed {
			log.Printf("metrics server: %v", err)
		}
	}()

	ctx, cancel := context.WithCancel(context.Background())
	defer func() {
		log.Println("Shutting down...")
		cancel()
		kafkaProducer.Close()
	}()

	mockSensorSource := sensorsource.NewMockSensorSource()
	err = core.Start(ctx, 1, mockSensorSource, kafkaProducer)
	if err != nil {
		log.Fatal("cannot start core:", err)
	}
}