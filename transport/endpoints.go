// Пакет transport: DTO и эндпоинты Go kit (не зависят от протокола).
package transport

import (
	"context"

	"github.com/go-kit/kit/endpoint"

	"example.com/restaurant-service/restaurant"
)

// ---------- DTO ----------

type CreateRequest struct {
	Restaurant restaurant.Restaurant
}

type CreateResponse struct {
	ID  string `json:"id,omitempty"`
	Err error  `json:"-"`
}

func (r CreateResponse) Failed() error { return r.Err }

type GetByIDRequest struct {
	ID string
}

type GetByIDResponse struct {
	Restaurant restaurant.Restaurant `json:"restaurant"`
	Err        error                 `json:"-"`
}

func (r GetByIDResponse) Failed() error { return r.Err }

type ChangeRatingRequest struct {
	ID     string
	Rating float32
}

type ChangeRatingResponse struct {
	Status string `json:"status,omitempty"`
	Err    error  `json:"-"`
}

func (r ChangeRatingResponse) Failed() error { return r.Err }

// ---------- Эндпоинты ----------

type Endpoints struct {
	Create       endpoint.Endpoint
	GetByID      endpoint.Endpoint
	ChangeRating endpoint.Endpoint
}

func MakeEndpoints(s restaurant.Service) Endpoints {
	return Endpoints{
		Create:       makeCreateEndpoint(s),
		GetByID:      makeGetByIDEndpoint(s),
		ChangeRating: makeChangeRatingEndpoint(s),
	}
}

func makeCreateEndpoint(s restaurant.Service) endpoint.Endpoint {
	return func(ctx context.Context, request interface{}) (interface{}, error) {
		req := request.(CreateRequest)
		id, err := s.Create(ctx, req.Restaurant)
		return CreateResponse{ID: id, Err: err}, nil
	}
}

func makeGetByIDEndpoint(s restaurant.Service) endpoint.Endpoint {
	return func(ctx context.Context, request interface{}) (interface{}, error) {
		req := request.(GetByIDRequest)
		r, err := s.GetByID(ctx, req.ID)
		return GetByIDResponse{Restaurant: r, Err: err}, nil
	}
}

func makeChangeRatingEndpoint(s restaurant.Service) endpoint.Endpoint {
	return func(ctx context.Context, request interface{}) (interface{}, error) {
		req := request.(ChangeRatingRequest)
		err := s.ChangeRating(ctx, req.ID, req.Rating)
		if err != nil {
			return ChangeRatingResponse{Err: err}, nil
		}
		return ChangeRatingResponse{Status: "ok"}, nil
	}
}
