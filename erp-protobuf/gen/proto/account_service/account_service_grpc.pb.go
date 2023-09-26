// Code generated by protoc-gen-go-grpc. DO NOT EDIT.
// versions:
// - protoc-gen-go-grpc v1.2.0
// - protoc             v4.23.4
// source: proto/account_service/account_service.proto

package account_service

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

// AccountServiceClient is the client API for AccountService service.
//
// For semantics around ctx use and closing/ending streaming RPCs, please refer to https://pkg.go.dev/google.golang.org/grpc/?tab=doc#ClientConn.NewStream.
type AccountServiceClient interface {
	GetUserList(ctx context.Context, in *GetUserListRequest, opts ...grpc.CallOption) (*GetUserListResponse, error)
	GetUserDetail(ctx context.Context, in *GetUserDetailRequest, opts ...grpc.CallOption) (*GetUserDetailResponse, error)
	GetUserEmailAuth(ctx context.Context, in *GetUserEmailAuthRequest, opts ...grpc.CallOption) (*GetUserEmailAuthResponse, error)
	GetUserRoleByUserId(ctx context.Context, in *GetUserRoleByUserIdRequest, opts ...grpc.CallOption) (*GetUserRoleByUserIdResponse, error)
	UpdateUserSalesAppToken(ctx context.Context, in *UpdateUserSalesAppTokenRequest, opts ...grpc.CallOption) (*GetUserDetailResponse, error)
	GetUserBySalesAppLoginToken(ctx context.Context, in *GetUserBySalesAppLoginTokenRequest, opts ...grpc.CallOption) (*GetUserDetailResponse, error)
	GetRoleDetail(ctx context.Context, in *GetRoleDetailRequest, opts ...grpc.CallOption) (*GetRoleDetailResponse, error)
	UpdateUserEdnAppToken(ctx context.Context, in *UpdateUserEdnAppTokenRequest, opts ...grpc.CallOption) (*GetUserDetailResponse, error)
	UpdateUserPurchaserAppToken(ctx context.Context, in *UpdateUserPurchaserAppTokenRequest, opts ...grpc.CallOption) (*GetUserDetailResponse, error)
	GetUserByEdnAppLoginToken(ctx context.Context, in *GetUserByEdnAppLoginTokenRequest, opts ...grpc.CallOption) (*GetUserDetailResponse, error)
	GetDivisionDetail(ctx context.Context, in *GetDivisionDetailRequest, opts ...grpc.CallOption) (*GetDivisionDetailResponse, error)
	GetDivisionDefaultByCustomerType(ctx context.Context, in *GetDivisionDefaultByCustomerTypeRequest, opts ...grpc.CallOption) (*GetDivisionDefaultByCustomerTypeResponse, error)
}

type accountServiceClient struct {
	cc grpc.ClientConnInterface
}

func NewAccountServiceClient(cc grpc.ClientConnInterface) AccountServiceClient {
	return &accountServiceClient{cc}
}

func (c *accountServiceClient) GetUserList(ctx context.Context, in *GetUserListRequest, opts ...grpc.CallOption) (*GetUserListResponse, error) {
	out := new(GetUserListResponse)
	err := c.cc.Invoke(ctx, "/proto.account_service.AccountService/GetUserList", in, out, opts...)
	if err != nil {
		return nil, err
	}
	return out, nil
}

func (c *accountServiceClient) GetUserDetail(ctx context.Context, in *GetUserDetailRequest, opts ...grpc.CallOption) (*GetUserDetailResponse, error) {
	out := new(GetUserDetailResponse)
	err := c.cc.Invoke(ctx, "/proto.account_service.AccountService/GetUserDetail", in, out, opts...)
	if err != nil {
		return nil, err
	}
	return out, nil
}

func (c *accountServiceClient) GetUserEmailAuth(ctx context.Context, in *GetUserEmailAuthRequest, opts ...grpc.CallOption) (*GetUserEmailAuthResponse, error) {
	out := new(GetUserEmailAuthResponse)
	err := c.cc.Invoke(ctx, "/proto.account_service.AccountService/GetUserEmailAuth", in, out, opts...)
	if err != nil {
		return nil, err
	}
	return out, nil
}

func (c *accountServiceClient) GetUserRoleByUserId(ctx context.Context, in *GetUserRoleByUserIdRequest, opts ...grpc.CallOption) (*GetUserRoleByUserIdResponse, error) {
	out := new(GetUserRoleByUserIdResponse)
	err := c.cc.Invoke(ctx, "/proto.account_service.AccountService/GetUserRoleByUserId", in, out, opts...)
	if err != nil {
		return nil, err
	}
	return out, nil
}

