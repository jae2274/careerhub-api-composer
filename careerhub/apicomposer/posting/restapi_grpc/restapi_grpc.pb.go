// Code generated by protoc-gen-go-grpc. DO NOT EDIT.
// versions:
// - protoc-gen-go-grpc v1.2.0
// - protoc             v4.25.1
// source: careerhub/apicomposer/posting/restapi_grpc/restapi.proto

package restapi_grpc

import (
	context "context"
	grpc "google.golang.org/grpc"
	codes "google.golang.org/grpc/codes"
	status "google.golang.org/grpc/status"
	emptypb "google.golang.org/protobuf/types/known/emptypb"
)

// This is a compile-time assertion to ensure that this generated file
// is compatible with the grpc package it is being compiled against.
// Requires gRPC-Go v1.32.0 or later.
const _ = grpc.SupportPackageIsVersion7

// RestApiGrpcClient is the client API for RestApiGrpc service.
//
// For semantics around ctx use and closing/ending streaming RPCs, please refer to https://pkg.go.dev/google.golang.org/grpc/?tab=doc#ClientConn.NewStream.
type RestApiGrpcClient interface {
	JobPostings(ctx context.Context, in *JobPostingsRequest, opts ...grpc.CallOption) (*JobPostingsResponse, error)
	JobPostingDetail(ctx context.Context, in *JobPostingDetailRequest, opts ...grpc.CallOption) (*JobPostingDetailResponse, error)
	Categories(ctx context.Context, in *emptypb.Empty, opts ...grpc.CallOption) (*CategoriesResponse, error)
	Skills(ctx context.Context, in *emptypb.Empty, opts ...grpc.CallOption) (*SkillsResponse, error)
	JobPostingsById(ctx context.Context, in *JobPostingsByIdRequest, opts ...grpc.CallOption) (*JobPostingsResponse, error)
}

type restApiGrpcClient struct {
	cc grpc.ClientConnInterface
}

func NewRestApiGrpcClient(cc grpc.ClientConnInterface) RestApiGrpcClient {
	return &restApiGrpcClient{cc}
}

func (c *restApiGrpcClient) JobPostings(ctx context.Context, in *JobPostingsRequest, opts ...grpc.CallOption) (*JobPostingsResponse, error) {
	out := new(JobPostingsResponse)
	err := c.cc.Invoke(ctx, "/careerhub.posting_service.restapi_grpc.RestApiGrpc/JobPostings", in, out, opts...)
	if err != nil {
		return nil, err
	}
	return out, nil
}

func (c *restApiGrpcClient) JobPostingDetail(ctx context.Context, in *JobPostingDetailRequest, opts ...grpc.CallOption) (*JobPostingDetailResponse, error) {
	out := new(JobPostingDetailResponse)
	err := c.cc.Invoke(ctx, "/careerhub.posting_service.restapi_grpc.RestApiGrpc/JobPostingDetail", in, out, opts...)
	if err != nil {
		return nil, err
	}
	return out, nil
}

func (c *restApiGrpcClient) Categories(ctx context.Context, in *emptypb.Empty, opts ...grpc.CallOption) (*CategoriesResponse, error) {
	out := new(CategoriesResponse)
	err := c.cc.Invoke(ctx, "/careerhub.posting_service.restapi_grpc.RestApiGrpc/Categories", in, out, opts...)
	if err != nil {
		return nil, err
	}
	return out, nil
}

func (c *restApiGrpcClient) Skills(ctx context.Context, in *emptypb.Empty, opts ...grpc.CallOption) (*SkillsResponse, error) {
	out := new(SkillsResponse)
	err := c.cc.Invoke(ctx, "/careerhub.posting_service.restapi_grpc.RestApiGrpc/Skills", in, out, opts...)
	if err != nil {
		return nil, err
	}
	return out, nil
}

func (c *restApiGrpcClient) JobPostingsById(ctx context.Context, in *JobPostingsByIdRequest, opts ...grpc.CallOption) (*JobPostingsResponse, error) {
	out := new(JobPostingsResponse)
	err := c.cc.Invoke(ctx, "/careerhub.posting_service.restapi_grpc.RestApiGrpc/JobPostingsById", in, out, opts...)
	if err != nil {
		return nil, err
	}
	return out, nil
}

