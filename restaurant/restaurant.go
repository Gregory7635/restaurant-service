// Пакет restaurant: сущности предметной области и интерфейсы (Вариант 5).
package restaurant

import "context"

// Restaurant - сущность «Ресторан».
type Restaurant struct {
	ID      string  `json:"id,omitempty"`
	Name    string  `json:"name"`
	Address string  `json:"address"`
	Rating  float32 `json:"rating"`
	Open    bool    `json:"open"`
}

// Service - бизнес-логика сервиса (ничего не знает про HTTP/gRPC).
type Service interface {
	Create(ctx context.Context, restaurant Restaurant) (string, error)
	GetByID(ctx context.Context, id string) (Restaurant, error)
	ChangeRating(ctx context.Context, id string, rating float32) error
}

// Repository - хранилище ресторанов.
type Repository interface {
	CreateRestaurant(ctx context.Context, r Restaurant) error
	GetRestaurantByID(ctx context.Context, id string) (Restaurant, error)
	ChangeRestaurantRating(ctx context.Context, id string, rating float32) error
}
