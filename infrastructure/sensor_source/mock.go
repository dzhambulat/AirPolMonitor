package sensorsource

import (
	"AirPolMonitor/core/types"
	"context"
	"log"
	"time"
)


type MockSensorSource struct{}

func (m *MockSensorSource) Subscribe(ctx context.Context) <-chan types.SensorData {
	channel := make(chan types.SensorData)
	go func() {
		defer close(channel)
		var data []int = []int{1, 2, 3, 4, 5, 6, 7, 8, 9, 10}
		for _, d := range data {
			log.Println("Sending sensor data:", d)
			channel <- types.SensorData{
				SensorID: "mock-sensor",
				SensorType: types.SensorTypeAirQuality,
				Value: float32(d),
				Timestamp: time.Now().Unix(),
			}
			time.Sleep(1 * time.Second)
		}
	}()
	
	go func() {
		<-ctx.Done()
		select {
			case _, ok := <-channel:
				if ok {
					close(channel)
				}
				return
			default:
				close(channel)
		}
	}()

	return channel
}

func NewMockSensorSource() *MockSensorSource {
	return &MockSensorSource{}
}
