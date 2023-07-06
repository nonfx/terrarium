package transport

import (
	"context"

	"github.com/cldcvr/terrarium/src/api/service"
	"github.com/cldcvr/terrarium/src/pkg/pb/terrariumpb"
	"github.com/cldcvr/terrarium/src/pkg/transporthelper"
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

func (t terrariumAPIImplementor) CodeCompletion(ctx context.Context, req *terrariumpb.CompletionRequest) (*terrariumpb.CompletionResponse, error) {
	return transporthelper.DefaultAPI(ctx, req, t.service.CodeCompletion, t.defaultMiddlewareOpts...)
}

func (t terrariumAPIImplementor) ListModuleAttributes(ctx context.Context, req *terrariumpb.ListModuleAttributesRequest) (*terrariumpb.ListModuleAttributesResponse, error) {
	return transporthelper.DefaultAPI(ctx, req, t.service.ListModuleAttributes, t.defaultMiddlewareOpts...)
}
