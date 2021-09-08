package auth

import (
	"context"
	"log"

	"google.golang.org/grpc"
	"google.golang.org/grpc/metadata"
)

func LoggingServerInterceptor(ctx context.Context, req interface{}, info *grpc.UnaryServerInfo, handler grpc.UnaryHandler) (resp interface{}, err error) {
	log.Printf("[BEFORE] server=%v method=%s req=%+v", info.Server, info.FullMethod, req)

	// Remember to call the original request
	response, err := handler(ctx, req)

	// This line is logged after we handed off the request to where it was supposed to go
	log.Printf("[ AFTER] response=%+v", response)

	return response, err
}

func AuthInterceptor(ctx context.Context, req interface{}, info *grpc.UnaryServerInfo, handler grpc.UnaryHandler) (resp interface{}, err error) {
	md, ok := metadata.FromIncomingContext(ctx)
	if !ok {
		log.Printf("metadata.FromIncomingContext was not OK")
	}

	log.Printf("[AUTH] called %#v", md)

	return handler(ctx, req)
}
