package implementation

import (
	"context"
	"errors"
	"testing"

	"github.com/go-kit/log"

	"example.com/restaurant-service/repository"
	"example.com/restaurant-service/restaurant"
)

func newTestService() restaurant.Service {
	return NewService(repository.NewMemoryRepository(), log.NewNopLogger())
}

func TestCreateAndGet(t *testing.T) {
	svc := newTestService()
	ctx := context.Background()

	id, err := svc.Create(ctx, restaurant.Restaurant{Name: "Пицца-Хаус", Address: "ул. Ленина, 1", Rating: 4.2, Open: true})
	if err != nil {
		t.Fatalf("Create: unexpected error: %v", err)
	}
	if id == "" {
		t.Fatal("Create: empty id")
	}

	r, err := svc.GetByID(ctx, id)
	if err != nil {
		t.Fatalf("GetByID: unexpected error: %v", err)
	}
	if r.Name != "Пицца-Хаус" || r.Rating != 4.2 || !r.Open {
		t.Errorf("GetByID: got %+v", r)
	}
}

func TestCreateValidation(t *testing.T) {
	svc := newTestService()
	tests := []struct {
		name string
		in   restaurant.Restaurant
		want error
	}{
		{"empty name", restaurant.Restaurant{Address: "addr", Rating: 3}, restaurant.ErrInvalidInput},
		{"empty address", restaurant.Restaurant{Name: "name", Rating: 3}, restaurant.ErrInvalidInput},
		{"rating too high", restaurant.Restaurant{Name: "n", Address: "a", Rating: 5.1}, restaurant.ErrInvalidRating},
		{"negative rating", restaurant.Restaurant{Name: "n", Address: "a", Rating: -1}, restaurant.ErrInvalidRating},
	}
	for _, tt := range tests {
		t.Run(tt.name, func(t *testing.T) {
			if _, err := svc.Create(context.Background(), tt.in); !errors.Is(err, tt.want) {
				t.Errorf("got %v, want %v", err, tt.want)
			}
		})
	}
}

func TestGetByIDNotFound(t *testing.T) {
	svc := newTestService()
	if _, err := svc.GetByID(context.Background(), "missing"); !errors.Is(err, restaurant.ErrNotFound) {
		t.Errorf("got %v, want %v", err, restaurant.ErrNotFound)
	}
}

func TestChangeRating(t *testing.T) {
	svc := newTestService()
	ctx := context.Background()
	id, _ := svc.Create(ctx, restaurant.Restaurant{Name: "Суши", Address: "пр. Мира, 5", Rating: 4})

	if err := svc.ChangeRating(ctx, id, 4.9); err != nil {
		t.Fatalf("ChangeRating: unexpected error: %v", err)
	}
	r, _ := svc.GetByID(ctx, id)
	if r.Rating != 4.9 {
		t.Errorf("rating = %v, want 4.9", r.Rating)
	}

	if err := svc.ChangeRating(ctx, id, 7); !errors.Is(err, restaurant.ErrInvalidRating) {
		t.Errorf("got %v, want %v", err, restaurant.ErrInvalidRating)
	}
	if err := svc.ChangeRating(ctx, "missing", 3); !errors.Is(err, restaurant.ErrNotFound) {
		t.Errorf("got %v, want %v", err, restaurant.ErrNotFound)
	}
}
