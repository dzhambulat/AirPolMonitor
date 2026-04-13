package app

import (
	"AirPolMonitor/services/places/domain"
	"AirPolMonitor/services/places/infrastructure"
	"context"
	"encoding/json"
	"log"
)

type AirSensorInput struct {
	Coordinates []float64
	SensorType domain.SensorType
	Value float32
	Timestamp int64
}

func getPlaceIdByCoordinates(ctx context.Context, coordinates []float64) (string, error) {
	return "place-id", nil
}

// RunSensorStream reads raw Kafka message bodies, unmarshals them as air readings,
// and applies them to the place aggregate via the repository.
func RunSensorStream(ctx context.Context, raw <-chan []byte, repo *infrastructure.MemoryPlaceSensorsRepository) {
	for {
		select {
		case <-ctx.Done():
			return
		case payload, ok := <-raw:
			if !ok {
				return
			}
			var msg AirSensorInput
			if err := json.Unmarshal(payload, &msg); err != nil {
				log.Printf("places: skip message: unmarshal: %v", err)
				continue
			}

			// ToDo - get place agg by coordinates
			placeID, err := getPlaceIdByCoordinates(ctx, msg.Coordinates)
			if err != nil {
				log.Printf("places: get place id by coordinates: %v", err)
				continue
			}

			repo.Save(ctx, placeID, infrastructure.SensorData{
				SensorType: msg.SensorType,
				Value: msg.Value,
				Timestamp: msg.Timestamp,
			})
		}
	}
}
