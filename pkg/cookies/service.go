package cookies

import (
	"context"
	"errors"

	"github.com/borud/grpc/pkg/apipb"
	"github.com/grpc-ecosystem/grpc-gateway/v2/runtime"
	"google.golang.org/grpc"
)

// CookieService is a service for demoing how to manage cookies
type CookieService struct {
	grpcListenAddr string
	grpcClientConn *grpc.ClientConn
	mux            *runtime.ServeMux
	ctx            context.Context
}

// New creates a new service instance
func New(grpcListenAddr string) *CookieService {
	return &CookieService{
		ctx:            context.Background(),
		grpcListenAddr: grpcListenAddr,
	}
}

// Start initializes the service
func (s *CookieService) Start() error {
	s.mux = runtime.NewServeMux()

	// Create internal connection to gRPC interface
	c, err := grpc.Dial(s.grpcListenAddr, grpc.WithInsecure())
	if err != nil {
		return err
	}
	s.grpcClientConn = c

	err = apipb.RegisterCookiesHandler(s.ctx, s.mux, c)
	if err != nil {
		return err
	}

	return nil
}

// Shutdown the service
func (s *CookieService) Shutdown() error {
	if s.grpcClientConn == nil {
		return errors.New("service was never started")
	}
	return s.grpcClientConn.Close()
}

// Mux returns pointer to the ServeMux instance
func (s *CookieService) Mux() *runtime.ServeMux {
	return s.mux
}