func (c *accountServiceClient) UpdateUserSalesAppToken(ctx context.Context, in *UpdateUserSalesAppTokenRequest, opts ...grpc.CallOption) (*GetUserDetailResponse, error) {
	out := new(GetUserDetailResponse)
	err := c.cc.Invoke(ctx, "/proto.account_service.AccountService/UpdateUserSalesAppToken", in, out, opts...)
	if err != nil {
		return nil, err
	}
	return out, nil
}

func (c *accountServiceClient) GetUserBySalesAppLoginToken(ctx context.Context, in *GetUserBySalesAppLoginTokenRequest, opts ...grpc.CallOption) (*GetUserDetailResponse, error) {
	out := new(GetUserDetailResponse)
	err := c.cc.Invoke(ctx, "/proto.account_service.AccountService/GetUserBySalesAppLoginToken", in, out, opts...)
	if err != nil {
		return nil, err
	}
	return out, nil
}

func (c *accountServiceClient) GetRoleDetail(ctx context.Context, in *GetRoleDetailRequest, opts ...grpc.CallOption) (*GetRoleDetailResponse, error) {
	out := new(GetRoleDetailResponse)
	err := c.cc.Invoke(ctx, "/proto.account_service.AccountService/GetRoleDetail", in, out, opts...)
	if err != nil {
		return nil, err
	}
	return out, nil
}

func (c *accountServiceClient) UpdateUserEdnAppToken(ctx context.Context, in *UpdateUserEdnAppTokenRequest, opts ...grpc.CallOption) (*GetUserDetailResponse, error) {
	out := new(GetUserDetailResponse)
	err := c.cc.Invoke(ctx, "/proto.account_service.AccountService/UpdateUserEdnAppToken", in, out, opts...)
	if err != nil {
		return nil, err
	}
	return out, nil
}

func (c *accountServiceClient) UpdateUserPurchaserAppToken(ctx context.Context, in *UpdateUserPurchaserAppTokenRequest, opts ...grpc.CallOption) (*GetUserDetailResponse, error) {
	out := new(GetUserDetailResponse)
	err := c.cc.Invoke(ctx, "/proto.account_service.AccountService/UpdateUserPurchaserAppToken", in, out, opts...)
	if err != nil {
		return nil, err
	}
	return out, nil
}

func (c *accountServiceClient) GetUserByEdnAppLoginToken(ctx context.Context, in *GetUserByEdnAppLoginTokenRequest, opts ...grpc.CallOption) (*GetUserDetailResponse, error) {
	out := new(GetUserDetailResponse)
	err := c.cc.Invoke(ctx, "/proto.account_service.AccountService/GetUserByEdnAppLoginToken", in, out, opts...)
	if err != nil {
		return nil, err
	}
	return out, nil
}

func (c *accountServiceClient) GetDivisionDetail(ctx context.Context, in *GetDivisionDetailRequest, opts ...grpc.CallOption) (*GetDivisionDetailResponse, error) {
	out := new(GetDivisionDetailResponse)
	err := c.cc.Invoke(ctx, "/proto.account_service.AccountService/GetDivisionDetail", in, out, opts...)
	if err != nil {
		return nil, err
	}
	return out, nil
}

func (c *accountServiceClient) GetDivisionDefaultByCustomerType(ctx context.Context, in *GetDivisionDefaultByCustomerTypeRequest, opts ...grpc.CallOption) (*GetDivisionDefaultByCustomerTypeResponse, error) {
	out := new(GetDivisionDefaultByCustomerTypeResponse)
	err := c.cc.Invoke(ctx, "/proto.account_service.AccountService/GetDivisionDefaultByCustomerType", in, out, opts...)
	if err != nil {
		return nil, err
	}
	return out, nil
}

