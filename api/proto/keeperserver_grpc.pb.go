// Code generated by protoc-gen-go-grpc. DO NOT EDIT.
// versions:
// - protoc-gen-go-grpc v1.3.0
// - protoc             v3.21.12
// source: api/proto/keeperserver.proto

package proto

import (
	context "context"
	grpc "google.golang.org/grpc"
	codes "google.golang.org/grpc/codes"
	status "google.golang.org/grpc/status"
)

// This is a compile-time assertion to ensure that this generated file
// is compatible with the grpc package it is being compiled against.
// Requires gRPC-Go v1.32.0 or later.
const _ = grpc.SupportPackageIsVersion7

const (
	KeeperService_Ping_FullMethodName = "/keeperservice.grpc.KeeperService/Ping"
)

// KeeperServiceClient is the client API for KeeperService service.
//
// For semantics around ctx use and closing/ending streaming RPCs, please refer to https://pkg.go.dev/google.golang.org/grpc/?tab=doc#ClientConn.NewStream.
type KeeperServiceClient interface {
	Ping(ctx context.Context, in *PingRequest, opts ...grpc.CallOption) (*PingResponse, error)
}

type keeperServiceClient struct {
	cc grpc.ClientConnInterface
}

func NewKeeperServiceClient(cc grpc.ClientConnInterface) KeeperServiceClient {
	return &keeperServiceClient{cc}
}

func (c *keeperServiceClient) Ping(ctx context.Context, in *PingRequest, opts ...grpc.CallOption) (*PingResponse, error) {
	out := new(PingResponse)
	err := c.cc.Invoke(ctx, KeeperService_Ping_FullMethodName, in, out, opts...)
	if err != nil {
		return nil, err
	}
	return out, nil
}

// KeeperServiceServer is the server API for KeeperService service.
// All implementations must embed UnimplementedKeeperServiceServer
// for forward compatibility
type KeeperServiceServer interface {
	Ping(context.Context, *PingRequest) (*PingResponse, error)
	mustEmbedUnimplementedKeeperServiceServer()
}

// UnimplementedKeeperServiceServer must be embedded to have forward compatible implementations.
type UnimplementedKeeperServiceServer struct {
}

func (UnimplementedKeeperServiceServer) Ping(context.Context, *PingRequest) (*PingResponse, error) {
	return nil, status.Errorf(codes.Unimplemented, "method Ping not implemented")
}
func (UnimplementedKeeperServiceServer) mustEmbedUnimplementedKeeperServiceServer() {}

// UnsafeKeeperServiceServer may be embedded to opt out of forward compatibility for this service.
// Use of this interface is not recommended, as added methods to KeeperServiceServer will
// result in compilation errors.
type UnsafeKeeperServiceServer interface {
	mustEmbedUnimplementedKeeperServiceServer()
}

func RegisterKeeperServiceServer(s grpc.ServiceRegistrar, srv KeeperServiceServer) {
	s.RegisterService(&KeeperService_ServiceDesc, srv)
}

func _KeeperService_Ping_Handler(srv interface{}, ctx context.Context, dec func(interface{}) error, interceptor grpc.UnaryServerInterceptor) (interface{}, error) {
	in := new(PingRequest)
	if err := dec(in); err != nil {
		return nil, err
	}
	if interceptor == nil {
		return srv.(KeeperServiceServer).Ping(ctx, in)
	}
	info := &grpc.UnaryServerInfo{
		Server:     srv,
		FullMethod: KeeperService_Ping_FullMethodName,
	}
	handler := func(ctx context.Context, req interface{}) (interface{}, error) {
		return srv.(KeeperServiceServer).Ping(ctx, req.(*PingRequest))
	}
	return interceptor(ctx, in, info, handler)
}

// KeeperService_ServiceDesc is the grpc.ServiceDesc for KeeperService service.
// It's only intended for direct use with grpc.RegisterService,
// and not to be introspected or modified (even as a copy)
var KeeperService_ServiceDesc = grpc.ServiceDesc{
	ServiceName: "keeperservice.grpc.KeeperService",
	HandlerType: (*KeeperServiceServer)(nil),
	Methods: []grpc.MethodDesc{
		{
			MethodName: "Ping",
			Handler:    _KeeperService_Ping_Handler,
		},
	},
	Streams:  []grpc.StreamDesc{},
	Metadata: "api/proto/keeperserver.proto",
}
