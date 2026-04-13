package infrastructure

import (
	"AirPolMonitor/services/places/domain"
	"context"
	"sync"
)

type SensorData struct {
	SensorType domain.SensorType
	Value float32
	Timestamp int64
}

type  MemoryPlaceSensorsRepository struct {
	data map[string][]SensorData
	mu sync.RWMutex
}

func NewMemoryPlaceSensorsRepository() *MemoryPlaceSensorsRepository {
	return &MemoryPlaceSensorsRepository{data: make(map[string][]SensorData)}
}

func (m *MemoryPlaceSensorsRepository) Save(ctx context.Context, placeID string, data SensorData) error {
	m.mu.Lock()
	defer m.mu.Unlock()
	m.data[placeID] = append(m.data[placeID], data)
	return nil
}

func (m *MemoryPlaceSensorsRepository) Get(ctx context.Context, placeID string) ([]SensorData, error) {
	m.mu.RLock()
	defer m.mu.RUnlock()
	return m.data[placeID], nil
}