// AccountServiceServer is the server API for AccountService service.
// All implementations should embed UnimplementedAccountServiceServer
// for forward compatibility
type AccountServiceServer interface {
	GetUserList(context.Context, *GetUserListRequest) (*GetUserListResponse, error)
	GetUserDetail(context.Context, *GetUserDetailRequest) (*GetUserDetailResponse, error)
	GetUserEmailAuth(context.Context, *GetUserEmailAuthRequest) (*GetUserEmailAuthResponse, error)
	GetUserRoleByUserId(context.Context, *GetUserRoleByUserIdRequest) (*GetUserRoleByUserIdResponse, error)
	UpdateUserSalesAppToken(context.Context, *UpdateUserSalesAppTokenRequest) (*GetUserDetailResponse, error)
	GetUserBySalesAppLoginToken(context.Context, *GetUserBySalesAppLoginTokenRequest) (*GetUserDetailResponse, error)
	GetRoleDetail(context.Context, *GetRoleDetailRequest) (*GetRoleDetailResponse, error)
	UpdateUserEdnAppToken(context.Context, *UpdateUserEdnAppTokenRequest) (*GetUserDetailResponse, error)
	UpdateUserPurchaserAppToken(context.Context, *UpdateUserPurchaserAppTokenRequest) (*GetUserDetailResponse, error)
	GetUserByEdnAppLoginToken(context.Context, *GetUserByEdnAppLoginTokenRequest) (*GetUserDetailResponse, error)
	GetDivisionDetail(context.Context, *GetDivisionDetailRequest) (*GetDivisionDetailResponse, error)
	GetDivisionDefaultByCustomerType(context.Context, *GetDivisionDefaultByCustomerTypeRequest) (*GetDivisionDefaultByCustomerTypeResponse, error)
}

// UnimplementedAccountServiceServer should be embedded to have forward compatible implementations.
type UnimplementedAccountServiceServer struct {
}

func (UnimplementedAccountServiceServer) GetUserList(context.Context, *GetUserListRequest) (*GetUserListResponse, error) {
	return nil, status.Errorf(codes.Unimplemented, "method GetUserList not implemented")
}
func (UnimplementedAccountServiceServer) GetUserDetail(context.Context, *GetUserDetailRequest) (*GetUserDetailResponse, error) {
	return nil, status.Errorf(codes.Unimplemented, "method GetUserDetail not implemented")
}
func (UnimplementedAccountServiceServer) GetUserEmailAuth(context.Context, *GetUserEmailAuthRequest) (*GetUserEmailAuthResponse, error) {
	return nil, status.Errorf(codes.Unimplemented, "method GetUserEmailAuth not implemented")
}
func (UnimplementedAccountServiceServer) GetUserRoleByUserId(context.Context, *GetUserRoleByUserIdRequest) (*GetUserRoleByUserIdResponse, error) {
	return nil, status.Errorf(codes.Unimplemented, "method GetUserRoleByUserId not implemented")
}
func (UnimplementedAccountServiceServer) UpdateUserSalesAppToken(context.Context, *UpdateUserSalesAppTokenRequest) (*GetUserDetailResponse, error) {
	return nil, status.Errorf(codes.Unimplemented, "method UpdateUserSalesAppToken not implemented")
}
func (UnimplementedAccountServiceServer) GetUserBySalesAppLoginToken(context.Context, *GetUserBySalesAppLoginTokenRequest) (*GetUserDetailResponse, error) {
	return nil, status.Errorf(codes.Unimplemented, "method GetUserBySalesAppLoginToken not implemented")
}
func (UnimplementedAccountServiceServer) GetRoleDetail(context.Context, *GetRoleDetailRequest) (*GetRoleDetailResponse, error) {
	return nil, status.Errorf(codes.Unimplemented, "method GetRoleDetail not implemented")
}
func (UnimplementedAccountServiceServer) UpdateUserEdnAppToken(context.Context, *UpdateUserEdnAppTokenRequest) (*GetUserDetailResponse, error) {
	return nil, status.Errorf(codes.Unimplemented, "method UpdateUserEdnAppToken not implemented")
}
func (UnimplementedAccountServiceServer) UpdateUserPurchaserAppToken(context.Context, *UpdateUserPurchaserAppTokenRequest) (*GetUserDetailResponse, error) {
	return nil, status.Errorf(codes.Unimplemented, "method UpdateUserPurchaserAppToken not implemented")
}
func (UnimplementedAccountServiceServer) GetUserByEdnAppLoginToken(context.Context, *GetUserByEdnAppLoginTokenRequest) (*GetUserDetailResponse, error) {
	return nil, status.Errorf(codes.Unimplemented, "method GetUserByEdnAppLoginToken not implemented")
}
func (UnimplementedAccountServiceServer) GetDivisionDetail(context.Context, *GetDivisionDetailRequest) (*GetDivisionDetailResponse, error) {
	return nil, status.Errorf(codes.Unimplemented, "method GetDivisionDetail not implemented")
}
func (UnimplementedAccountServiceServer) GetDivisionDefaultByCustomerType(context.Context, *GetDivisionDefaultByCustomerTypeRequest) (*GetDivisionDefaultByCustomerTypeResponse, error) {
	return nil, status.Errorf(codes.Unimplemented, "method GetDivisionDefaultByCustomerType not implemented")
}

