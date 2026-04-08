package domain

type PlaceRepositoryPort interface {

}

type PlaceAggregate struct {
	id string
	name string
	coordinates []float64
	
}

func (p *PlaceAggregate) addSensorData() error {
	return nil
}

func AddSensorDataToPlace(placeId string) error {
	return nil
}