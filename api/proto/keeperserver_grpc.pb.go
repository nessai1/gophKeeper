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
	KeeperService_Ping_FullMethodName                = "/keeperservice.grpc.KeeperService/Ping"
	KeeperService_Register_FullMethodName            = "/keeperservice.grpc.KeeperService/Register"
	KeeperService_Login_FullMethodName               = "/keeperservice.grpc.KeeperService/Login"
	KeeperService_UploadMediaSecret_FullMethodName   = "/keeperservice.grpc.KeeperService/UploadMediaSecret"
	KeeperService_DownloadMediaSecret_FullMethodName = "/keeperservice.grpc.KeeperService/DownloadMediaSecret"
	KeeperService_SecretList_FullMethodName          = "/keeperservice.grpc.KeeperService/SecretList"
	KeeperService_SecretSet_FullMethodName           = "/keeperservice.grpc.KeeperService/SecretSet"
	KeeperService_SecretGet_FullMethodName           = "/keeperservice.grpc.KeeperService/SecretGet"
	KeeperService_SecretUpdate_FullMethodName        = "/keeperservice.grpc.KeeperService/SecretUpdate"
	KeeperService_SecretDelete_FullMethodName        = "/keeperservice.grpc.KeeperService/SecretDelete"
)

// KeeperServiceClient is the client API for KeeperService service.
//
// For semantics around ctx use and closing/ending streaming RPCs, please refer to https://pkg.go.dev/google.golang.org/grpc/?tab=doc#ClientConn.NewStream.
type KeeperServiceClient interface {
	Ping(ctx context.Context, in *PingRequest, opts ...grpc.CallOption) (*PingResponse, error)
	Register(ctx context.Context, in *UserCredentialsRequest, opts ...grpc.CallOption) (*UserCredentialsResponse, error)
	Login(ctx context.Context, in *UserCredentialsRequest, opts ...grpc.CallOption) (*UserCredentialsResponse, error)
	UploadMediaSecret(ctx context.Context, opts ...grpc.CallOption) (KeeperService_UploadMediaSecretClient, error)
	DownloadMediaSecret(ctx context.Context, in *DownloadMediaSecretRequest, opts ...grpc.CallOption) (KeeperService_DownloadMediaSecretClient, error)
	SecretList(ctx context.Context, in *SecretListRequest, opts ...grpc.CallOption) (*SecretListResponse, error)
	SecretSet(ctx context.Context, in *SecretSetRequest, opts ...grpc.CallOption) (*SecretSetResponse, error)
	SecretGet(ctx context.Context, in *SecretGetRequest, opts ...grpc.CallOption) (*SecretGetResponse, error)
	SecretUpdate(ctx context.Context, in *SecretUpdateRequest, opts ...grpc.CallOption) (*SecretUpdateResponse, error)
	SecretDelete(ctx context.Context, in *SecretDeleteRequest, opts ...grpc.CallOption) (*SecretDeleteResponse, error)
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

func (c *keeperServiceClient) Register(ctx context.Context, in *UserCredentialsRequest, opts ...grpc.CallOption) (*UserCredentialsResponse, error) {
	out := new(UserCredentialsResponse)
	err := c.cc.Invoke(ctx, KeeperService_Register_FullMethodName, in, out, opts...)
	if err != nil {
		return nil, err
	}
	return out, nil
}

func (c *keeperServiceClient) Login(ctx context.Context, in *UserCredentialsRequest, opts ...grpc.CallOption) (*UserCredentialsResponse, error) {
	out := new(UserCredentialsResponse)
	err := c.cc.Invoke(ctx, KeeperService_Login_FullMethodName, in, out, opts...)
	if err != nil {
		return nil, err
	}
	return out, nil
}

func (c *keeperServiceClient) UploadMediaSecret(ctx context.Context, opts ...grpc.CallOption) (KeeperService_UploadMediaSecretClient, error) {
	stream, err := c.cc.NewStream(ctx, &KeeperService_ServiceDesc.Streams[0], KeeperService_UploadMediaSecret_FullMethodName, opts...)
	if err != nil {
		return nil, err
	}
	x := &keeperServiceUploadMediaSecretClient{stream}
	return x, nil
}

type KeeperService_UploadMediaSecretClient interface {
	Send(*UploadMediaSecretRequest) error
	CloseAndRecv() (*UploadMediaSecretResponse, error)
	grpc.ClientStream
}

type keeperServiceUploadMediaSecretClient struct {
	grpc.ClientStream
}

func (x *keeperServiceUploadMediaSecretClient) Send(m *UploadMediaSecretRequest) error {
	return x.ClientStream.SendMsg(m)
}

func (x *keeperServiceUploadMediaSecretClient) CloseAndRecv() (*UploadMediaSecretResponse, error) {
	if err := x.ClientStream.CloseSend(); err != nil {
		return nil, err
	}
	m := new(UploadMediaSecretResponse)
	if err := x.ClientStream.RecvMsg(m); err != nil {
		return nil, err
	}
	return m, nil
}

func (c *keeperServiceClient) DownloadMediaSecret(ctx context.Context, in *DownloadMediaSecretRequest, opts ...grpc.CallOption) (KeeperService_DownloadMediaSecretClient, error) {
	stream, err := c.cc.NewStream(ctx, &KeeperService_ServiceDesc.Streams[1], KeeperService_DownloadMediaSecret_FullMethodName, opts...)
	if err != nil {
		return nil, err
	}
	x := &keeperServiceDownloadMediaSecretClient{stream}
	if err := x.ClientStream.SendMsg(in); err != nil {
		return nil, err
	}
	if err := x.ClientStream.CloseSend(); err != nil {
		return nil, err
	}
	return x, nil
}

type KeeperService_DownloadMediaSecretClient interface {
	Recv() (*DownloadMediaSecretResponse, error)
	grpc.ClientStream
}

type keeperServiceDownloadMediaSecretClient struct {
	grpc.ClientStream
}

func (x *keeperServiceDownloadMediaSecretClient) Recv() (*DownloadMediaSecretResponse, error) {
	m := new(DownloadMediaSecretResponse)
	if err := x.ClientStream.RecvMsg(m); err != nil {
		return nil, err
	}
	return m, nil
}

func (c *keeperServiceClient) SecretList(ctx context.Context, in *SecretListRequest, opts ...grpc.CallOption) (*SecretListResponse, error) {
	out := new(SecretListResponse)
	err := c.cc.Invoke(ctx, KeeperService_SecretList_FullMethodName, in, out, opts...)
	if err != nil {
		return nil, err
	}
	return out, nil
}

func (c *keeperServiceClient) SecretSet(ctx context.Context, in *SecretSetRequest, opts ...grpc.CallOption) (*SecretSetResponse, error) {
	out := new(SecretSetResponse)
	err := c.cc.Invoke(ctx, KeeperService_SecretSet_FullMethodName, in, out, opts...)
	if err != nil {
		return nil, err
	}
	return out, nil
}

func (c *keeperServiceClient) SecretGet(ctx context.Context, in *SecretGetRequest, opts ...grpc.CallOption) (*SecretGetResponse, error) {
	out := new(SecretGetResponse)
	err := c.cc.Invoke(ctx, KeeperService_SecretGet_FullMethodName, in, out, opts...)
	if err != nil {
		return nil, err
	}
	return out, nil
}

func (c *keeperServiceClient) SecretUpdate(ctx context.Context, in *SecretUpdateRequest, opts ...grpc.CallOption) (*SecretUpdateResponse, error) {
	out := new(SecretUpdateResponse)
	err := c.cc.Invoke(ctx, KeeperService_SecretUpdate_FullMethodName, in, out, opts...)
	if err != nil {
		return nil, err
	}
	return out, nil
}

func (c *keeperServiceClient) SecretDelete(ctx context.Context, in *SecretDeleteRequest, opts ...grpc.CallOption) (*SecretDeleteResponse, error) {
	out := new(SecretDeleteResponse)
	err := c.cc.Invoke(ctx, KeeperService_SecretDelete_FullMethodName, in, out, opts...)
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
	Register(context.Context, *UserCredentialsRequest) (*UserCredentialsResponse, error)
	Login(context.Context, *UserCredentialsRequest) (*UserCredentialsResponse, error)
	UploadMediaSecret(KeeperService_UploadMediaSecretServer) error
	DownloadMediaSecret(*DownloadMediaSecretRequest, KeeperService_DownloadMediaSecretServer) error
	SecretList(context.Context, *SecretListRequest) (*SecretListResponse, error)
	SecretSet(context.Context, *SecretSetRequest) (*SecretSetResponse, error)
	SecretGet(context.Context, *SecretGetRequest) (*SecretGetResponse, error)
	SecretUpdate(context.Context, *SecretUpdateRequest) (*SecretUpdateResponse, error)
	SecretDelete(context.Context, *SecretDeleteRequest) (*SecretDeleteResponse, error)
	mustEmbedUnimplementedKeeperServiceServer()
}

// UnimplementedKeeperServiceServer must be embedded to have forward compatible implementations.
type UnimplementedKeeperServiceServer struct {
}

func (UnimplementedKeeperServiceServer) Ping(context.Context, *PingRequest) (*PingResponse, error) {
	return nil, status.Errorf(codes.Unimplemented, "method Ping not implemented")
}
func (UnimplementedKeeperServiceServer) Register(context.Context, *UserCredentialsRequest) (*UserCredentialsResponse, error) {
	return nil, status.Errorf(codes.Unimplemented, "method Register not implemented")
}
func (UnimplementedKeeperServiceServer) Login(context.Context, *UserCredentialsRequest) (*UserCredentialsResponse, error) {
	return nil, status.Errorf(codes.Unimplemented, "method Login not implemented")
}
func (UnimplementedKeeperServiceServer) UploadMediaSecret(KeeperService_UploadMediaSecretServer) error {
	return status.Errorf(codes.Unimplemented, "method UploadMediaSecret not implemented")
}
func (UnimplementedKeeperServiceServer) DownloadMediaSecret(*DownloadMediaSecretRequest, KeeperService_DownloadMediaSecretServer) error {
	return status.Errorf(codes.Unimplemented, "method DownloadMediaSecret not implemented")
}
func (UnimplementedKeeperServiceServer) SecretList(context.Context, *SecretListRequest) (*SecretListResponse, error) {
	return nil, status.Errorf(codes.Unimplemented, "method SecretList not implemented")
}
func (UnimplementedKeeperServiceServer) SecretSet(context.Context, *SecretSetRequest) (*SecretSetResponse, error) {
	return nil, status.Errorf(codes.Unimplemented, "method SecretSet not implemented")
}
func (UnimplementedKeeperServiceServer) SecretGet(context.Context, *SecretGetRequest) (*SecretGetResponse, error) {
	return nil, status.Errorf(codes.Unimplemented, "method SecretGet not implemented")
}
func (UnimplementedKeeperServiceServer) SecretUpdate(context.Context, *SecretUpdateRequest) (*SecretUpdateResponse, error) {
	return nil, status.Errorf(codes.Unimplemented, "method SecretUpdate not implemented")
}
func (UnimplementedKeeperServiceServer) SecretDelete(context.Context, *SecretDeleteRequest) (*SecretDeleteResponse, error) {
	return nil, status.Errorf(codes.Unimplemented, "method SecretDelete not implemented")
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

func _KeeperService_Register_Handler(srv interface{}, ctx context.Context, dec func(interface{}) error, interceptor grpc.UnaryServerInterceptor) (interface{}, error) {
	in := new(UserCredentialsRequest)
	if err := dec(in); err != nil {
		return nil, err
	}
	if interceptor == nil {
		return srv.(KeeperServiceServer).Register(ctx, in)
	}
	info := &grpc.UnaryServerInfo{
		Server:     srv,
		FullMethod: KeeperService_Register_FullMethodName,
	}
	handler := func(ctx context.Context, req interface{}) (interface{}, error) {
		return srv.(KeeperServiceServer).Register(ctx, req.(*UserCredentialsRequest))
	}
	return interceptor(ctx, in, info, handler)
}

func _KeeperService_Login_Handler(srv interface{}, ctx context.Context, dec func(interface{}) error, interceptor grpc.UnaryServerInterceptor) (interface{}, error) {
	in := new(UserCredentialsRequest)
	if err := dec(in); err != nil {
		return nil, err
	}
	if interceptor == nil {
		return srv.(KeeperServiceServer).Login(ctx, in)
	}
	info := &grpc.UnaryServerInfo{
		Server:     srv,
		FullMethod: KeeperService_Login_FullMethodName,
	}
	handler := func(ctx context.Context, req interface{}) (interface{}, error) {
		return srv.(KeeperServiceServer).Login(ctx, req.(*UserCredentialsRequest))
	}
	return interceptor(ctx, in, info, handler)
}

func _KeeperService_UploadMediaSecret_Handler(srv interface{}, stream grpc.ServerStream) error {
	return srv.(KeeperServiceServer).UploadMediaSecret(&keeperServiceUploadMediaSecretServer{stream})
}

type KeeperService_UploadMediaSecretServer interface {
	SendAndClose(*UploadMediaSecretResponse) error
	Recv() (*UploadMediaSecretRequest, error)
	grpc.ServerStream
}

type keeperServiceUploadMediaSecretServer struct {
	grpc.ServerStream
}

func (x *keeperServiceUploadMediaSecretServer) SendAndClose(m *UploadMediaSecretResponse) error {
	return x.ServerStream.SendMsg(m)
}

func (x *keeperServiceUploadMediaSecretServer) Recv() (*UploadMediaSecretRequest, error) {
	m := new(UploadMediaSecretRequest)
	if err := x.ServerStream.RecvMsg(m); err != nil {
		return nil, err
	}
	return m, nil
}

func _KeeperService_DownloadMediaSecret_Handler(srv interface{}, stream grpc.ServerStream) error {
	m := new(DownloadMediaSecretRequest)
	if err := stream.RecvMsg(m); err != nil {
		return err
	}
	return srv.(KeeperServiceServer).DownloadMediaSecret(m, &keeperServiceDownloadMediaSecretServer{stream})
}

type KeeperService_DownloadMediaSecretServer interface {
	Send(*DownloadMediaSecretResponse) error
	grpc.ServerStream
}

type keeperServiceDownloadMediaSecretServer struct {
	grpc.ServerStream
}

func (x *keeperServiceDownloadMediaSecretServer) Send(m *DownloadMediaSecretResponse) error {
	return x.ServerStream.SendMsg(m)
}

func _KeeperService_SecretList_Handler(srv interface{}, ctx context.Context, dec func(interface{}) error, interceptor grpc.UnaryServerInterceptor) (interface{}, error) {
	in := new(SecretListRequest)
	if err := dec(in); err != nil {
		return nil, err
	}
	if interceptor == nil {
		return srv.(KeeperServiceServer).SecretList(ctx, in)
	}
	info := &grpc.UnaryServerInfo{
		Server:     srv,
		FullMethod: KeeperService_SecretList_FullMethodName,
	}
	handler := func(ctx context.Context, req interface{}) (interface{}, error) {
		return srv.(KeeperServiceServer).SecretList(ctx, req.(*SecretListRequest))
	}
	return interceptor(ctx, in, info, handler)
}

func _KeeperService_SecretSet_Handler(srv interface{}, ctx context.Context, dec func(interface{}) error, interceptor grpc.UnaryServerInterceptor) (interface{}, error) {
	in := new(SecretSetRequest)
	if err := dec(in); err != nil {
		return nil, err
	}
	if interceptor == nil {
		return srv.(KeeperServiceServer).SecretSet(ctx, in)
	}
	info := &grpc.UnaryServerInfo{
		Server:     srv,
		FullMethod: KeeperService_SecretSet_FullMethodName,
	}
	handler := func(ctx context.Context, req interface{}) (interface{}, error) {
		return srv.(KeeperServiceServer).SecretSet(ctx, req.(*SecretSetRequest))
	}
	return interceptor(ctx, in, info, handler)
}

func _KeeperService_SecretGet_Handler(srv interface{}, ctx context.Context, dec func(interface{}) error, interceptor grpc.UnaryServerInterceptor) (interface{}, error) {
	in := new(SecretGetRequest)
	if err := dec(in); err != nil {
		return nil, err
	}
	if interceptor == nil {
		return srv.(KeeperServiceServer).SecretGet(ctx, in)
	}
	info := &grpc.UnaryServerInfo{
		Server:     srv,
		FullMethod: KeeperService_SecretGet_FullMethodName,
	}
	handler := func(ctx context.Context, req interface{}) (interface{}, error) {
		return srv.(KeeperServiceServer).SecretGet(ctx, req.(*SecretGetRequest))
	}
	return interceptor(ctx, in, info, handler)
}

func _KeeperService_SecretUpdate_Handler(srv interface{}, ctx context.Context, dec func(interface{}) error, interceptor grpc.UnaryServerInterceptor) (interface{}, error) {
	in := new(SecretUpdateRequest)
	if err := dec(in); err != nil {
		return nil, err
	}
	if interceptor == nil {
		return srv.(KeeperServiceServer).SecretUpdate(ctx, in)
	}
	info := &grpc.UnaryServerInfo{
		Server:     srv,
		FullMethod: KeeperService_SecretUpdate_FullMethodName,
	}
	handler := func(ctx context.Context, req interface{}) (interface{}, error) {
		return srv.(KeeperServiceServer).SecretUpdate(ctx, req.(*SecretUpdateRequest))
	}
	return interceptor(ctx, in, info, handler)
}

func _KeeperService_SecretDelete_Handler(srv interface{}, ctx context.Context, dec func(interface{}) error, interceptor grpc.UnaryServerInterceptor) (interface{}, error) {
	in := new(SecretDeleteRequest)
	if err := dec(in); err != nil {
		return nil, err
	}
	if interceptor == nil {
		return srv.(KeeperServiceServer).SecretDelete(ctx, in)
	}
	info := &grpc.UnaryServerInfo{
		Server:     srv,
		FullMethod: KeeperService_SecretDelete_FullMethodName,
	}
	handler := func(ctx context.Context, req interface{}) (interface{}, error) {
		return srv.(KeeperServiceServer).SecretDelete(ctx, req.(*SecretDeleteRequest))
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
		{
			MethodName: "Register",
			Handler:    _KeeperService_Register_Handler,
		},
		{
			MethodName: "Login",
			Handler:    _KeeperService_Login_Handler,
		},
		{
			MethodName: "SecretList",
			Handler:    _KeeperService_SecretList_Handler,
		},
		{
			MethodName: "SecretSet",
			Handler:    _KeeperService_SecretSet_Handler,
		},
		{
			MethodName: "SecretGet",
			Handler:    _KeeperService_SecretGet_Handler,
		},
		{
			MethodName: "SecretUpdate",
			Handler:    _KeeperService_SecretUpdate_Handler,
		},
		{
			MethodName: "SecretDelete",
			Handler:    _KeeperService_SecretDelete_Handler,
		},
	},
	Streams: []grpc.StreamDesc{
		{
			StreamName:    "UploadMediaSecret",
			Handler:       _KeeperService_UploadMediaSecret_Handler,
			ClientStreams: true,
		},
		{
			StreamName:    "DownloadMediaSecret",
			Handler:       _KeeperService_DownloadMediaSecret_Handler,
			ServerStreams: true,
		},
	},
	Metadata: "api/proto/keeperserver.proto",
}
