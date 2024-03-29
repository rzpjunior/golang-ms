// Code generated by protoc-gen-go-grpc. DO NOT EDIT.
// versions:
// - protoc-gen-go-grpc v1.2.0
// - protoc             v4.23.4
// source: proto/boilerplate_service/boilerplate_service.proto

package boilerplate_service

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

// BoilerplateServiceClient is the client API for BoilerplateService service.
//
// For semantics around ctx use and closing/ending streaming RPCs, please refer to https://pkg.go.dev/google.golang.org/grpc/?tab=doc#ClientConn.NewStream.
type BoilerplateServiceClient interface {
	GetPerson(ctx context.Context, in *GetPersonRequest, opts ...grpc.CallOption) (*GetPersonResponse, error)
	GetPersonByID(ctx context.Context, in *GetPersonByIDRequest, opts ...grpc.CallOption) (*GetPersonByIDResponse, error)
	CreatePerson(ctx context.Context, in *CreatePersonRequest, opts ...grpc.CallOption) (*CreatePersonResponse, error)
	UpdatePerson(ctx context.Context, in *UpdatePersonRequest, opts ...grpc.CallOption) (*UpdatePersonResponse, error)
	DeletePerson(ctx context.Context, in *DeletePersonRequest, opts ...grpc.CallOption) (*DeletePersonResponse, error)
}

type boilerplateServiceClient struct {
	cc grpc.ClientConnInterface
}

func NewBoilerplateServiceClient(cc grpc.ClientConnInterface) BoilerplateServiceClient {
	return &boilerplateServiceClient{cc}
}

func (c *boilerplateServiceClient) GetPerson(ctx context.Context, in *GetPersonRequest, opts ...grpc.CallOption) (*GetPersonResponse, error) {
	out := new(GetPersonResponse)
	err := c.cc.Invoke(ctx, "/proto.boilerplate_service.BoilerplateService/GetPerson", in, out, opts...)
	if err != nil {
		return nil, err
	}
	return out, nil
}

func (c *boilerplateServiceClient) GetPersonByID(ctx context.Context, in *GetPersonByIDRequest, opts ...grpc.CallOption) (*GetPersonByIDResponse, error) {
	out := new(GetPersonByIDResponse)
	err := c.cc.Invoke(ctx, "/proto.boilerplate_service.BoilerplateService/GetPersonByID", in, out, opts...)
	if err != nil {
		return nil, err
	}
	return out, nil
}

func (c *boilerplateServiceClient) CreatePerson(ctx context.Context, in *CreatePersonRequest, opts ...grpc.CallOption) (*CreatePersonResponse, error) {
	out := new(CreatePersonResponse)
	err := c.cc.Invoke(ctx, "/proto.boilerplate_service.BoilerplateService/CreatePerson", in, out, opts...)
	if err != nil {
		return nil, err
	}
	return out, nil
}

func (c *boilerplateServiceClient) UpdatePerson(ctx context.Context, in *UpdatePersonRequest, opts ...grpc.CallOption) (*UpdatePersonResponse, error) {
	out := new(UpdatePersonResponse)
	err := c.cc.Invoke(ctx, "/proto.boilerplate_service.BoilerplateService/UpdatePerson", in, out, opts...)
	if err != nil {
		return nil, err
	}
	return out, nil
}

func (c *boilerplateServiceClient) DeletePerson(ctx context.Context, in *DeletePersonRequest, opts ...grpc.CallOption) (*DeletePersonResponse, error) {
	out := new(DeletePersonResponse)
	err := c.cc.Invoke(ctx, "/proto.boilerplate_service.BoilerplateService/DeletePerson", in, out, opts...)
	if err != nil {
		return nil, err
	}
	return out, nil
}

