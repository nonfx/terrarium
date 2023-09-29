// Code generated by protoc-gen-go-grpc. DO NOT EDIT.

package terrariumpb

import (
	context "context"
	grpc "google.golang.org/grpc"
	codes "google.golang.org/grpc/codes"
	status "google.golang.org/grpc/status"
	emptypb "google.golang.org/protobuf/types/known/emptypb"
)

// This is a compile-time assertion to ensure that this generated file
// is compatible with the grpc package it is being compiled against.
const _ = grpc.SupportPackageIsVersion7

// TerrariumServiceClient is the client API for TerrariumService service.
//
// For semantics around ctx use and closing/ending streaming RPCs, please refer to https://pkg.go.dev/google.golang.org/grpc/?tab=doc#ClientConn.NewStream.
type TerrariumServiceClient interface {
	// HealthCheck check endpoint
	HealthCheck(ctx context.Context, in *emptypb.Empty, opts ...grpc.CallOption) (*emptypb.Empty, error)
	// ListModules modules matching the source pattern
	ListModules(ctx context.Context, in *ListModulesRequest, opts ...grpc.CallOption) (*ListModulesResponse, error)
	// ListModuleAttributes returns a list of attributes of the given module.
	// Optionally, it can also include output suggestions that is attributes from other modules that can fullfil this module.
	ListModuleAttributes(ctx context.Context, in *ListModuleAttributesRequest, opts ...grpc.CallOption) (*ListModuleAttributesResponse, error)
}

type terrariumServiceClient struct {
	cc grpc.ClientConnInterface
}

func NewTerrariumServiceClient(cc grpc.ClientConnInterface) TerrariumServiceClient {
	return &terrariumServiceClient{cc}
}

func (c *terrariumServiceClient) HealthCheck(ctx context.Context, in *emptypb.Empty, opts ...grpc.CallOption) (*emptypb.Empty, error) {
	out := new(emptypb.Empty)
	err := c.cc.Invoke(ctx, "/terrarium.v0.TerrariumService/HealthCheck", in, out, opts...)
	if err != nil {
		return nil, err
	}
	return out, nil
}

func (c *terrariumServiceClient) ListModules(ctx context.Context, in *ListModulesRequest, opts ...grpc.CallOption) (*ListModulesResponse, error) {
	out := new(ListModulesResponse)
	err := c.cc.Invoke(ctx, "/terrarium.v0.TerrariumService/ListModules", in, out, opts...)
	if err != nil {
		return nil, err
	}
	return out, nil
}

func (c *terrariumServiceClient) ListModuleAttributes(ctx context.Context, in *ListModuleAttributesRequest, opts ...grpc.CallOption) (*ListModuleAttributesResponse, error) {
	out := new(ListModuleAttributesResponse)
	err := c.cc.Invoke(ctx, "/terrarium.v0.TerrariumService/ListModuleAttributes", in, out, opts...)
	if err != nil {
		return nil, err
	}
	return out, nil
}

// TerrariumServiceServer is the server API for TerrariumService service.
// All implementations must embed UnimplementedTerrariumServiceServer
// for forward compatibility
type TerrariumServiceServer interface {
	// HealthCheck check endpoint
	HealthCheck(context.Context, *emptypb.Empty) (*emptypb.Empty, error)
	// ListModules modules matching the source pattern
	ListModules(context.Context, *ListModulesRequest) (*ListModulesResponse, error)
	// ListModuleAttributes returns a list of attributes of the given module.
	// Optionally, it can also include output suggestions that is attributes from other modules that can fullfil this module.
	ListModuleAttributes(context.Context, *ListModuleAttributesRequest) (*ListModuleAttributesResponse, error)
	mustEmbedUnimplementedTerrariumServiceServer()
}

// UnimplementedTerrariumServiceServer must be embedded to have forward compatible implementations.
type UnimplementedTerrariumServiceServer struct {
}

