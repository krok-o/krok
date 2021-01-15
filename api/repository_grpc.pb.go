// Code generated by protoc-gen-go-grpc. DO NOT EDIT.

package api

import (
	context "context"
	grpc "google.golang.org/grpc"
	codes "google.golang.org/grpc/codes"
	status "google.golang.org/grpc/status"
)

// This is a compile-time assertion to ensure that this generated file
// is compatible with the grpc package it is being compiled against.
const _ = grpc.SupportPackageIsVersion7

// RepositoryServiceClient is the client API for RepositoryService service.
//
// For semantics around ctx use and closing/ending streaming RPCs, please refer to https://pkg.go.dev/google.golang.org/grpc/?tab=doc#ClientConn.NewStream.
type RepositoryServiceClient interface {
	CreateRepository(ctx context.Context, in *Repository, opts ...grpc.CallOption) (*Repository, error)
	GetRepository(ctx context.Context, in *GetRepositoryRequest, opts ...grpc.CallOption) (*Repository, error)
}

type repositoryServiceClient struct {
	cc grpc.ClientConnInterface
}

func NewRepositoryServiceClient(cc grpc.ClientConnInterface) RepositoryServiceClient {
	return &repositoryServiceClient{cc}
}

func (c *repositoryServiceClient) CreateRepository(ctx context.Context, in *Repository, opts ...grpc.CallOption) (*Repository, error) {
	out := new(Repository)
	err := c.cc.Invoke(ctx, "/v1.RepositoryService/CreateRepository", in, out, opts...)
	if err != nil {
		return nil, err
	}
	return out, nil
}

func (c *repositoryServiceClient) GetRepository(ctx context.Context, in *GetRepositoryRequest, opts ...grpc.CallOption) (*Repository, error) {
	out := new(Repository)
	err := c.cc.Invoke(ctx, "/v1.RepositoryService/GetRepository", in, out, opts...)
	if err != nil {
		return nil, err
	}
	return out, nil
}

// RepositoryServiceServer is the server API for RepositoryService service.
// All implementations must embed UnimplementedRepositoryServiceServer
// for forward compatibility
type RepositoryServiceServer interface {
	CreateRepository(context.Context, *Repository) (*Repository, error)
	GetRepository(context.Context, *GetRepositoryRequest) (*Repository, error)
	mustEmbedUnimplementedRepositoryServiceServer()
}

// UnimplementedRepositoryServiceServer must be embedded to have forward compatible implementations.
type UnimplementedRepositoryServiceServer struct {
}

func (UnimplementedRepositoryServiceServer) CreateRepository(context.Context, *Repository) (*Repository, error) {
	return nil, status.Errorf(codes.Unimplemented, "method CreateRepository not implemented")
}
func (UnimplementedRepositoryServiceServer) GetRepository(context.Context, *GetRepositoryRequest) (*Repository, error) {
	return nil, status.Errorf(codes.Unimplemented, "method GetRepository not implemented")
}
func (UnimplementedRepositoryServiceServer) mustEmbedUnimplementedRepositoryServiceServer() {}

// UnsafeRepositoryServiceServer may be embedded to opt out of forward compatibility for this service.
// Use of this interface is not recommended, as added methods to RepositoryServiceServer will
// result in compilation errors.
type UnsafeRepositoryServiceServer interface {
	mustEmbedUnimplementedRepositoryServiceServer()
}

func RegisterRepositoryServiceServer(s grpc.ServiceRegistrar, srv RepositoryServiceServer) {
	s.RegisterService(&_RepositoryService_serviceDesc, srv)
}

func _RepositoryService_CreateRepository_Handler(srv interface{}, ctx context.Context, dec func(interface{}) error, interceptor grpc.UnaryServerInterceptor) (interface{}, error) {
	in := new(Repository)
	if err := dec(in); err != nil {
		return nil, err
	}
	if interceptor == nil {
		return srv.(RepositoryServiceServer).CreateRepository(ctx, in)
	}
	info := &grpc.UnaryServerInfo{
		Server:     srv,
		FullMethod: "/v1.RepositoryService/CreateRepository",
	}
	handler := func(ctx context.Context, req interface{}) (interface{}, error) {
		return srv.(RepositoryServiceServer).CreateRepository(ctx, req.(*Repository))
	}
	return interceptor(ctx, in, info, handler)
}

func _RepositoryService_GetRepository_Handler(srv interface{}, ctx context.Context, dec func(interface{}) error, interceptor grpc.UnaryServerInterceptor) (interface{}, error) {
	in := new(GetRepositoryRequest)
	if err := dec(in); err != nil {
		return nil, err
	}
	if interceptor == nil {
		return srv.(RepositoryServiceServer).GetRepository(ctx, in)
	}
	info := &grpc.UnaryServerInfo{
		Server:     srv,
		FullMethod: "/v1.RepositoryService/GetRepository",
	}
	handler := func(ctx context.Context, req interface{}) (interface{}, error) {
		return srv.(RepositoryServiceServer).GetRepository(ctx, req.(*GetRepositoryRequest))
	}
	return interceptor(ctx, in, info, handler)
}

var _RepositoryService_serviceDesc = grpc.ServiceDesc{
	ServiceName: "v1.RepositoryService",
	HandlerType: (*RepositoryServiceServer)(nil),
	Methods: []grpc.MethodDesc{
		{
			MethodName: "CreateRepository",
			Handler:    _RepositoryService_CreateRepository_Handler,
		},
		{
			MethodName: "GetRepository",
			Handler:    _RepositoryService_GetRepository_Handler,
		},
	},
	Streams:  []grpc.StreamDesc{},
	Metadata: "api/repository.proto",
}
