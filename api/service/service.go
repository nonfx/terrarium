package service

import (
	"context"

	"github.com/cldcvr/terrarium/api/db"
	"github.com/cldcvr/terrarium/api/pkg/pb/terrariumpb"
	"github.com/cldcvr/terrarium/api/service/terrariumsrv"
	"google.golang.org/protobuf/types/known/emptypb"
)

type Service interface {
	HealthCheck(ctx context.Context, req *emptypb.Empty) (*emptypb.Empty, error)
	ListModules(ctx context.Context, req *terrariumpb.ListModulesRequest) (*terrariumpb.ListModulesResponse, error)
	CodeCompletion(ctx context.Context, req *terrariumpb.CompletionRequest) (*terrariumpb.CompletionResponse, error)
}

func New() (Service, error) {
	d, err := db.Connect()
	if err != nil {
		return nil, err
	}

	s := terrariumsrv.New(d)

	return &struct {
		*healthservice
		*terrariumsrv.Service
	}{
		healthservice: &healthservice{},
		Service:       s,
	}, nil
}

type healthservice struct{}

func (s healthservice) HealthCheck(ctx context.Context, req *emptypb.Empty) (*emptypb.Empty, error) {
	return &emptypb.Empty{}, nil
}