// UnsafeAccountServiceServer may be embedded to opt out of forward compatibility for this service.
// Use of this interface is not recommended, as added methods to AccountServiceServer will
// result in compilation errors.
type UnsafeAccountServiceServer interface {
	mustEmbedUnimplementedAccountServiceServer()
}

func RegisterAccountServiceServer(s grpc.ServiceRegistrar, srv AccountServiceServer) {
	s.RegisterService(&AccountService_ServiceDesc, srv)
}

func _AccountService_GetUserList_Handler(srv interface{}, ctx context.Context, dec func(interface{}) error, interceptor grpc.UnaryServerInterceptor) (interface{}, error) {
	in := new(GetUserListRequest)
	if err := dec(in); err != nil {
		return nil, err
	}
	if interceptor == nil {
		return srv.(AccountServiceServer).GetUserList(ctx, in)
	}
	info := &grpc.UnaryServerInfo{
		Server:     srv,
		FullMethod: "/proto.account_service.AccountService/GetUserList",
	}
	handler := func(ctx context.Context, req interface{}) (interface{}, error) {
		return srv.(AccountServiceServer).GetUserList(ctx, req.(*GetUserListRequest))
	}
	return interceptor(ctx, in, info, handler)
}

func _AccountService_GetUserDetail_Handler(srv interface{}, ctx context.Context, dec func(interface{}) error, interceptor grpc.UnaryServerInterceptor) (interface{}, error) {
	in := new(GetUserDetailRequest)
	if err := dec(in); err != nil {
		return nil, err
	}
	if interceptor == nil {
		return srv.(AccountServiceServer).GetUserDetail(ctx, in)
	}
	info := &grpc.UnaryServerInfo{
		Server:     srv,
		FullMethod: "/proto.account_service.AccountService/GetUserDetail",
	}
	handler := func(ctx context.Context, req interface{}) (interface{}, error) {
		return srv.(AccountServiceServer).GetUserDetail(ctx, req.(*GetUserDetailRequest))
	}
	return interceptor(ctx, in, info, handler)
}

func _AccountService_GetUserEmailAuth_Handler(srv interface{}, ctx context.Context, dec func(interface{}) error, interceptor grpc.UnaryServerInterceptor) (interface{}, error) {
	in := new(GetUserEmailAuthRequest)
	if err := dec(in); err != nil {
		return nil, err
	}
	if interceptor == nil {
		return srv.(AccountServiceServer).GetUserEmailAuth(ctx, in)
	}
	info := &grpc.UnaryServerInfo{
		Server:     srv,
		FullMethod: "/proto.account_service.AccountService/GetUserEmailAuth",
	}
	handler := func(ctx context.Context, req interface{}) (interface{}, error) {
		return srv.(AccountServiceServer).GetUserEmailAuth(ctx, req.(*GetUserEmailAuthRequest))
	}
	return interceptor(ctx, in, info, handler)
}

func _AccountService_GetUserRoleByUserId_Handler(srv interface{}, ctx context.Context, dec func(interface{}) error, interceptor grpc.UnaryServerInterceptor) (interface{}, error) {
	in := new(GetUserRoleByUserIdRequest)
	if err := dec(in); err != nil {
		return nil, err
	}
	if interceptor == nil {
		return srv.(AccountServiceServer).GetUserRoleByUserId(ctx, in)
	}
	info := &grpc.UnaryServerInfo{
		Server:     srv,
		FullMethod: "/proto.account_service.AccountService/GetUserRoleByUserId",
	}
	handler := func(ctx context.Context, req interface{}) (interface{}, error) {
		return srv.(AccountServiceServer).GetUserRoleByUserId(ctx, req.(*GetUserRoleByUserIdRequest))
	}
	return interceptor(ctx, in, info, handler)
}

