package service

import (
	"context"

	"github.com/cldcvr/terrarium/src/api/internal/config"
	"github.com/cldcvr/terrarium/src/api/service/terrariumsrv"
	"github.com/cldcvr/terrarium/src/pkg/pb/terrariumpb"
	"google.golang.org/protobuf/types/known/emptypb"
)

type Service interface {
	HealthCheck(ctx context.Context, req *emptypb.Empty) (*emptypb.Empty, error)
	ListModules(ctx context.Context, req *terrariumpb.ListModulesRequest) (*terrariumpb.ListModulesResponse, error)
	CodeCompletion(ctx context.Context, req *terrariumpb.CompletionRequest) (*terrariumpb.CompletionResponse, error)
	ListModuleAttributes(ctx context.Context, req *terrariumpb.ListModuleAttributesRequest) (*terrariumpb.ListModuleAttributesResponse, error)
}

func New() (Service, error) {
	d, err := config.DBConnect()
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
