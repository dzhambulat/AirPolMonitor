package readmodel

func GetPlaceSensorData(placeId string, sensorType []SensorType) ([]SensorData, error) {
	return nil, nil
}

type SensorType int

const (
	SensorTypeAirQuality SensorType = iota
	SensorTypeTemperature
	SensorTypeHumidity
	SensorTypeCO2
	SensorTypeVOC
)

type SensorData struct {
	sensorType SensorType
}