// RestApiGrpcServer is the server API for RestApiGrpc service.
// All implementations must embed UnimplementedRestApiGrpcServer
// for forward compatibility
type RestApiGrpcServer interface {
	JobPostings(context.Context, *JobPostingsRequest) (*JobPostingsResponse, error)
	JobPostingDetail(context.Context, *JobPostingDetailRequest) (*JobPostingDetailResponse, error)
	Categories(context.Context, *emptypb.Empty) (*CategoriesResponse, error)
	Skills(context.Context, *emptypb.Empty) (*SkillsResponse, error)
	JobPostingsById(context.Context, *JobPostingsByIdRequest) (*JobPostingsResponse, error)
	mustEmbedUnimplementedRestApiGrpcServer()
}

// UnimplementedRestApiGrpcServer must be embedded to have forward compatible implementations.
type UnimplementedRestApiGrpcServer struct {
}

func (UnimplementedRestApiGrpcServer) JobPostings(context.Context, *JobPostingsRequest) (*JobPostingsResponse, error) {
	return nil, status.Errorf(codes.Unimplemented, "method JobPostings not implemented")
}
func (UnimplementedRestApiGrpcServer) JobPostingDetail(context.Context, *JobPostingDetailRequest) (*JobPostingDetailResponse, error) {
	return nil, status.Errorf(codes.Unimplemented, "method JobPostingDetail not implemented")
}
func (UnimplementedRestApiGrpcServer) Categories(context.Context, *emptypb.Empty) (*CategoriesResponse, error) {
	return nil, status.Errorf(codes.Unimplemented, "method Categories not implemented")
}
func (UnimplementedRestApiGrpcServer) Skills(context.Context, *emptypb.Empty) (*SkillsResponse, error) {
	return nil, status.Errorf(codes.Unimplemented, "method Skills not implemented")
}
func (UnimplementedRestApiGrpcServer) JobPostingsById(context.Context, *JobPostingsByIdRequest) (*JobPostingsResponse, error) {
	return nil, status.Errorf(codes.Unimplemented, "method JobPostingsById not implemented")
}
func (UnimplementedRestApiGrpcServer) mustEmbedUnimplementedRestApiGrpcServer() {}

// UnsafeRestApiGrpcServer may be embedded to opt out of forward compatibility for this service.
// Use of this interface is not recommended, as added methods to RestApiGrpcServer will
// result in compilation errors.
type UnsafeRestApiGrpcServer interface {
	mustEmbedUnimplementedRestApiGrpcServer()
}

func RegisterRestApiGrpcServer(s grpc.ServiceRegistrar, srv RestApiGrpcServer) {
	s.RegisterService(&RestApiGrpc_ServiceDesc, srv)
}

func _RestApiGrpc_JobPostings_Handler(srv interface{}, ctx context.Context, dec func(interface{}) error, interceptor grpc.UnaryServerInterceptor) (interface{}, error) {
	in := new(JobPostingsRequest)
	if err := dec(in); err != nil {
		return nil, err
	}
	if interceptor == nil {
		return srv.(RestApiGrpcServer).JobPostings(ctx, in)
	}
	info := &grpc.UnaryServerInfo{
		Server:     srv,
		FullMethod: "/careerhub.posting_service.restapi_grpc.RestApiGrpc/JobPostings",
	}
	handler := func(ctx context.Context, req interface{}) (interface{}, error) {
		return srv.(RestApiGrpcServer).JobPostings(ctx, req.(*JobPostingsRequest))
	}
	return interceptor(ctx, in, info, handler)
}

func _RestApiGrpc_JobPostingDetail_Handler(srv interface{}, ctx context.Context, dec func(interface{}) error, interceptor grpc.UnaryServerInterceptor) (interface{}, error) {
	in := new(JobPostingDetailRequest)
	if err := dec(in); err != nil {
		return nil, err
	}
	if interceptor == nil {
		return srv.(RestApiGrpcServer).JobPostingDetail(ctx, in)
	}
	info := &grpc.UnaryServerInfo{
		Server:     srv,
		FullMethod: "/careerhub.posting_service.restapi_grpc.RestApiGrpc/JobPostingDetail",
	}
	handler := func(ctx context.Context, req interface{}) (interface{}, error) {
		return srv.(RestApiGrpcServer).JobPostingDetail(ctx, req.(*JobPostingDetailRequest))
	}
	return interceptor(ctx, in, info, handler)
}

