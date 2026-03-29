package observability

import (
	"time"

	"github.com/prometheus/client_golang/prometheus"
	"github.com/prometheus/client_golang/prometheus/promauto"
)

// EventToKafkaDuration records time from event received (by core) to Kafka produce ack.
var EventToKafkaDuration = promauto.NewHistogram(prometheus.HistogramOpts{
	Name:    "event_to_kafka_duration_seconds",
	Help:    "Time from event received to Kafka produce acknowledged (MQTT/core → Kafka).",
	Buckets: prometheus.ExponentialBuckets(0.001, 2, 12), // 1ms to ~2s
})

// ObserveEventToKafka records one observation for the event→Kafka latency.
func ObserveEventToKafka(d time.Duration) {
	EventToKafkaDuration.Observe(d.Seconds())
}
