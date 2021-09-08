package cookies

import (
	"context"
	"fmt"
	"log"
	"net/http"

	"github.com/borud/grpc/pkg/apipb"
	"github.com/grpc-ecosystem/grpc-gateway/v2/runtime"
	"google.golang.org/grpc"
	"google.golang.org/grpc/metadata"
	"google.golang.org/protobuf/reflect/protoreflect"
)

// CookieService is a service for demoing how to manage cookies
type CookieService struct {
	grpcListenAddr string
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

// Mux returns pointer to the ServeMux instance
func (s *CookieService) Mux() *runtime.ServeMux {
	return s.mux
}

// Start initializes the service
func (s *CookieService) Start() error {
	s.mux = runtime.NewServeMux(
		runtime.WithForwardResponseOption(gatewayResponseModifier),
		runtime.WithMetadata(gatewayMetadataAnnotator),
	)

	err := apipb.RegisterCookiesHandlerFromEndpoint(s.ctx, s.mux, s.grpcListenAddr, []grpc.DialOption{grpc.WithInsecure()})
	if err != nil {
		return err
	}

	return nil
}

// gatewayResponseModifier is called after
func gatewayResponseModifier(ctx context.Context, w http.ResponseWriter, msg protoreflect.ProtoMessage) error {
	log.Printf("called gatewayResponseModifier")

	md, ok := runtime.ServerMetadataFromContext(ctx)
	if !ok {
		return fmt.Errorf("failed to extract metadata from context")
	}

	// get cookie name and value
	name := md.HeaderMD.Get(cookieValueHeaderName)
	if len(name) == 0 {
		return nil
	}
	value := md.HeaderMD.Get(cookieValueHeaderValue)
	if len(value) == 0 {
		return nil
	}

	// then set the cookie
	http.SetCookie(w, &http.Cookie{
		Name:  name[0],
		Value: value[0],
	})

	return nil
}

// gatewayMetadataAnnotator is called before
func gatewayMetadataAnnotator(ctx context.Context, r *http.Request) metadata.MD {
	log.Printf("called gatewayMetadataAnnotator")

	c, err := r.Cookie("foobar")
	if err == http.ErrNoCookie {
		log.Printf("  -> no cookie for you")
		return nil
	}
	if err != nil {
		return metadata.Pairs()
	}

	log.Printf("  -> cookie: %s", c.String())

	return metadata.Pairs("foobar-cookie", c.Value)
}
