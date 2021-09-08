package main

import (
	"log"
	"os"
	"time"

	"github.com/borud/grpc/pkg/client"
	"github.com/jessevdk/go-flags"
	"google.golang.org/grpc"
)

var opt struct {
	GRPCListenAddr string        `long:"grpc-listen" default:":3333" description:"grpc listen address" required:"yes"`
	Timeout        time.Duration `long:"timeout" default:"1s" description:"grpc call timeout"`
	Add            addCmd        `command:"add" description:"add sample"`
	List           listCmd       `command:"list" alias:"ls" description:"list samples"`
}

func main() {
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

func newClient() *client.Client {
	c, err := client.Dial(opt.GRPCListenAddr, grpc.WithInsecure())
	if err != nil {
		log.Fatalf("unable to connect to '%s': %v", opt.GRPCListenAddr, err)
	}
	return c
}
