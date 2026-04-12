package domain

import (
	"context"
	"fmt"
)

type SensorType int

const (
	SensorTypeAirQuality SensorType = iota
	SensorTypeTemperature
	SensorTypeHumidity
	SensorTypeCO2
	SensorTypeVOC
)

type AirSensorDataPayload struct {
	PlaceID   string  `json:"PlaceID,omitempty"`
	Value     float32 `json:"Value"`
	Timestamp int64   `json:"Timestamp"`
	SensorType SensorType `json:"SensorType"`
}

// PlaceRepositoryPort persists and loads place aggregates.
type PlaceRepositoryPort interface {
	Get(ctx context.Context, placeID string) (*PlaceAggregate, error)
	Save(ctx context.Context, agg *PlaceAggregate) error
}

// PlaceAggregate is a place-centric projection updated from sensor stream events.
type PlaceAggregate struct {
	id          string
	name        string
}

func NewPlaceAggregate(id, name string, repo PlaceRepositoryPort) (*PlaceAggregate, error) {
	if id == "" {
		return &PlaceAggregate{
			id: id,
			name: name,
		}, nil
	}

	ctx:=context.Background()
	agg, err := repo.Get(ctx, id)
	if err != nil {
		return nil, err
	}
	return agg, nil
}

// ApplySensorData merges one air reading into this place aggregate.
func (p *PlaceAggregate) ApplySensorData(r AirSensorDataPayload) error {
	
	return nil
}


// AddSensorDataToPlace loads the aggregate, applies the reading, and saves it.
func AddSensorDataToPlace(ctx context.Context, repo PlaceRepositoryPort, placeID string, r AirSensorDataPayload) error {
	if repo == nil {
		return fmt.Errorf("place repository is nil")
	}
	agg, err := repo.Get(ctx, placeID)
	if err != nil {
		return err
	}
	if err := agg.ApplySensorData(r); err != nil {
		return err
	}
	return repo.Save(ctx, agg)
}