func _AccountService_UpdateUserSalesAppToken_Handler(srv interface{}, ctx context.Context, dec func(interface{}) error, interceptor grpc.UnaryServerInterceptor) (interface{}, error) {
	in := new(UpdateUserSalesAppTokenRequest)
	if err := dec(in); err != nil {
		return nil, err
	}
	if interceptor == nil {
		return srv.(AccountServiceServer).UpdateUserSalesAppToken(ctx, in)
	}
	info := &grpc.UnaryServerInfo{
		Server:     srv,
		FullMethod: "/proto.account_service.AccountService/UpdateUserSalesAppToken",
	}
	handler := func(ctx context.Context, req interface{}) (interface{}, error) {
		return srv.(AccountServiceServer).UpdateUserSalesAppToken(ctx, req.(*UpdateUserSalesAppTokenRequest))
	}
	return interceptor(ctx, in, info, handler)
}

func _AccountService_GetUserBySalesAppLoginToken_Handler(srv interface{}, ctx context.Context, dec func(interface{}) error, interceptor grpc.UnaryServerInterceptor) (interface{}, error) {
	in := new(GetUserBySalesAppLoginTokenRequest)
	if err := dec(in); err != nil {
		return nil, err
	}
	if interceptor == nil {
		return srv.(AccountServiceServer).GetUserBySalesAppLoginToken(ctx, in)
	}
	info := &grpc.UnaryServerInfo{
		Server:     srv,
		FullMethod: "/proto.account_service.AccountService/GetUserBySalesAppLoginToken",
	}
	handler := func(ctx context.Context, req interface{}) (interface{}, error) {
		return srv.(AccountServiceServer).GetUserBySalesAppLoginToken(ctx, req.(*GetUserBySalesAppLoginTokenRequest))
	}
	return interceptor(ctx, in, info, handler)
}

func _AccountService_GetRoleDetail_Handler(srv interface{}, ctx context.Context, dec func(interface{}) error, interceptor grpc.UnaryServerInterceptor) (interface{}, error) {
	in := new(GetRoleDetailRequest)
	if err := dec(in); err != nil {
		return nil, err
	}
	if interceptor == nil {
		return srv.(AccountServiceServer).GetRoleDetail(ctx, in)
	}
	info := &grpc.UnaryServerInfo{
		Server:     srv,
		FullMethod: "/proto.account_service.AccountService/GetRoleDetail",
	}
	handler := func(ctx context.Context, req interface{}) (interface{}, error) {
		return srv.(AccountServiceServer).GetRoleDetail(ctx, req.(*GetRoleDetailRequest))
	}
	return interceptor(ctx, in, info, handler)
}

func _AccountService_UpdateUserEdnAppToken_Handler(srv interface{}, ctx context.Context, dec func(interface{}) error, interceptor grpc.UnaryServerInterceptor) (interface{}, error) {
	in := new(UpdateUserEdnAppTokenRequest)
	if err := dec(in); err != nil {
		return nil, err
	}
	if interceptor == nil {
		return srv.(AccountServiceServer).UpdateUserEdnAppToken(ctx, in)
	}
	info := &grpc.UnaryServerInfo{
		Server:     srv,
		FullMethod: "/proto.account_service.AccountService/UpdateUserEdnAppToken",
	}
	handler := func(ctx context.Context, req interface{}) (interface{}, error) {
		return srv.(AccountServiceServer).UpdateUserEdnAppToken(ctx, req.(*UpdateUserEdnAppTokenRequest))
	}
	return interceptor(ctx, in, info, handler)
}

func _AccountService_UpdateUserPurchaserAppToken_Handler(srv interface{}, ctx context.Context, dec func(interface{}) error, interceptor grpc.UnaryServerInterceptor) (interface{}, error) {
	in := new(UpdateUserPurchaserAppTokenRequest)
	if err := dec(in); err != nil {
		return nil, err
	}
	if interceptor == nil {
		return srv.(AccountServiceServer).UpdateUserPurchaserAppToken(ctx, in)
	}
	info := &grpc.UnaryServerInfo{
		Server:     srv,
		FullMethod: "/proto.account_service.AccountService/UpdateUserPurchaserAppToken",
	}
	handler := func(ctx context.Context, req interface{}) (interface{}, error) {
		return srv.(AccountServiceServer).UpdateUserPurchaserAppToken(ctx, req.(*UpdateUserPurchaserAppTokenRequest))
	}
	return interceptor(ctx, in, info, handler)
}