// BoilerplateServiceServer is the server API for BoilerplateService service.
// All implementations should embed UnimplementedBoilerplateServiceServer
// for forward compatibility
type BoilerplateServiceServer interface {
	GetPerson(context.Context, *GetPersonRequest) (*GetPersonResponse, error)
	GetPersonByID(context.Context, *GetPersonByIDRequest) (*GetPersonByIDResponse, error)
	CreatePerson(context.Context, *CreatePersonRequest) (*CreatePersonResponse, error)
	UpdatePerson(context.Context, *UpdatePersonRequest) (*UpdatePersonResponse, error)
	DeletePerson(context.Context, *DeletePersonRequest) (*DeletePersonResponse, error)
}

// UnimplementedBoilerplateServiceServer should be embedded to have forward compatible implementations.
type UnimplementedBoilerplateServiceServer struct {
}

func (UnimplementedBoilerplateServiceServer) GetPerson(context.Context, *GetPersonRequest) (*GetPersonResponse, error) {
	return nil, status.Errorf(codes.Unimplemented, "method GetPerson not implemented")
}
func (UnimplementedBoilerplateServiceServer) GetPersonByID(context.Context, *GetPersonByIDRequest) (*GetPersonByIDResponse, error) {
	return nil, status.Errorf(codes.Unimplemented, "method GetPersonByID not implemented")
}
func (UnimplementedBoilerplateServiceServer) CreatePerson(context.Context, *CreatePersonRequest) (*CreatePersonResponse, error) {
	return nil, status.Errorf(codes.Unimplemented, "method CreatePerson not implemented")
}
func (UnimplementedBoilerplateServiceServer) UpdatePerson(context.Context, *UpdatePersonRequest) (*UpdatePersonResponse, error) {
	return nil, status.Errorf(codes.Unimplemented, "method UpdatePerson not implemented")
}
func (UnimplementedBoilerplateServiceServer) DeletePerson(context.Context, *DeletePersonRequest) (*DeletePersonResponse, error) {
	return nil, status.Errorf(codes.Unimplemented, "method DeletePerson not implemented")
}

// UnsafeBoilerplateServiceServer may be embedded to opt out of forward compatibility for this service.
// Use of this interface is not recommended, as added methods to BoilerplateServiceServer will
// result in compilation errors.
type UnsafeBoilerplateServiceServer interface {
	mustEmbedUnimplementedBoilerplateServiceServer()
}

func RegisterBoilerplateServiceServer(s grpc.ServiceRegistrar, srv BoilerplateServiceServer) {
	s.RegisterService(&BoilerplateService_ServiceDesc, srv)
}

func _BoilerplateService_GetPerson_Handler(srv interface{}, ctx context.Context, dec func(interface{}) error, interceptor grpc.UnaryServerInterceptor) (interface{}, error) {
	in := new(GetPersonRequest)
	if err := dec(in); err != nil {
		return nil, err
	}
	if interceptor == nil {
		return srv.(BoilerplateServiceServer).GetPerson(ctx, in)
	}
	info := &grpc.UnaryServerInfo{
		Server:     srv,
		FullMethod: "/proto.boilerplate_service.BoilerplateService/GetPerson",
	}
	handler := func(ctx context.Context, req interface{}) (interface{}, error) {
		return srv.(BoilerplateServiceServer).GetPerson(ctx, req.(*GetPersonRequest))
	}
	return interceptor(ctx, in, info, handler)
}

func _BoilerplateService_GetPersonByID_Handler(srv interface{}, ctx context.Context, dec func(interface{}) error, interceptor grpc.UnaryServerInterceptor) (interface{}, error) {
	in := new(GetPersonByIDRequest)
	if err := dec(in); err != nil {
		return nil, err
	}
	if interceptor == nil {
		return srv.(BoilerplateServiceServer).GetPersonByID(ctx, in)
	}
	info := &grpc.UnaryServerInfo{
		Server:     srv,
		FullMethod: "/proto.boilerplate_service.BoilerplateService/GetPersonByID",
	}
	handler := func(ctx context.Context, req interface{}) (interface{}, error) {
		return srv.(BoilerplateServiceServer).GetPersonByID(ctx, req.(*GetPersonByIDRequest))
	}
	return interceptor(ctx, in, info, handler)
}

