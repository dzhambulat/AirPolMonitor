package readmodel

import (
	"context"
)

// MockAnalyticsRepository is a fake repository for local development.
// It will be replaced later by a ClickHouse-backed implementation.
type MockAnalyticsRepository struct {
	data []SensorData
}

func NewMockAnalyticsRepository(seed []SensorData) *MockAnalyticsRepository {
	cp := append([]SensorData(nil), seed...)
	return &MockAnalyticsRepository{data: cp}
}

func (m *MockAnalyticsRepository) ListByPlaceAndTypes(ctx context.Context, placeID string, sensorTypes []SensorType) ([]SensorData, error) {
	_ = ctx

	typeFilter := make(map[SensorType]struct{}, len(sensorTypes))
	for _, t := range sensorTypes {
		typeFilter[t] = struct{}{}
	}
	filterByType := len(typeFilter) > 0

	out := make([]SensorData, 0, len(m.data))
	for _, row := range m.data {
		if placeID != "" && row.PlaceID != placeID {
			continue
		}
		if filterByType {
			if _, ok := typeFilter[row.SensorType]; !ok {
				continue
			}
		}
		out = append(out, row)
	}
	return out, nil
}

func SeedDemoData() []SensorData {
	return []SensorData{
		{PlaceID: "place-a", SensorType: SensorTypeAirQuality, Value: 32.1, Timestamp: 1710000001},
		{PlaceID: "place-a", SensorType: SensorTypeAirQuality, Value: 28.6, Timestamp: 1710000301},
		{PlaceID: "place-a", SensorType: SensorTypeCO2, Value: 640, Timestamp: 1710000601},
		{PlaceID: "place-a", SensorType: SensorTypeCO2, Value: 720, Timestamp: 1710000901},
		{PlaceID: "place-b", SensorType: SensorTypeHumidity, Value: 42.5, Timestamp: 1710000001},
		{PlaceID: "place-b", SensorType: SensorTypeHumidity, Value: 48.2, Timestamp: 1710000601},
	}
}
