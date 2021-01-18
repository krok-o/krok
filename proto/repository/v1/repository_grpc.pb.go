// Code generated by protoc-gen-go-grpc. DO NOT EDIT.

package repov1

import (
	context "context"
	empty "github.com/golang/protobuf/ptypes/empty"
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
	CreateRepository(ctx context.Context, in *CreateRepositoryRequest, opts ...grpc.CallOption) (*Repository, error)
	UpdateRepository(ctx context.Context, in *UpdateRepositoryRequest, opts ...grpc.CallOption) (*Repository, error)
	GetRepository(ctx context.Context, in *GetRepositoryRequest, opts ...grpc.CallOption) (*Repository, error)
	ListRepositories(ctx context.Context, in *ListRepositoryRequest, opts ...grpc.CallOption) (*Repositories, error)
	DeleteRepository(ctx context.Context, in *DeleteRepositoryRequest, opts ...grpc.CallOption) (*empty.Empty, error)
}

type repositoryServiceClient struct {
	cc grpc.ClientConnInterface
}

func NewRepositoryServiceClient(cc grpc.ClientConnInterface) RepositoryServiceClient {
	return &repositoryServiceClient{cc}
}

func (c *repositoryServiceClient) CreateRepository(ctx context.Context, in *CreateRepositoryRequest, opts ...grpc.CallOption) (*Repository, error) {
	out := new(Repository)
	err := c.cc.Invoke(ctx, "/repository.v1.RepositoryService/CreateRepository", in, out, opts...)
	if err != nil {
		return nil, err
	}
	return out, nil
}

func (c *repositoryServiceClient) UpdateRepository(ctx context.Context, in *UpdateRepositoryRequest, opts ...grpc.CallOption) (*Repository, error) {
	out := new(Repository)
	err := c.cc.Invoke(ctx, "/repository.v1.RepositoryService/UpdateRepository", in, out, opts...)
	if err != nil {
		return nil, err
	}
	return out, nil
}

func (c *repositoryServiceClient) GetRepository(ctx context.Context, in *GetRepositoryRequest, opts ...grpc.CallOption) (*Repository, error) {
	out := new(Repository)
	err := c.cc.Invoke(ctx, "/repository.v1.RepositoryService/GetRepository", in, out, opts...)
	if err != nil {
		return nil, err
	}
	return out, nil
}

func (c *repositoryServiceClient) ListRepositories(ctx context.Context, in *ListRepositoryRequest, opts ...grpc.CallOption) (*Repositories, error) {
	out := new(Repositories)
	err := c.cc.Invoke(ctx, "/repository.v1.RepositoryService/ListRepositories", in, out, opts...)
	if err != nil {
		return nil, err
	}
	return out, nil
}

func (c *repositoryServiceClient) DeleteRepository(ctx context.Context, in *DeleteRepositoryRequest, opts ...grpc.CallOption) (*empty.Empty, error) {
	out := new(empty.Empty)
	err := c.cc.Invoke(ctx, "/repository.v1.RepositoryService/DeleteRepository", in, out, opts...)
	if err != nil {
		return nil, err
	}
	return out, nil
}

// RepositoryServiceServer is the server API for RepositoryService service.
// All implementations must embed UnimplementedRepositoryServiceServer
// for forward compatibility
type RepositoryServiceServer interface {
	CreateRepository(context.Context, *CreateRepositoryRequest) (*Repository, error)
	UpdateRepository(context.Context, *UpdateRepositoryRequest) (*Repository, error)
	GetRepository(context.Context, *GetRepositoryRequest) (*Repository, error)
	ListRepositories(context.Context, *ListRepositoryRequest) (*Repositories, error)
	DeleteRepository(context.Context, *DeleteRepositoryRequest) (*empty.Empty, error)
	mustEmbedUnimplementedRepositoryServiceServer()
}

// UnimplementedRepositoryServiceServer must be embedded to have forward compatible implementations.
type UnimplementedRepositoryServiceServer struct {
}

