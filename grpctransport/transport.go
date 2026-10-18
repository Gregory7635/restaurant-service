// Пакет grpctransport: gRPC-транспорт (те же эндпоинты, что и у HTTP).
package grpctransport

import (
	"context"
	"errors"

	"github.com/go-kit/kit/endpoint"
	kitgrpc "github.com/go-kit/kit/transport/grpc"
	"google.golang.org/grpc/codes"
	"google.golang.org/grpc/status"

	"example.com/restaurant-service/pb"
	"example.com/restaurant-service/restaurant"
	"example.com/restaurant-service/transport"
)

type grpcServer struct {
	pb.UnimplementedRestaurantServiceServer

	create       kitgrpc.Handler
	getByID      kitgrpc.Handler
	changeRating kitgrpc.Handler
}

// NewServer создаёт реализацию pb.RestaurantServiceServer поверх эндпоинтов.
func NewServer(e transport.Endpoints) pb.RestaurantServiceServer {
	return &grpcServer{
		create:       kitgrpc.NewServer(e.Create, decodeCreateRequest, encodeCreateResponse),
		getByID:      kitgrpc.NewServer(e.GetByID, decodeGetByIDRequest, encodeGetByIDResponse),
		changeRating: kitgrpc.NewServer(e.ChangeRating, decodeChangeRatingRequest, encodeChangeRatingResponse),
	}
}

func (s *grpcServer) Create(ctx context.Context, req *pb.CreateRequest) (*pb.CreateResponse, error) {
	_, resp, err := s.create.ServeGRPC(ctx, req)
	if err != nil {
		return nil, err
	}
	return resp.(*pb.CreateResponse), nil
}

func (s *grpcServer) GetByID(ctx context.Context, req *pb.GetByIDRequest) (*pb.GetByIDResponse, error) {
	_, resp, err := s.getByID.ServeGRPC(ctx, req)
	if err != nil {
		return nil, err
	}
	return resp.(*pb.GetByIDResponse), nil
}

func (s *grpcServer) ChangeRating(ctx context.Context, req *pb.ChangeRatingRequest) (*pb.ChangeRatingResponse, error) {
	_, resp, err := s.changeRating.ServeGRPC(ctx, req)
	if err != nil {
		return nil, err
	}
	return resp.(*pb.ChangeRatingResponse), nil
}

// ---------- decode: protobuf -> DTO ----------

func decodeCreateRequest(_ context.Context, request interface{}) (interface{}, error) {
	req := request.(*pb.CreateRequest)
	r := req.GetRestaurant()
	return transport.CreateRequest{Restaurant: restaurant.Restaurant{
		Name:    r.GetName(),
		Address: r.GetAddress(),
		Rating:  r.GetRating(),
		Open:    r.GetOpen(),
	}}, nil
}

func decodeGetByIDRequest(_ context.Context, request interface{}) (interface{}, error) {
	req := request.(*pb.GetByIDRequest)
	return transport.GetByIDRequest{ID: req.GetId()}, nil
}

func decodeChangeRatingRequest(_ context.Context, request interface{}) (interface{}, error) {
	req := request.(*pb.ChangeRatingRequest)
	return transport.ChangeRatingRequest{ID: req.GetId(), Rating: req.GetRating()}, nil
}

// ---------- encode: DTO -> protobuf ----------

func encodeCreateResponse(_ context.Context, response interface{}) (interface{}, error) {
	if err := failure(response); err != nil {
		return nil, err
	}
	resp := response.(transport.CreateResponse)
	return &pb.CreateResponse{Id: resp.ID}, nil
}

func encodeGetByIDResponse(_ context.Context, response interface{}) (interface{}, error) {
	if err := failure(response); err != nil {
		return nil, err
	}
	r := response.(transport.GetByIDResponse).Restaurant
	return &pb.GetByIDResponse{Restaurant: &pb.Restaurant{
		Id:      r.ID,
		Name:    r.Name,
		Address: r.Address,
		Rating:  r.Rating,
		Open:    r.Open,
	}}, nil
}

func encodeChangeRatingResponse(_ context.Context, response interface{}) (interface{}, error) {
	if err := failure(response); err != nil {
		return nil, err
	}
	return &pb.ChangeRatingResponse{}, nil
}

// failure превращает бизнес-ошибку из DTO в gRPC-статус.
func failure(response interface{}) error {
	f, ok := response.(endpoint.Failer)
	if !ok || f.Failed() == nil {
		return nil
	}
	err := f.Failed()
	switch {
	case errors.Is(err, restaurant.ErrNotFound):
		return status.Error(codes.NotFound, err.Error())
	case errors.Is(err, restaurant.ErrInvalidInput), errors.Is(err, restaurant.ErrInvalidRating):
		return status.Error(codes.InvalidArgument, err.Error())
	default:
		return status.Error(codes.Internal, err.Error())
	}
}
