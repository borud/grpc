package main

import (
	"context"
	"fmt"
	"log"

	"github.com/golang/protobuf/ptypes/empty"
)

type listCmd struct{}

func (r *listCmd) Execute([]string) error {
	client := newClient()

	ctx, cancel := context.WithTimeout(context.Background(), opt.Timeout)
	defer cancel()

	response, err := client.Sample.List(ctx, &empty.Empty{})
	if err != nil {
		return fmt.Errorf("unable to list: %v", err)
	}

	for _, sample := range response.Samples {
		log.Printf("id='%d' name='%s' data='%s'", sample.Id, sample.Name, sample.Data)
	}

	return nil
}
