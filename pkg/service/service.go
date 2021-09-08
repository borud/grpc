package service

import (
	"context"
	"fmt"
	"strconv"
	"strings"
	"sync"

	"github.com/borud/grpc/pkg/apipb"
	"github.com/borud/grpc/pkg/auth"
	"github.com/borud/grpc/pkg/model"
	"github.com/grpc-ecosystem/grpc-gateway/v2/runtime" // Note: make sure this is v2
	"google.golang.org/grpc"
)

// Service is the API implementation
type Service struct {
	mu             sync.RWMutex
	lastID         uint64
	data           map[uint64]model.Sample
	grpcClientConn *grpc.ClientConn
}

// New creates a new service
func New() *Service {
	return &Service{
		mu:     sync.RWMutex{},
		lastID: 0,
		data:   map[uint64]model.Sample{},
	}
}

// CloseInternalClientConn closes the internal client connection if we use the RESTMuxViaGRPC()
// style calls.  While this isn't strictly necessary, being able to shut it down will clean things
// up so unit tests don't complain about leaked goroutines.
func (s *Service) CloseInternalClientConn() error {
	if s.grpcClientConn != nil {
		return s.grpcClientConn.Close()
	}
	return nil
}

// Return a gRPC server instance.
func (s *Service) GRPCServer() *grpc.Server {
	grpcServer := grpc.NewServer(
		grpc.ChainUnaryInterceptor(
			auth.AuthInterceptor,
			auth.LoggingServerInterceptor,
		),
	)

	apipb.RegisterSamplesServer(grpcServer, s)

	return grpcServer
}

// RESTMuxViaGRPC creates a mux which uses an internal gRPC client.  This has the benefit that
// REST accesses will go through interceptors.
func (s *Service) RESTMuxViaGRPC(ctx context.Context, listenAddr string) (*runtime.ServeMux, error) {
	mux := runtime.NewServeMux()

	// Get the port
	parts := strings.Split(listenAddr, ":")
	if len(parts) < 2 {
		return nil, fmt.Errorf("cannot get port from listenAddr '%s'", listenAddr)
	}
	port, err := strconv.ParseInt(parts[1], 10, 64)
	if err != nil {
		return nil, fmt.Errorf("error parsing port number in listenAddr %s: %w", listenAddr, err)
	}

	// Create the client
	c, err := grpc.Dial(fmt.Sprintf(":%d", port), grpc.WithInsecure())
	if err != nil {
		return nil, fmt.Errorf("unable to wire up internal REST client to gRPC interface: %w", err)
	}

	s.grpcClientConn = c

	// Register handler
	if err := apipb.RegisterSamplesHandler(ctx, mux, c); err != nil {
		return nil, fmt.Errorf("RegisterUsersHandler failed: %w", err)
	}

	return mux, nil
}

// RESTMux that uses the service directly.  This means requests will not go through interceptors.
func (s *Service) RESTMux(ctx context.Context) (*runtime.ServeMux, error) {
	mux := runtime.NewServeMux()

	if err := apipb.RegisterSamplesHandlerServer(ctx, mux, s); err != nil {
		return nil, fmt.Errorf("RegisterUsersHandler failed: %w", err)
	}

	return mux, nil
}
