package transport

import (
	"context"

	"github.com/cldcvr/terrarium/api/pkg/pb/terrariumpb"
	"github.com/cldcvr/terrarium/api/pkg/transporthelper"
	"github.com/cldcvr/terrarium/api/service"
	"github.com/go-kit/kit/endpoint"
	"google.golang.org/protobuf/types/known/emptypb"
)

type terrariumAPIImplementor struct {
	terrariumpb.UnimplementedTerrariumServiceServer

	service               service.Service
	defaultMiddlewareOpts []endpoint.Middleware
}

func NewTerrariumAPI(service service.Service) terrariumpb.TerrariumServiceServer {
	return &terrariumAPIImplementor{service: service, defaultMiddlewareOpts: []endpoint.Middleware{
		transporthelper.WithReqValidatorEPMiddleware(),
		transporthelper.WithLoggingEPMiddleware(),
	}}
}

func (t terrariumAPIImplementor) HealthCheck(ctx context.Context, req *emptypb.Empty) (*emptypb.Empty, error) {
	return transporthelper.DefaultAPI(ctx, req, t.service.HealthCheck, t.defaultMiddlewareOpts...)
}

func (t terrariumAPIImplementor) ListModules(ctx context.Context, req *terrariumpb.ListModulesRequest) (*terrariumpb.ListModulesResponse, error) {
	return transporthelper.DefaultAPI(ctx, req, t.service.ListModules, t.defaultMiddlewareOpts...)
}

func (t terrariumAPIImplementor) ListResources(ctx context.Context, req *terrariumpb.ListResourcesRequest) (*terrariumpb.ListResourcesResponse, error) {
	return transporthelper.DefaultAPI(ctx, req, t.service.ListResources, t.defaultMiddlewareOpts...)
}

func (t terrariumAPIImplementor) GetModuleDependencies(ctx context.Context, req *terrariumpb.DependencyRequest) (*terrariumpb.DependencyResponse, error) {
	return transporthelper.DefaultAPI(ctx, req, t.service.GetModuleDependencies, t.defaultMiddlewareOpts...)
}

func (t terrariumAPIImplementor) GetResourceDependencies(ctx context.Context, req *terrariumpb.DependencyRequest) (*terrariumpb.DependencyResponse, error) {
	return transporthelper.DefaultAPI(ctx, req, t.service.GetResourceDependencies, t.defaultMiddlewareOpts...)
}

func (t terrariumAPIImplementor) GetModuleConsumers(ctx context.Context, req *terrariumpb.ConsumerRequest) (*terrariumpb.ConsumerResponse, error) {
	return transporthelper.DefaultAPI(ctx, req, t.service.GetModuleConsumers, t.defaultMiddlewareOpts...)
}

func (t terrariumAPIImplementor) GetResourceConsumers(ctx context.Context, req *terrariumpb.ConsumerRequest) (*terrariumpb.ConsumerResponse, error) {
	return transporthelper.DefaultAPI(ctx, req, t.service.GetResourceConsumers, t.defaultMiddlewareOpts...)
}

func (t terrariumAPIImplementor) CodeCompletion(ctx context.Context, req *terrariumpb.CompletionRequest) (*terrariumpb.CompletionResponse, error) {
	return transporthelper.DefaultAPI(ctx, req, t.service.CodeCompletion, t.defaultMiddlewareOpts...)
}