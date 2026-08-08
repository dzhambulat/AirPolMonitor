package mqtt

import (
	"AirPolMonitor/services/sensors/types"
	"context"
	"encoding/json"
	"log"
	"strings"
)

type SensorSource struct {
	subscriber *Subscriber
}

type sensorDataPayload struct {
	SensorID  string  `json:"sensor_id"`
	Value     float32 `json:"value"`
	Timestamp int64   `json:"timestamp"`
	Pollutant string  `json:"pollutant"`
}

func NewSensorSource(cfg SubscriberConfig) *SensorSource {
	return &SensorSource{subscriber: NewSubscriber(cfg)}
}

func (s *SensorSource) Subscribe(ctx context.Context) <-chan types.SensorData {
	out := make(chan types.SensorData)
	go func() {
		defer close(out)

		msgCh, err := s.subscriber.Subscribe(ctx)
		if err != nil {
			log.Printf("mqtt sensor source: subscribe failed: %v", err)
			return
		}

		for {
			select {
			case <-ctx.Done():
				return
			case msg, ok := <-msgCh:
				if !ok {
					return
				}
				var p sensorDataPayload
				if err := json.Unmarshal(msg.Payload, &p); err != nil {
					log.Printf("mqtt sensor source: bad payload: %v", err)
					continue
				}
				if p.SensorID == "" {
					log.Printf("mqtt sensor source: skip empty sensor_id")
					continue
				}
				out <- types.SensorData{
					SensorID:   p.SensorID,
					SensorType: sensorTypeFromPollutant(p.Pollutant),
					Value:      p.Value,
					Timestamp:  p.Timestamp,
				}
			}
		}
	}()
	return out
}

func sensorTypeFromPollutant(pollutant string) int {
	switch strings.ToLower(strings.TrimSpace(pollutant)) {
	case "pm25", "pm10", "aqi":
		return types.SensorTypeAirQuality
	case "temperature", "temp":
		return types.SensorTypeTemperature
	case "humidity":
		return types.SensorTypeHumidity
	case "co2":
		return types.SensorTypeCO2
	case "voc":
		return types.SensorTypeVOC
	default:
		return types.SensorTypeAirQuality
	}
}
