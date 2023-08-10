// Code generated by protoc-gen-go-grpc. DO NOT EDIT.
// versions:
// - protoc-gen-go-grpc v1.3.0
// - protoc             v4.23.4
// source: service_bank_app.proto

package pb

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
	BankApp_CreateUser_FullMethodName = "/pb.BankApp/CreateUser"
	BankApp_LoginUser_FullMethodName  = "/pb.BankApp/LoginUser"
)

// BankAppClient is the client API for BankApp service.
//
// For semantics around ctx use and closing/ending streaming RPCs, please refer to https://pkg.go.dev/google.golang.org/grpc/?tab=doc#ClientConn.NewStream.
type BankAppClient interface {
	CreateUser(ctx context.Context, in *CreateUserRequest, opts ...grpc.CallOption) (*CreateUserResponse, error)
	LoginUser(ctx context.Context, in *LoginUserRequest, opts ...grpc.CallOption) (*LoginUserResponse, error)
}

type bankAppClient struct {
	cc grpc.ClientConnInterface
}

func NewBankAppClient(cc grpc.ClientConnInterface) BankAppClient {
	return &bankAppClient{cc}
}

func (c *bankAppClient) CreateUser(ctx context.Context, in *CreateUserRequest, opts ...grpc.CallOption) (*CreateUserResponse, error) {
	out := new(CreateUserResponse)
	err := c.cc.Invoke(ctx, BankApp_CreateUser_FullMethodName, in, out, opts...)
	if err != nil {
		return nil, err
	}
	return out, nil
}

func (c *bankAppClient) LoginUser(ctx context.Context, in *LoginUserRequest, opts ...grpc.CallOption) (*LoginUserResponse, error) {
	out := new(LoginUserResponse)
	err := c.cc.Invoke(ctx, BankApp_LoginUser_FullMethodName, in, out, opts...)
	if err != nil {
		return nil, err
	}
	return out, nil
}

// BankAppServer is the server API for BankApp service.
// All implementations must embed UnimplementedBankAppServer
// for forward compatibility
type BankAppServer interface {
	CreateUser(context.Context, *CreateUserRequest) (*CreateUserResponse, error)
	LoginUser(context.Context, *LoginUserRequest) (*LoginUserResponse, error)
	mustEmbedUnimplementedBankAppServer()
}

// UnimplementedBankAppServer must be embedded to have forward compatible implementations.
type UnimplementedBankAppServer struct {
}

func (UnimplementedBankAppServer) CreateUser(context.Context, *CreateUserRequest) (*CreateUserResponse, error) {
	return nil, status.Errorf(codes.Unimplemented, "method CreateUser not implemented")
}
func (UnimplementedBankAppServer) LoginUser(context.Context, *LoginUserRequest) (*LoginUserResponse, error) {
	return nil, status.Errorf(codes.Unimplemented, "method LoginUser not implemented")
}
func (UnimplementedBankAppServer) mustEmbedUnimplementedBankAppServer() {}

// UnsafeBankAppServer may be embedded to opt out of forward compatibility for this service.
// Use of this interface is not recommended, as added methods to BankAppServer will
// result in compilation errors.
type UnsafeBankAppServer interface {
	mustEmbedUnimplementedBankAppServer()
}

func RegisterBankAppServer(s grpc.ServiceRegistrar, srv BankAppServer) {
	s.RegisterService(&BankApp_ServiceDesc, srv)
}

func _BankApp_CreateUser_Handler(srv interface{}, ctx context.Context, dec func(interface{}) error, interceptor grpc.UnaryServerInterceptor) (interface{}, error) {
	in := new(CreateUserRequest)
	if err := dec(in); err != nil {
		return nil, err
	}
	if interceptor == nil {
		return srv.(BankAppServer).CreateUser(ctx, in)
	}
	info := &grpc.UnaryServerInfo{
		Server:     srv,
		FullMethod: BankApp_CreateUser_FullMethodName,
	}
	handler := func(ctx context.Context, req interface{}) (interface{}, error) {
		return srv.(BankAppServer).CreateUser(ctx, req.(*CreateUserRequest))
	}
	return interceptor(ctx, in, info, handler)
}

func _BankApp_LoginUser_Handler(srv interface{}, ctx context.Context, dec func(interface{}) error, interceptor grpc.UnaryServerInterceptor) (interface{}, error) {
	in := new(LoginUserRequest)
	if err := dec(in); err != nil {
		return nil, err
	}
	if interceptor == nil {
		return srv.(BankAppServer).LoginUser(ctx, in)
	}
	info := &grpc.UnaryServerInfo{
		Server:     srv,
		FullMethod: BankApp_LoginUser_FullMethodName,
	}
	handler := func(ctx context.Context, req interface{}) (interface{}, error) {
		return srv.(BankAppServer).LoginUser(ctx, req.(*LoginUserRequest))
	}
	return interceptor(ctx, in, info, handler)
}

// BankApp_ServiceDesc is the grpc.ServiceDesc for BankApp service.
// It's only intended for direct use with grpc.RegisterService,
// and not to be introspected or modified (even as a copy)
var BankApp_ServiceDesc = grpc.ServiceDesc{
	ServiceName: "pb.BankApp",
	HandlerType: (*BankAppServer)(nil),
	Methods: []grpc.MethodDesc{
		{
			MethodName: "CreateUser",
			Handler:    _BankApp_CreateUser_Handler,
		},
		{
			MethodName: "LoginUser",
			Handler:    _BankApp_LoginUser_Handler,
		},
	},
	Streams:  []grpc.StreamDesc{},
	Metadata: "service_bank_app.proto",
}
