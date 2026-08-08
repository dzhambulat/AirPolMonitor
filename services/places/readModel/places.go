package readmodel

import (
	"context"
	"fmt"
	"sort"
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
	PlaceID    string
	SensorType SensorType
	Value      float64
	Timestamp  int64
}

type AnalyticsRepository interface {
	ListByPlaceAndTypes(ctx context.Context, placeID string, sensorTypes []SensorType) ([]SensorData, error)
}

type SensorStats struct {
	Count int
	Avg   float64
	Min   float64
	Max   float64
}

func GetPlaceSensorData(ctx context.Context, repo AnalyticsRepository, placeID string, sensorTypes []SensorType) ([]SensorData, error) {
	if repo == nil {
		return nil, fmt.Errorf("analytics repository is nil")
	}
	return repo.ListByPlaceAndTypes(ctx, placeID, sensorTypes)
}

// GetLatestReadingsBySensorType returns one latest point per sensor type.
func GetLatestReadingsBySensorType(ctx context.Context, repo AnalyticsRepository, placeID string, sensorTypes []SensorType) (map[SensorType]SensorData, error) {
	rows, err := GetPlaceSensorData(ctx, repo, placeID, sensorTypes)
	if err != nil {
		return nil, err
	}

	latest := make(map[SensorType]SensorData)
	for _, row := range rows {
		existing, ok := latest[row.SensorType]
		if !ok || row.Timestamp > existing.Timestamp {
			latest[row.SensorType] = row
		}
	}
	return latest, nil
}

// GetSensorTypeStats calculates count/avg/min/max for one sensor type.
func GetSensorTypeStats(ctx context.Context, repo AnalyticsRepository, placeID string, sensorType SensorType) (SensorStats, error) {
	rows, err := GetPlaceSensorData(ctx, repo, placeID, []SensorType{sensorType})
	if err != nil {
		return SensorStats{}, err
	}
	if len(rows) == 0 {
		return SensorStats{}, nil
	}

	stats := SensorStats{
		Count: len(rows),
		Min:   rows[0].Value,
		Max:   rows[0].Value,
	}
	sum := 0.0
	for _, row := range rows {
		sum += row.Value
		if row.Value < stats.Min {
			stats.Min = row.Value
		}
		if row.Value > stats.Max {
			stats.Max = row.Value
		}
	}
	stats.Avg = sum / float64(stats.Count)
	return stats, nil
}

// GetTimeSeries sorts and returns readings for graph-ready output.
func GetTimeSeries(ctx context.Context, repo AnalyticsRepository, placeID string, sensorTypes []SensorType) ([]SensorData, error) {
	rows, err := GetPlaceSensorData(ctx, repo, placeID, sensorTypes)
	if err != nil {
		return nil, err
	}
	sort.Slice(rows, func(i, j int) bool {
		return rows[i].Timestamp < rows[j].Timestamp
	})
	return rows, nil
}

// GetThresholdExceedances returns readings strictly above threshold.
func GetThresholdExceedances(ctx context.Context, repo AnalyticsRepository, placeID string, sensorType SensorType, threshold float64) ([]SensorData, error) {
	rows, err := GetPlaceSensorData(ctx, repo, placeID, []SensorType{sensorType})
	if err != nil {
		return nil, err
	}
	out := make([]SensorData, 0, len(rows))
	for _, row := range rows {
		if row.Value > threshold {
			out = append(out, row)
		}
	}
	return out, nil
}