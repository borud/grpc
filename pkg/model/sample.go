package model

import "github.com/borud/grpc/pkg/apipb"

// Sample is a sample data structure
type Sample struct {
	ID   uint64
	Name string
	Data []byte
}

// Proto returns the protobuffer equivalent of the Sample
func (s Sample) Proto() *apipb.Sample {
	return &apipb.Sample{
		Id:   s.ID,
		Name: s.Name,
		Data: s.Data,
	}
}

// SampleFromProto returns Sample corresponding to apipb.Sample
func SampleFromProto(p *apipb.Sample) Sample {
	return Sample{
		ID:   p.Id,
		Name: p.Name,
		Data: p.Data,
	}
}
