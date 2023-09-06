// Copyright (c) CloudCover
// SPDX-License-Identifier: Apache-2.0

package transporthelper

import (
	"context"
	"time"

	"github.com/go-kit/kit/endpoint"
	"github.com/grpc-ecosystem/grpc-gateway/v2/runtime"
	"github.com/rotisserie/eris"
	log "github.com/sirupsen/logrus"
	"google.golang.org/grpc/codes"
	"google.golang.org/grpc/metadata"
	"google.golang.org/grpc/status"
)

// Errors
var (
	ErrReqInvalidType = status.Error(codes.Internal, "invalid request type")
	ErrReqInvalid     = status.Error(codes.InvalidArgument, "request validation failed")
)

// ServiceFunc represents the function signature for a service.
type ServiceFunc[REQ, RES any] func(context.Context, *REQ) (*RES, error)

// DefaultEPCall validates the request and calls the given handler.
func DefaultEPCall[RES any](ctx context.Context, ep endpoint.Endpoint, req interface{}) (*RES, error) {
	res, err := ep(ctx, req)
	if err != nil {
		return nil, err
	}
	return res.(*RES), nil
}

// ValidateReq runs validation on the proto request object.
func ValidateReq(req interface{}) error {
	validator, hasValidator := req.(interface{ Validate() error })
	if !hasValidator {
		// No validator, skip validation.
		return nil
	}
	err := validator.Validate()
	if err != nil {
		return eris.Wrapf(ErrReqInvalid, "%v", err)
	}
	return nil
}

// DefaultEP returns an endpoint that calls the service function.
func DefaultEP[REQ, RES any](ctx context.Context, serviceFunc ServiceFunc[REQ, RES]) endpoint.Endpoint {
	return func(ctx context.Context, request interface{}) (response interface{}, err error) {
		reqTyped, ok := request.(*REQ)
		if !ok {
			var t REQ
			return nil, eris.Wrapf(ErrReqInvalidType, "expected %T got %T", &t, request)
		}
		return serviceFunc(ctx, reqTyped)
	}
}

// WithLoggingEPMiddleware returns a middleware that logs API calls.
func WithLoggingEPMiddleware() endpoint.Middleware {
	return func(next endpoint.Endpoint) endpoint.Endpoint {
		return func(ctx context.Context, request interface{}) (response interface{}, err error) {
			logger := log.WithFields(log.Fields{
				MdKeyHTTPRequestID: GetReqIdFromCtx(ctx),
			})
			logger = logWithReqMeta(ctx, logger)
			defer func(begin time.Time) {
				logger = logger.WithFields(log.Fields{
					"took": time.Since(begin).String(),
				})
				if err != nil {
					logger = logWithErrInfo(err, logger)
					logger.Error("API call failed")
				} else {
					logger.Info("API call passed")
				}
			}(time.Now())
			response, err = next(ctx, request)
			return
		}
	}
}

// WithReqValidatorEPMiddleware returns a middleware that validates the request.
func WithReqValidatorEPMiddleware() endpoint.Middleware {
	return func(next endpoint.Endpoint) endpoint.Endpoint {
		return func(ctx context.Context, request interface{}) (response interface{}, err error) {
			err = ValidateReq(request)
			if err != nil {
				return nil, err
			}
			return next(ctx, request)
		}
	}
}

// attachMiddlewares attaches middleware to an endpoint.
func attachMiddlewares(ep endpoint.Endpoint, middlewareOptions ...endpoint.Middleware) endpoint.Endpoint {
	for _, m := range middlewareOptions {
		ep = m(ep)
	}
	return ep
}

// DefaultAPI is a convenience function for handling API requests.
func DefaultAPI[REQ, RES any](ctx context.Context, req *REQ, serviceFunc ServiceFunc[REQ, RES], middlewareOptions ...endpoint.Middleware) (*RES, error) {
	ep := attachMiddlewares(DefaultEP(ctx, serviceFunc), middlewareOptions...)
	return DefaultEPCall[RES](ctx, ep, req)
}

// logWithReqMeta adds request metadata to the logger.
func logWithReqMeta(ctx context.Context, logger *log.Entry) *log.Entry {
	f := log.Fields{}
	// Extract from incoming metadata.
	for _, k := range []string{MdKeyHTTPMethod, MdKeyHTTPPath, MdKeyHTTPHost, MdKeyHTTPRemote} {
		v := metadata.ValueFromIncomingContext(ctx, k)
		if len(v) > 0 {
			f[k] = v[0]
		}
	}
	// Extract from context private key.
	for _, k := range []string{MdKeyGRPCMethod} {
		v := ctx.Value(ctxPrivateKey(k))
		if v != nil {
			f[k] = v
		}
	}
	return logger.WithFields(f)
}

// logWithErrInfo adds error information to the logger.
func logWithErrInfo(err error, logger *log.Entry) *log.Entry {
	grpcCode := status.Convert(err).Code()
	httpStatus := runtime.HTTPStatusFromCode(grpcCode)
	return logger.WithFields(log.Fields{
		"error":       eris.ToJSON(err, true),
		"http_status": httpStatus,
		"grpc_code":   grpcCode,
	})
}