func (UnimplementedTerrariumServiceServer) HealthCheck(context.Context, *emptypb.Empty) (*emptypb.Empty, error) {
	return nil, status.Errorf(codes.Unimplemented, "method HealthCheck not implemented")
}
func (UnimplementedTerrariumServiceServer) ListModules(context.Context, *ListModulesRequest) (*ListModulesResponse, error) {
	return nil, status.Errorf(codes.Unimplemented, "method ListModules not implemented")
}
func (UnimplementedTerrariumServiceServer) ListModuleAttributes(context.Context, *ListModuleAttributesRequest) (*ListModuleAttributesResponse, error) {
	return nil, status.Errorf(codes.Unimplemented, "method ListModuleAttributes not implemented")
}
func (UnimplementedTerrariumServiceServer) mustEmbedUnimplementedTerrariumServiceServer() {}

// UnsafeTerrariumServiceServer may be embedded to opt out of forward compatibility for this service.
// Use of this interface is not recommended, as added methods to TerrariumServiceServer will
// result in compilation errors.
type UnsafeTerrariumServiceServer interface {
	mustEmbedUnimplementedTerrariumServiceServer()
}

func RegisterTerrariumServiceServer(s grpc.ServiceRegistrar, srv TerrariumServiceServer) {
	s.RegisterService(&_TerrariumService_serviceDesc, srv)
}

func _TerrariumService_HealthCheck_Handler(srv interface{}, ctx context.Context, dec func(interface{}) error, interceptor grpc.UnaryServerInterceptor) (interface{}, error) {
	in := new(emptypb.Empty)
	if err := dec(in); err != nil {
		return nil, err
	}
	if interceptor == nil {
		return srv.(TerrariumServiceServer).HealthCheck(ctx, in)
	}
	info := &grpc.UnaryServerInfo{
		Server:     srv,
		FullMethod: "/terrarium.v0.TerrariumService/HealthCheck",
	}
	handler := func(ctx context.Context, req interface{}) (interface{}, error) {
		return srv.(TerrariumServiceServer).HealthCheck(ctx, req.(*emptypb.Empty))
	}
	return interceptor(ctx, in, info, handler)
}

func _TerrariumService_ListModules_Handler(srv interface{}, ctx context.Context, dec func(interface{}) error, interceptor grpc.UnaryServerInterceptor) (interface{}, error) {
	in := new(ListModulesRequest)
	if err := dec(in); err != nil {
		return nil, err
	}
	if interceptor == nil {
		return srv.(TerrariumServiceServer).ListModules(ctx, in)
	}
	info := &grpc.UnaryServerInfo{
		Server:     srv,
		FullMethod: "/terrarium.v0.TerrariumService/ListModules",
	}
	handler := func(ctx context.Context, req interface{}) (interface{}, error) {
		return srv.(TerrariumServiceServer).ListModules(ctx, req.(*ListModulesRequest))
	}
	return interceptor(ctx, in, info, handler)
}

func _TerrariumService_ListModuleAttributes_Handler(srv interface{}, ctx context.Context, dec func(interface{}) error, interceptor grpc.UnaryServerInterceptor) (interface{}, error) {
	in := new(ListModuleAttributesRequest)
	if err := dec(in); err != nil {
		return nil, err
	}
	if interceptor == nil {
		return srv.(TerrariumServiceServer).ListModuleAttributes(ctx, in)
	}
	info := &grpc.UnaryServerInfo{
		Server:     srv,
		FullMethod: "/terrarium.v0.TerrariumService/ListModuleAttributes",
	}
	handler := func(ctx context.Context, req interface{}) (interface{}, error) {
		return srv.(TerrariumServiceServer).ListModuleAttributes(ctx, req.(*ListModuleAttributesRequest))
	}
	return interceptor(ctx, in, info, handler)
}

var _TerrariumService_serviceDesc = grpc.ServiceDesc{
	ServiceName: "terrarium.v0.TerrariumService",
	HandlerType: (*TerrariumServiceServer)(nil),
	Methods: []grpc.MethodDesc{
		{
			MethodName: "HealthCheck",
			Handler:    _TerrariumService_HealthCheck_Handler,
		},
		{
			MethodName: "ListModules",
			Handler:    _TerrariumService_ListModules_Handler,
		},
		{
			MethodName: "ListModuleAttributes",
			Handler:    _TerrariumService_ListModuleAttributes_Handler,
		},
	},
	Streams:  []grpc.StreamDesc{},
	Metadata: "terrariumpb/service.proto",
}
