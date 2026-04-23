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
	coordinates []float64
	requestedData []SensorData
}

// CreatePlace is an application command for creating a place aggregate.
type EditPlaceCommand struct {
	ID          string
	Name        string
	Coordinates []float64
}

func GetPlaceAggregate(id, name string, repo PlaceRepositoryPort) (*PlaceAggregate, error) {
	if id == "" {
		agg:= PlaceAggregate{
			id: id,
			name: name,
		}

		repo.Save(context.Background(), &agg)
		return &agg, nil
	}

	ctx:=context.Background()
	agg, err := repo.Get(ctx, id)
	if err != nil {
		return nil, err
	}
	return agg, nil
}

// HandleCreatePlace applies CreatePlace command data to this aggregate.
func (p *PlaceAggregate) HandleEditPlaceCommand(cmd EditPlaceCommand) error {
	if cmd.ID == "" {
		return fmt.Errorf("place id is required")
	}
	if len(cmd.Coordinates) == 0 {
		return fmt.Errorf("coordinates are required")
	}

	p.id = cmd.ID
	p.name = cmd.Name
	p.coordinates = append([]float64(nil), cmd.Coordinates...)
	return nil
}


