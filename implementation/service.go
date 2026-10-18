// Пакет implementation: реализация сервисного слоя.
package implementation

import (
	"context"
	"errors"

	"github.com/go-kit/log"
	"github.com/go-kit/log/level"
	"github.com/google/uuid"

	"example.com/restaurant-service/restaurant"
)

type service struct {
	repository restaurant.Repository
	logger     log.Logger
}

// NewService - конструктор сервиса ресторанов.
func NewService(repo restaurant.Repository, logger log.Logger) restaurant.Service {
	return &service{repository: repo, logger: logger}
}

func validRating(r float32) bool { return r >= 0 && r <= 5 }

func (s *service) Create(ctx context.Context, r restaurant.Restaurant) (string, error) {
	logger := log.With(s.logger, "method", "Create")

	if r.Name == "" || r.Address == "" {
		return "", restaurant.ErrInvalidInput
	}
	if !validRating(r.Rating) {
		return "", restaurant.ErrInvalidRating
	}

	r.ID = uuid.New().String()

	if err := s.repository.CreateRestaurant(ctx, r); err != nil {
		level.Error(logger).Log("err", err)
		return "", restaurant.ErrCmdRepository
	}
	return r.ID, nil
}

func (s *service) GetByID(ctx context.Context, id string) (restaurant.Restaurant, error) {
	logger := log.With(s.logger, "method", "GetByID")

	r, err := s.repository.GetRestaurantByID(ctx, id)
	if err != nil {
		level.Error(logger).Log("err", err)
		if errors.Is(err, restaurant.ErrNotFound) {
			return restaurant.Restaurant{}, restaurant.ErrNotFound
		}
		return restaurant.Restaurant{}, restaurant.ErrQueryRepository
	}
	return r, nil
}

func (s *service) ChangeRating(ctx context.Context, id string, rating float32) error {
	logger := log.With(s.logger, "method", "ChangeRating")

	if !validRating(rating) {
		return restaurant.ErrInvalidRating
	}
	if err := s.repository.ChangeRestaurantRating(ctx, id, rating); err != nil {
		level.Error(logger).Log("err", err)
		if errors.Is(err, restaurant.ErrNotFound) {
			return restaurant.ErrNotFound
		}
		return restaurant.ErrCmdRepository
	}
	return nil
}
