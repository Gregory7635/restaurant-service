package main

import (
	"context"
	"flag"
	"net"
	"net/http"
	"os"
	"os/signal"
	"syscall"
	"time"

	"github.com/go-kit/log"
	"github.com/go-kit/log/level"
	"google.golang.org/grpc"
	"google.golang.org/grpc/reflection"

	"example.com/restaurant-service/grpctransport"
	"example.com/restaurant-service/httptransport"
	"example.com/restaurant-service/implementation"
	"example.com/restaurant-service/pb"
	"example.com/restaurant-service/repository"
	"example.com/restaurant-service/transport"
)

func main() {
	httpAddr := flag.String("http.addr", ":8080", "HTTP listen address")
	grpcAddr := flag.String("grpc.addr", ":9090", "gRPC listen address")
	flag.Parse()

	// Логгер
	var logger log.Logger
	{
		logger = log.NewLogfmtLogger(os.Stderr)
		logger = log.NewSyncLogger(logger)
		logger = level.NewFilter(logger, level.AllowDebug())
		logger = log.With(logger, "svc", "restaurant", "ts", log.DefaultTimestampUTC, "caller", log.DefaultCaller)
	}
	_ = level.Info(logger).Log("msg", "service started")
	defer func() { _ = level.Info(logger).Log("msg", "service ended") }()

	// Сборка слоёв: репозиторий -> сервис -> эндпоинты -> транспорты
	repo := repository.NewMemoryRepository()
	svc := implementation.NewService(repo, logger)
	endpoints := transport.MakeEndpoints(svc)

	// HTTP-сервер
	httpServer := &http.Server{
		Addr:    *httpAddr,
		Handler: httptransport.NewHandler(endpoints, logger),
	}

	// gRPC-сервер
	grpcServer := grpc.NewServer()
	pb.RegisterRestaurantServiceServer(grpcServer, grpctransport.NewServer(endpoints))
	reflection.Register(grpcServer) // чтобы grpcurl работал без .proto-файла

	errs := make(chan error, 3)

	go func() {
		c := make(chan os.Signal, 1)
		signal.Notify(c, syscall.SIGINT, syscall.SIGTERM)
		sig := <-c
		errs <- &signalError{sig.String()}
	}()

	go func() {
		_ = level.Info(logger).Log("transport", "HTTP", "addr", *httpAddr)
		errs <- httpServer.ListenAndServe()
	}()

	go func() {
		lis, err := net.Listen("tcp", *grpcAddr)
		if err != nil {
			errs <- err
			return
		}
		_ = level.Info(logger).Log("transport", "gRPC", "addr", *grpcAddr)
		errs <- grpcServer.Serve(lis)
	}()

	_ = level.Error(logger).Log("exit", <-errs)

	// Graceful shutdown
	ctx, cancel := context.WithTimeout(context.Background(), 5*time.Second)
	defer cancel()
	_ = httpServer.Shutdown(ctx)
	grpcServer.GracefulStop()
}

type signalError struct{ s string }

func (e *signalError) Error() string { return "received signal: " + e.s }
