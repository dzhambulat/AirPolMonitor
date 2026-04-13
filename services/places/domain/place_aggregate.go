package domain

import (
	"context"
)

type SensorType int

const (
	SensorTypeAirQuality SensorType = iota
	SensorTypeTemperature
	SensorTypeHumidity
	SensorTypeCO2
	SensorTypeVOC
)

type SensorData struct {
	Value     float32
	Timestamp int64
	SensorType SensorType
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
	requestedData []SensorData
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


