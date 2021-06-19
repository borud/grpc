package service

import (
	"context"
	"sync/atomic"

	"github.com/borud/grpc/pkg/apipb"
	"github.com/borud/grpc/pkg/model"
	"github.com/golang/protobuf/ptypes/empty"
	"google.golang.org/grpc/codes"
	"google.golang.org/grpc/status"
	"google.golang.org/protobuf/types/known/emptypb"
)

func (s *Service) Create(ctx context.Context, sample *apipb.Sample) (*apipb.SampleCreateResponse, error) {
	s.mu.Lock()
	defer s.mu.Unlock()

	sample.Id = atomic.AddUint64(&s.lastID, 1)
	s.data[sample.Id] = model.SampleFromProto(sample)

	return &apipb.SampleCreateResponse{Id: sample.Id}, nil
}

func (s *Service) Get(ctx context.Context, req *apipb.SampleGetRequest) (*apipb.Sample, error) {
	s.mu.RLock()
	defer s.mu.RUnlock()

	sample, ok := s.data[req.Id]
	if !ok {
		return nil, status.Error(codes.NotFound, "not found")
	}

	return sample.Proto(), nil
}

func (s *Service) Update(ctx context.Context, req *apipb.Sample) (*empty.Empty, error) {
	s.mu.Lock()
	defer s.mu.Unlock()

	_, ok := s.data[req.Id]
	if !ok {
		return nil, status.Error(codes.NotFound, "not found")
	}
	s.data[req.Id] = model.SampleFromProto(req)
	return &empty.Empty{}, nil
}

func (s *Service) Delete(ctx context.Context, req *apipb.SampleDeleteRequest) (*empty.Empty, error) {
	s.mu.Lock()
	defer s.mu.Unlock()

	_, ok := s.data[req.Id]
	if !ok {
		return nil, status.Error(codes.NotFound, "not found")
	}
	delete(s.data, req.Id)
	return &emptypb.Empty{}, nil
}

func (s *Service) List(ctx context.Context, req *empty.Empty) (*apipb.SampleListResponse, error) {
	s.mu.RLock()
	defer s.mu.RUnlock()

	list := make([]*apipb.Sample, len(s.data))
	for _, v := range s.data {
		list = append(list, v.Proto())
	}

	return &apipb.SampleListResponse{Samples: list}, nil
}
