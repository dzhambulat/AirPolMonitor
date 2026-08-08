package types

import "context"

const (
	SensorTypeAirQuality int = iota
	SensorTypeTemperature
	SensorTypeHumidity
	SensorTypeCO2
	SensorTypeVOC
)

type SensorData struct {
	SensorID string
	SensorType int
	Value    float32
	Timestamp int64
}

type ISensorSource interface {
	Subscribe(ctx context.Context) <-chan SensorData
}
