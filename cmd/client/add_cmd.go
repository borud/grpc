package main

import (
	"context"
	"fmt"

	"github.com/borud/grpc/pkg/apipb"
)

type addCmd struct {
	Name string `long:"name" required:"yes"`
	Data string `long:"data" required:"yes"`
}

func (r *addCmd) Execute([]string) error {
	client := newClient()

	ctx, cancel := context.WithTimeout(context.Background(), opt.Timeout)
	defer cancel()

	response, err := client.Sample.Create(ctx, &apipb.Sample{
		Name: r.Name,
		Data: []byte(r.Data),
	})
	if err != nil {
		return fmt.Errorf("Unable to create: %v", err)
	}

	fmt.Printf("added %d\n", response.Id)

	return nil
}