func _BoilerplateService_CreatePerson_Handler(srv interface{}, ctx context.Context, dec func(interface{}) error, interceptor grpc.UnaryServerInterceptor) (interface{}, error) {
	in := new(CreatePersonRequest)
	if err := dec(in); err != nil {
		return nil, err
	}
	if interceptor == nil {
		return srv.(BoilerplateServiceServer).CreatePerson(ctx, in)
	}
	info := &grpc.UnaryServerInfo{
		Server:     srv,
		FullMethod: "/proto.boilerplate_service.BoilerplateService/CreatePerson",
	}
	handler := func(ctx context.Context, req interface{}) (interface{}, error) {
		return srv.(BoilerplateServiceServer).CreatePerson(ctx, req.(*CreatePersonRequest))
	}
	return interceptor(ctx, in, info, handler)
}

func _BoilerplateService_UpdatePerson_Handler(srv interface{}, ctx context.Context, dec func(interface{}) error, interceptor grpc.UnaryServerInterceptor) (interface{}, error) {
	in := new(UpdatePersonRequest)
	if err := dec(in); err != nil {
		return nil, err
	}
	if interceptor == nil {
		return srv.(BoilerplateServiceServer).UpdatePerson(ctx, in)
	}
	info := &grpc.UnaryServerInfo{
		Server:     srv,
		FullMethod: "/proto.boilerplate_service.BoilerplateService/UpdatePerson",
	}
	handler := func(ctx context.Context, req interface{}) (interface{}, error) {
		return srv.(BoilerplateServiceServer).UpdatePerson(ctx, req.(*UpdatePersonRequest))
	}
	return interceptor(ctx, in, info, handler)
}

func _BoilerplateService_DeletePerson_Handler(srv interface{}, ctx context.Context, dec func(interface{}) error, interceptor grpc.UnaryServerInterceptor) (interface{}, error) {
	in := new(DeletePersonRequest)
	if err := dec(in); err != nil {
		return nil, err
	}
	if interceptor == nil {
		return srv.(BoilerplateServiceServer).DeletePerson(ctx, in)
	}
	info := &grpc.UnaryServerInfo{
		Server:     srv,
		FullMethod: "/proto.boilerplate_service.BoilerplateService/DeletePerson",
	}
	handler := func(ctx context.Context, req interface{}) (interface{}, error) {
		return srv.(BoilerplateServiceServer).DeletePerson(ctx, req.(*DeletePersonRequest))
	}
	return interceptor(ctx, in, info, handler)
}

// BoilerplateService_ServiceDesc is the grpc.ServiceDesc for BoilerplateService service.
// It's only intended for direct use with grpc.RegisterService,
// and not to be introspected or modified (even as a copy)
var BoilerplateService_ServiceDesc = grpc.ServiceDesc{
	ServiceName: "proto.boilerplate_service.BoilerplateService",
	HandlerType: (*BoilerplateServiceServer)(nil),
	Methods: []grpc.MethodDesc{
		{
			MethodName: "GetPerson",
			Handler:    _BoilerplateService_GetPerson_Handler,
		},
		{
			MethodName: "GetPersonByID",
			Handler:    _BoilerplateService_GetPersonByID_Handler,
		},
		{
			MethodName: "CreatePerson",
			Handler:    _BoilerplateService_CreatePerson_Handler,
		},
		{
			MethodName: "UpdatePerson",
			Handler:    _BoilerplateService_UpdatePerson_Handler,
		},
		{
			MethodName: "DeletePerson",
			Handler:    _BoilerplateService_DeletePerson_Handler,
		},
	},
	Streams:  []grpc.StreamDesc{},
	Metadata: "proto/boilerplate_service/boilerplate_service.proto",
}
