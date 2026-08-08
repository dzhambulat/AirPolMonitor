package sensors

import (
	"AirPolMonitor/services/sensors/types"
	"context"
	"errors"
	"sync"
)


func Start(ctx context.Context, numWorkers int, sensorSource types.ISensorSource, airDataProducer types.IAirPolDataProducer) error {
	wg := sync.WaitGroup{}
	var sensorData <-chan types.SensorData = sensorSource.Subscribe(ctx)
	if sensorData == nil {
		return errors.New("sensor data source is nil")
	}
	for i :=0; i < numWorkers; i++ {
		wg.Add(1)
		go func() {
			
			var airData chan types.AirData = make(chan types.AirData)
			
			defer close(airData)
			defer wg.Done()
			
			go func() {
				airDataProducer.Produce(ctx, airData)
			}()
			
			for {
				select {
				case data, ok := <-sensorData:
					if !ok {
						return
					}
					airData <- types.AirData{
						SensorID: data.SensorID,
						Value: data.Value,
						Timestamp: data.Timestamp,
					}
				case <-ctx.Done():
					return
				}
			}
		}()
	}
	wg.Wait()
	return nil
}