func (UnimplementedRepositoryServiceServer) CreateRepository(context.Context, *CreateRepositoryRequest) (*Repository, error) {
	return nil, status.Errorf(codes.Unimplemented, "method CreateRepository not implemented")
}
func (UnimplementedRepositoryServiceServer) UpdateRepository(context.Context, *UpdateRepositoryRequest) (*Repository, error) {
	return nil, status.Errorf(codes.Unimplemented, "method UpdateRepository not implemented")
}
func (UnimplementedRepositoryServiceServer) GetRepository(context.Context, *GetRepositoryRequest) (*Repository, error) {
	return nil, status.Errorf(codes.Unimplemented, "method GetRepository not implemented")
}
func (UnimplementedRepositoryServiceServer) ListRepositories(context.Context, *ListRepositoryRequest) (*Repositories, error) {
	return nil, status.Errorf(codes.Unimplemented, "method ListRepositories not implemented")
}
func (UnimplementedRepositoryServiceServer) DeleteRepository(context.Context, *DeleteRepositoryRequest) (*empty.Empty, error) {
	return nil, status.Errorf(codes.Unimplemented, "method DeleteRepository not implemented")
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
	in := new(CreateRepositoryRequest)
	if err := dec(in); err != nil {
		return nil, err
	}
	if interceptor == nil {
		return srv.(RepositoryServiceServer).CreateRepository(ctx, in)
	}
	info := &grpc.UnaryServerInfo{
		Server:     srv,
		FullMethod: "/repository.v1.RepositoryService/CreateRepository",
	}
	handler := func(ctx context.Context, req interface{}) (interface{}, error) {
		return srv.(RepositoryServiceServer).CreateRepository(ctx, req.(*CreateRepositoryRequest))
	}
	return interceptor(ctx, in, info, handler)
}

func _RepositoryService_UpdateRepository_Handler(srv interface{}, ctx context.Context, dec func(interface{}) error, interceptor grpc.UnaryServerInterceptor) (interface{}, error) {
	in := new(UpdateRepositoryRequest)
	if err := dec(in); err != nil {
		return nil, err
	}
	if interceptor == nil {
		return srv.(RepositoryServiceServer).UpdateRepository(ctx, in)
	}
	info := &grpc.UnaryServerInfo{
		Server:     srv,
		FullMethod: "/repository.v1.RepositoryService/UpdateRepository",
	}
	handler := func(ctx context.Context, req interface{}) (interface{}, error) {
		return srv.(RepositoryServiceServer).UpdateRepository(ctx, req.(*UpdateRepositoryRequest))
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
		FullMethod: "/repository.v1.RepositoryService/GetRepository",
	}
	handler := func(ctx context.Context, req interface{}) (interface{}, error) {
		return srv.(RepositoryServiceServer).GetRepository(ctx, req.(*GetRepositoryRequest))
	}
	return interceptor(ctx, in, info, handler)
}

func _RepositoryService_ListRepositories_Handler(srv interface{}, ctx context.Context, dec func(interface{}) error, interceptor grpc.UnaryServerInterceptor) (interface{}, error) {
	in := new(ListRepositoryRequest)
	if err := dec(in); err != nil {
		return nil, err
	}
	if interceptor == nil {
		return srv.(RepositoryServiceServer).ListRepositories(ctx, in)
	}
	info := &grpc.UnaryServerInfo{
		Server:     srv,
		FullMethod: "/repository.v1.RepositoryService/ListRepositories",
	}
	handler := func(ctx context.Context, req interface{}) (interface{}, error) {
		return srv.(RepositoryServiceServer).ListRepositories(ctx, req.(*ListRepositoryRequest))
	}
	return interceptor(ctx, in, info, handler)
}

func _RepositoryService_DeleteRepository_Handler(srv interface{}, ctx context.Context, dec func(interface{}) error, interceptor grpc.UnaryServerInterceptor) (interface{}, error) {
	in := new(DeleteRepositoryRequest)
	if err := dec(in); err != nil {
		return nil, err
	}
	if interceptor == nil {
		return srv.(RepositoryServiceServer).DeleteRepository(ctx, in)
	}
	info := &grpc.UnaryServerInfo{
		Server:     srv,
		FullMethod: "/repository.v1.RepositoryService/DeleteRepository",
	}
	handler := func(ctx context.Context, req interface{}) (interface{}, error) {
		return srv.(RepositoryServiceServer).DeleteRepository(ctx, req.(*DeleteRepositoryRequest))
	}
	return interceptor(ctx, in, info, handler)
}

var _RepositoryService_serviceDesc = grpc.ServiceDesc{
	ServiceName: "repository.v1.RepositoryService",
	HandlerType: (*RepositoryServiceServer)(nil),
	Methods: []grpc.MethodDesc{
		{
			MethodName: "CreateRepository",
			Handler:    _RepositoryService_CreateRepository_Handler,
		},
		{
			MethodName: "UpdateRepository",
			Handler:    _RepositoryService_UpdateRepository_Handler,
		},
		{
			MethodName: "GetRepository",
			Handler:    _RepositoryService_GetRepository_Handler,
		},
		{
			MethodName: "ListRepositories",
			Handler:    _RepositoryService_ListRepositories_Handler,
		},
		{
			MethodName: "DeleteRepository",
			Handler:    _RepositoryService_DeleteRepository_Handler,
		},
	},
	Streams:  []grpc.StreamDesc{},
	Metadata: "repository/v1/repository.proto",
}
