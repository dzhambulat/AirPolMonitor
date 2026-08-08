package infrastructure

import (
	"AirPolMonitor/services/places/domain"
	"context"
	"fmt"
	"sync"

	"github.com/google/uuid"
)

// MemoryPlaceRepository is an in-process store for development.
type MemoryPlaceRepository struct {
	mu    sync.Mutex
	places map[string]*domain.PlaceAggregate
}

func NewMemoryPlaceRepository() *MemoryPlaceRepository {
	return &MemoryPlaceRepository{places: make(map[string]*domain.PlaceAggregate)}
}

func (m *MemoryPlaceRepository) SeedPlace(agg *domain.PlaceAggregate) {
	m.mu.Lock()
	defer m.mu.Unlock()
	m.places[agg.Id] = agg
}

func (m *MemoryPlaceRepository) Get(ctx context.Context, placeID string) (*domain.PlaceAggregate, error) {
	_ = ctx
	m.mu.Lock()
	defer m.mu.Unlock()
	agg, ok := m.places[placeID]
	if !ok {
		return nil, fmt.Errorf("place not found: %s", placeID)
	}
	return agg, nil
}

func (m *MemoryPlaceRepository) Save(ctx context.Context, agg *domain.PlaceAggregate) error {
	_ = ctx
	if agg == nil {
		return fmt.Errorf("aggregate is nil")
	}
	m.mu.Lock()
	defer m.mu.Unlock()

	if agg.Id == "" {
		agg.Id = uuid.New().String()
	}
	m.places[agg.Id] = agg
	return nil
}
