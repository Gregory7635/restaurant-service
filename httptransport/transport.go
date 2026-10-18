// Пакет httptransport: HTTP-транспорт.
package httptransport

import (
	"context"
	"encoding/json"
	"errors"
	"net/http"

	"github.com/go-kit/kit/endpoint"
	kithttp "github.com/go-kit/kit/transport/http"
	"github.com/go-kit/log"
	"github.com/gorilla/mux"

	"example.com/restaurant-service/restaurant"
	"example.com/restaurant-service/transport"
)

// NewHandler связывает эндпоинты Go kit с HTTP.
func NewHandler(endpoints transport.Endpoints, logger log.Logger) http.Handler {
	r := mux.NewRouter()

	options := []kithttp.ServerOption{
		kithttp.ServerErrorHandler(transportErrorHandler{logger}),
		kithttp.ServerErrorEncoder(encodeError),
	}

	r.Methods("POST").Path("/restaurants").Handler(kithttp.NewServer(
		endpoints.Create, decodeCreateRequest, encodeResponse, options...))

	r.Methods("GET").Path("/restaurants/{id}").Handler(kithttp.NewServer(
		endpoints.GetByID, decodeGetByIDRequest, encodeResponse, options...))

	r.Methods("PATCH").Path("/restaurants/{id}/rating").Handler(kithttp.NewServer(
		endpoints.ChangeRating, decodeChangeRatingRequest, encodeResponse, options...))

	return r
}

type transportErrorHandler struct{ logger log.Logger }

func (h transportErrorHandler) Handle(_ context.Context, err error) {
	h.logger.Log("transport_error", err)
}

func decodeCreateRequest(_ context.Context, r *http.Request) (interface{}, error) {
	var req transport.CreateRequest
	if err := json.NewDecoder(r.Body).Decode(&req.Restaurant); err != nil {
		return nil, err
	}
	return req, nil
}

func decodeGetByIDRequest(_ context.Context, r *http.Request) (interface{}, error) {
	return transport.GetByIDRequest{ID: mux.Vars(r)["id"]}, nil
}

func decodeChangeRatingRequest(_ context.Context, r *http.Request) (interface{}, error) {
	var body struct {
		Rating float32 `json:"rating"`
	}
	if err := json.NewDecoder(r.Body).Decode(&body); err != nil {
		return nil, err
	}
	return transport.ChangeRatingRequest{ID: mux.Vars(r)["id"], Rating: body.Rating}, nil
}

// encodeResponse: если сервис вернул бизнес-ошибку, отдаём её с нужным HTTP-статусом.
func encodeResponse(ctx context.Context, w http.ResponseWriter, response interface{}) error {
	if f, ok := response.(endpoint.Failer); ok && f.Failed() != nil {
		encodeError(ctx, f.Failed(), w)
		return nil
	}
	w.Header().Set("Content-Type", "application/json; charset=utf-8")
	return json.NewEncoder(w).Encode(response)
}

func encodeError(_ context.Context, err error, w http.ResponseWriter) {
	w.Header().Set("Content-Type", "application/json; charset=utf-8")
	w.WriteHeader(statusFor(err))
	json.NewEncoder(w).Encode(map[string]string{"error": err.Error()})
}

func statusFor(err error) int {
	switch {
	case errors.Is(err, restaurant.ErrNotFound):
		return http.StatusNotFound
	case errors.Is(err, restaurant.ErrInvalidInput), errors.Is(err, restaurant.ErrInvalidRating):
		return http.StatusBadRequest
	default:
		var syn *json.SyntaxError
		var ute *json.UnmarshalTypeError
		if errors.As(err, &syn) || errors.As(err, &ute) {
			return http.StatusBadRequest
		}
		return http.StatusInternalServerError
	}
}
