// Copyright (c) Ollion
// SPDX-License-Identifier: Apache-2.0

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

func (t terrariumAPIImplementor) ListModuleAttributes(ctx context.Context, req *terrariumpb.ListModuleAttributesRequest) (*terrariumpb.ListModuleAttributesResponse, error) {
	return transporthelper.DefaultAPI(ctx, req, t.service.ListModuleAttributes, t.defaultMiddlewareOpts...)
}

func (t terrariumAPIImplementor) ListTaxonomy(ctx context.Context, req *terrariumpb.ListTaxonomyRequest) (*terrariumpb.ListTaxonomyResponse, error) {
	return transporthelper.DefaultAPI(ctx, req, t.service.ListTaxonomy, t.defaultMiddlewareOpts...)
}

func (t terrariumAPIImplementor) ListPlatforms(ctx context.Context, req *terrariumpb.ListPlatformsRequest) (*terrariumpb.ListPlatformsResponse, error) {
	return transporthelper.DefaultAPI(ctx, req, t.service.ListPlatforms, t.defaultMiddlewareOpts...)
}

func (t terrariumAPIImplementor) ListComponents(ctx context.Context, req *terrariumpb.ListComponentsRequest) (*terrariumpb.ListComponentsResponse, error) {
	return transporthelper.DefaultAPI(ctx, req, t.service.ListComponents, t.defaultMiddlewareOpts...)
}

func (t terrariumAPIImplementor) ListDependencies(ctx context.Context, req *terrariumpb.ListDependenciesRequest) (*terrariumpb.ListDependenciesResponse, error) {
	return transporthelper.DefaultAPI(ctx, req, t.service.ListDependencies, t.defaultMiddlewareOpts...)
}
