package main

import (
	"log"
	"os"
	"os/signal"
	"syscall"

	"github.com/borud/grpc/pkg/server"
	"github.com/jessevdk/go-flags"
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
	server := server.New(server.Config{
		GRPCListenAddr: opt.GRPCListenAddr,
		HTTPListenAddr: opt.HTTPListenAddr,
	})
	err := server.Start()
	if err != nil {
		log.Fatalf("error starting server: %v", err)
	}

	// Block forever
	// Capture Ctrl-C
	c := make(chan os.Signal, 1)
	signal.Notify(c, os.Interrupt, syscall.SIGTERM)
	<-c

	server.Shutdown()
}