func _RestApiGrpc_Categories_Handler(srv interface{}, ctx context.Context, dec func(interface{}) error, interceptor grpc.UnaryServerInterceptor) (interface{}, error) {
	in := new(emptypb.Empty)
	if err := dec(in); err != nil {
		return nil, err
	}
	if interceptor == nil {
		return srv.(RestApiGrpcServer).Categories(ctx, in)
	}
	info := &grpc.UnaryServerInfo{
		Server:     srv,
		FullMethod: "/careerhub.posting_service.restapi_grpc.RestApiGrpc/Categories",
	}
	handler := func(ctx context.Context, req interface{}) (interface{}, error) {
		return srv.(RestApiGrpcServer).Categories(ctx, req.(*emptypb.Empty))
	}
	return interceptor(ctx, in, info, handler)
}

func _RestApiGrpc_Skills_Handler(srv interface{}, ctx context.Context, dec func(interface{}) error, interceptor grpc.UnaryServerInterceptor) (interface{}, error) {
	in := new(emptypb.Empty)
	if err := dec(in); err != nil {
		return nil, err
	}
	if interceptor == nil {
		return srv.(RestApiGrpcServer).Skills(ctx, in)
	}
	info := &grpc.UnaryServerInfo{
		Server:     srv,
		FullMethod: "/careerhub.posting_service.restapi_grpc.RestApiGrpc/Skills",
	}
	handler := func(ctx context.Context, req interface{}) (interface{}, error) {
		return srv.(RestApiGrpcServer).Skills(ctx, req.(*emptypb.Empty))
	}
	return interceptor(ctx, in, info, handler)
}

func _RestApiGrpc_JobPostingsById_Handler(srv interface{}, ctx context.Context, dec func(interface{}) error, interceptor grpc.UnaryServerInterceptor) (interface{}, error) {
	in := new(JobPostingsByIdRequest)
	if err := dec(in); err != nil {
		return nil, err
	}
	if interceptor == nil {
		return srv.(RestApiGrpcServer).JobPostingsById(ctx, in)
	}
	info := &grpc.UnaryServerInfo{
		Server:     srv,
		FullMethod: "/careerhub.posting_service.restapi_grpc.RestApiGrpc/JobPostingsById",
	}
	handler := func(ctx context.Context, req interface{}) (interface{}, error) {
		return srv.(RestApiGrpcServer).JobPostingsById(ctx, req.(*JobPostingsByIdRequest))
	}
	return interceptor(ctx, in, info, handler)
}

// RestApiGrpc_ServiceDesc is the grpc.ServiceDesc for RestApiGrpc service.
// It's only intended for direct use with grpc.RegisterService,
// and not to be introspected or modified (even as a copy)
var RestApiGrpc_ServiceDesc = grpc.ServiceDesc{
	ServiceName: "careerhub.posting_service.restapi_grpc.RestApiGrpc",
	HandlerType: (*RestApiGrpcServer)(nil),
	Methods: []grpc.MethodDesc{
		{
			MethodName: "JobPostings",
			Handler:    _RestApiGrpc_JobPostings_Handler,
		},
		{
			MethodName: "JobPostingDetail",
			Handler:    _RestApiGrpc_JobPostingDetail_Handler,
		},
		{
			MethodName: "Categories",
			Handler:    _RestApiGrpc_Categories_Handler,
		},
		{
			MethodName: "Skills",
			Handler:    _RestApiGrpc_Skills_Handler,
		},
		{
			MethodName: "JobPostingsById",
			Handler:    _RestApiGrpc_JobPostingsById_Handler,
		},
	},
	Streams:  []grpc.StreamDesc{},
	Metadata: "careerhub/apicomposer/posting/restapi_grpc/restapi.proto",
}
