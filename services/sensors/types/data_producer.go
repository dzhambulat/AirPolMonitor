package types

import "context"

type AirData struct {
	SensorID string
	Value    float32
	Timestamp int64
}

type IAirPolDataProducer interface {
	Produce(ctx context.Context, data <-chan AirData) error
}