func _AccountService_GetUserByEdnAppLoginToken_Handler(srv interface{}, ctx context.Context, dec func(interface{}) error, interceptor grpc.UnaryServerInterceptor) (interface{}, error) {
	in := new(GetUserByEdnAppLoginTokenRequest)
	if err := dec(in); err != nil {
		return nil, err
	}
	if interceptor == nil {
		return srv.(AccountServiceServer).GetUserByEdnAppLoginToken(ctx, in)
	}
	info := &grpc.UnaryServerInfo{
		Server:     srv,
		FullMethod: "/proto.account_service.AccountService/GetUserByEdnAppLoginToken",
	}
	handler := func(ctx context.Context, req interface{}) (interface{}, error) {
		return srv.(AccountServiceServer).GetUserByEdnAppLoginToken(ctx, req.(*GetUserByEdnAppLoginTokenRequest))
	}
	return interceptor(ctx, in, info, handler)
}

func _AccountService_GetDivisionDetail_Handler(srv interface{}, ctx context.Context, dec func(interface{}) error, interceptor grpc.UnaryServerInterceptor) (interface{}, error) {
	in := new(GetDivisionDetailRequest)
	if err := dec(in); err != nil {
		return nil, err
	}
	if interceptor == nil {
		return srv.(AccountServiceServer).GetDivisionDetail(ctx, in)
	}
	info := &grpc.UnaryServerInfo{
		Server:     srv,
		FullMethod: "/proto.account_service.AccountService/GetDivisionDetail",
	}
	handler := func(ctx context.Context, req interface{}) (interface{}, error) {
		return srv.(AccountServiceServer).GetDivisionDetail(ctx, req.(*GetDivisionDetailRequest))
	}
	return interceptor(ctx, in, info, handler)
}

func _AccountService_GetDivisionDefaultByCustomerType_Handler(srv interface{}, ctx context.Context, dec func(interface{}) error, interceptor grpc.UnaryServerInterceptor) (interface{}, error) {
	in := new(GetDivisionDefaultByCustomerTypeRequest)
	if err := dec(in); err != nil {
		return nil, err
	}
	if interceptor == nil {
		return srv.(AccountServiceServer).GetDivisionDefaultByCustomerType(ctx, in)
	}
	info := &grpc.UnaryServerInfo{
		Server:     srv,
		FullMethod: "/proto.account_service.AccountService/GetDivisionDefaultByCustomerType",
	}
	handler := func(ctx context.Context, req interface{}) (interface{}, error) {
		return srv.(AccountServiceServer).GetDivisionDefaultByCustomerType(ctx, req.(*GetDivisionDefaultByCustomerTypeRequest))
	}
	return interceptor(ctx, in, info, handler)
}

// AccountService_ServiceDesc is the grpc.ServiceDesc for AccountService service.
// It's only intended for direct use with grpc.RegisterService,
// and not to be introspected or modified (even as a copy)
var AccountService_ServiceDesc = grpc.ServiceDesc{
	ServiceName: "proto.account_service.AccountService",
	HandlerType: (*AccountServiceServer)(nil),
	Methods: []grpc.MethodDesc{
		{
			MethodName: "GetUserList",
			Handler:    _AccountService_GetUserList_Handler,
		},
		{
			MethodName: "GetUserDetail",
			Handler:    _AccountService_GetUserDetail_Handler,
		},
		{
			MethodName: "GetUserEmailAuth",
			Handler:    _AccountService_GetUserEmailAuth_Handler,
		},
		{
			MethodName: "GetUserRoleByUserId",
			Handler:    _AccountService_GetUserRoleByUserId_Handler,
		},
		{
			MethodName: "UpdateUserSalesAppToken",
			Handler:    _AccountService_UpdateUserSalesAppToken_Handler,
		},
		{
			MethodName: "GetUserBySalesAppLoginToken",
			Handler:    _AccountService_GetUserBySalesAppLoginToken_Handler,
		},
		{
			MethodName: "GetRoleDetail",
			Handler:    _AccountService_GetRoleDetail_Handler,
		},
		{
			MethodName: "UpdateUserEdnAppToken",
			Handler:    _AccountService_UpdateUserEdnAppToken_Handler,
		},
		{
			MethodName: "UpdateUserPurchaserAppToken",
			Handler:    _AccountService_UpdateUserPurchaserAppToken_Handler,
		},
		{
			MethodName: "GetUserByEdnAppLoginToken",
			Handler:    _AccountService_GetUserByEdnAppLoginToken_Handler,
		},
		{
			MethodName: "GetDivisionDetail",
			Handler:    _AccountService_GetDivisionDetail_Handler,
		},
		{
			MethodName: "GetDivisionDefaultByCustomerType",
			Handler:    _AccountService_GetDivisionDefaultByCustomerType_Handler,
		},
	},
	Streams:  []grpc.StreamDesc{},
	Metadata: "proto/account_service/account_service.proto",
}