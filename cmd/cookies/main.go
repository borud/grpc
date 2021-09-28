package main

import (
	"log"
	"net"
	"net/http"
	"os"

	"github.com/borud/grpc/pkg/apipb"
	"github.com/borud/grpc/pkg/cookies"
	"github.com/jessevdk/go-flags"
	"google.golang.org/grpc"
)

var opt struct {
	GRPCListenAddr string `long:"grpc-listen" default:":3333" description:"grpc listen address" required:"yes"`
	HTTPListenAddr string `long:"http-listen" default:":4444" description:"http listen address" required:"yes"`
}

func init() {
	p := flags.NewParser(&opt, flags.Default)
	_, err := p.Parse()
	switch flagsErr := err.(type) {
	case flags.ErrorType:
		if flagsErr == flags.ErrHelp {
			os.Exit(0)
		}
		os.Exit(1)

	default:
		// os.Exit(1)
	}
}

func main() {
	// Set up listener
	listener, err := net.Listen("tcp", opt.GRPCListenAddr)
	if err != nil {
		log.Fatalf("error listening to %s: %v", opt.GRPCListenAddr, err)
	}

	// Create cookieservice
	cookieService := cookies.New(opt.GRPCListenAddr)
	err = cookieService.Start()
	if err != nil {
		log.Fatalf("error starting cookie service: %v", err)
	}

	// Create and start gRPC server
	grpcServer := grpc.NewServer()
	apipb.RegisterCookiesServer(grpcServer, cookieService)

	go func() {
		log.Printf("grpc server returned %v", grpcServer.Serve(listener))
	}()
	log.Printf("grpc: %s", listener.Addr().String())

	// Create mux
	mux := cookieService.Mux()


	// Set up web interface
	httpServer := &http.Server{
		Addr:    opt.HTTPListenAddr,
		Handler: mux,
	}

	log.Printf("http: %s", opt.HTTPListenAddr)
	httpServer.ListenAndServe()
}
