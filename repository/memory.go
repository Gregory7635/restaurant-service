// Пакет repository: in-memory хранилище.
package repository

import (
	"context"
	"sync"

	"example.com/restaurant-service/restaurant"
)

type memoryRepo struct {
	mu          sync.RWMutex
	restaurants map[string]restaurant.Restaurant
}

// NewMemoryRepository создаёт хранилище в памяти (map + мьютекс).
func NewMemoryRepository() restaurant.Repository {
	return &memoryRepo{restaurants: make(map[string]restaurant.Restaurant)}
}

func (m *memoryRepo) CreateRestaurant(ctx context.Context, r restaurant.Restaurant) error {
	m.mu.Lock()
	defer m.mu.Unlock()
	m.restaurants[r.ID] = r
	return nil
}

func (m *memoryRepo) GetRestaurantByID(ctx context.Context, id string) (restaurant.Restaurant, error) {
	m.mu.RLock()
	defer m.mu.RUnlock()
	r, ok := m.restaurants[id]
	if !ok {
		return restaurant.Restaurant{}, restaurant.ErrNotFound
	}
	return r, nil
}

func (m *memoryRepo) ChangeRestaurantRating(ctx context.Context, id string, rating float32) error {
	m.mu.Lock()
	defer m.mu.Unlock()
	r, ok := m.restaurants[id]
	if !ok {
		return restaurant.ErrNotFound
	}
	r.Rating = rating
	m.restaurants[id] = r
	return nil
}
