package cookies

import (
	"context"
	"log"

	"github.com/borud/grpc/pkg/apipb"
	"github.com/golang/protobuf/ptypes/empty"
)

func (s *CookieService) SetCookie(ctx context.Context, cookie *apipb.Cookie) (*empty.Empty, error) {
	log.Printf("SetCookie '%s' = '%s'", cookie.Name, cookie.Value)
	return &empty.Empty{}, nil
}

func (s *CookieService) GetCookie(ctx context.Context, req *apipb.GetCookieRequest) (*apipb.Cookie, error) {
	log.Printf("GetCookie '%s'", req.Name)
	return &apipb.Cookie{}, nil
}

func (s *CookieService) DeleteCookie(ctx context.Context, req *apipb.DeleteCookieRequest) (*empty.Empty, error) {
	log.Printf("DeleteCookie '%s'", req.Name)
	return &empty.Empty{}, nil
}
