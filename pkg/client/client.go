package client

import (
	"github.com/borud/grpc/pkg/apipb"
	"google.golang.org/grpc"
)

// Client is a structure that contains clients for all the gRPC services in the
// project.  In this example we only have one, but if you have multiple services
// this comes in handy.
type Client struct {
	addr   string
	conn   *grpc.ClientConn
	Sample apipb.SamplesClient
}

// Dial the gRPC server
func Dial(addr string, opts ...grpc.DialOption) (*Client, error) {
	c, err := grpc.Dial(addr, opts...)
	if err != nil {
		return nil, err
	}

	return &Client{
		addr:   addr,
		conn:   c,
		Sample: apipb.NewSamplesClient(c),
	}, nil
}

// Close the client connection
func (c *Client) Close() error {
	return c.conn.Close()
}
