package service

import (
	"context"
	"time"

	"github.com/cldcvr/terrarium/api/pkg/pb/terrariumpb"
	"google.golang.org/protobuf/types/known/emptypb"
)

type Service interface {
	HealthCheck(ctx context.Context, req *emptypb.Empty) (*emptypb.Empty, error)
	ListModules(ctx context.Context, req *terrariumpb.ListModulesRequest) (*terrariumpb.ListModulesResponse, error)
	ListResources(ctx context.Context, req *terrariumpb.ListResourcesRequest) (*terrariumpb.ListResourcesResponse, error)
	GetModuleDependencies(ctx context.Context, req *terrariumpb.DependencyRequest) (*terrariumpb.DependencyResponse, error)
	GetResourceDependencies(ctx context.Context, req *terrariumpb.DependencyRequest) (*terrariumpb.DependencyResponse, error)
	GetModuleConsumers(ctx context.Context, req *terrariumpb.ConsumerRequest) (*terrariumpb.ConsumerResponse, error)
	GetResourceConsumers(ctx context.Context, req *terrariumpb.ConsumerRequest) (*terrariumpb.ConsumerResponse, error)
	CodeCompletion(ctx context.Context, req *terrariumpb.CompletionRequest) (*terrariumpb.CompletionResponse, error)
}

func New() Service {
	return &service{}
}

type service struct {
	terrariumpb.UnimplementedTerrariumServiceServer
}

func (s service) HealthCheck(ctx context.Context, req *emptypb.Empty) (*emptypb.Empty, error) {
	time.Sleep(200 * time.Millisecond)
	return &emptypb.Empty{}, nil
}
