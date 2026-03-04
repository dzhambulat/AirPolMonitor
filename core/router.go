package core

import (
	"AirPolMonitor/core/types"
	"context"
	"sync"
)


func Start(ctx context.Context, numWorkers int, sensorSource types.ISensorSource, airDataProducer types.IAirPolDataProducer) error {
	wg := sync.WaitGroup{}
	for i :=0; i < numWorkers; i++ {
		wg.Add(1)
		go func() {
			
			var airData chan types.AirData = make(chan types.AirData)
			var sensorData <-chan types.SensorData = sensorSource.Subscribe(ctx)
			
			defer close(airData)

			go func() {
				airDataProducer.Produce(ctx, airData)
			}()
			
			for {
				select {
				case data := <-sensorData:
					airData <- types.AirData{
						SensorID: data.SensorID,
						Value: data.Value,
						Timestamp: data.Timestamp,
					}
				case <-ctx.Done():
					wg.Done()
					return
				}
			}
		}()
	}
	wg.Wait()
	return nil
}