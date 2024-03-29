// Code generated by protoc-gen-go-grpc. DO NOT EDIT.
// versions:
// - protoc-gen-go-grpc v1.2.0
// - protoc             v3.12.4
// source: protos/token.proto

package token

import (
	context "context"
	wrappers "github.com/golang/protobuf/ptypes/wrappers"
	grpc "google.golang.org/grpc"
	codes "google.golang.org/grpc/codes"
	status "google.golang.org/grpc/status"
)

// This is a compile-time assertion to ensure that this generated file
// is compatible with the grpc package it is being compiled against.
// Requires gRPC-Go v1.32.0 or later.
const _ = grpc.SupportPackageIsVersion7

// TokenServiceClient is the client API for TokenService service.
//
// For semantics around ctx use and closing/ending streaming RPCs, please refer to https://pkg.go.dev/google.golang.org/grpc/?tab=doc#ClientConn.NewStream.
type TokenServiceClient interface {
	NewAuthCodeToken(ctx context.Context, in *TokenRequest, opts ...grpc.CallOption) (*AuthCodeToken, error)
	NewAuthCodeGrant(ctx context.Context, in *TokenRequest, opts ...grpc.CallOption) (*wrappers.StringValue, error)
	VerifyAuthCodeToken(ctx context.Context, in *TokenRequest, opts ...grpc.CallOption) (*wrappers.BoolValue, error)
	NewAuthCodeRefreshToken(ctx context.Context, in *RefreshTokenRequest, opts ...grpc.CallOption) (*AuthCodeToken, error)
	AuthCodeRefreshTokenExists(ctx context.Context, in *RefreshTokenRequest, opts ...grpc.CallOption) (*wrappers.BoolValue, error)
	AddUserIdAuthCodeFlow(ctx context.Context, in *User, opts ...grpc.CallOption) (*wrappers.BoolValue, error)
	GetUserIdfromAccesstoken(ctx context.Context, in *User, opts ...grpc.CallOption) (*wrappers.StringValue, error)
	CreateToken(ctx context.Context, in *User, opts ...grpc.CallOption) (*AuthCodeToken, error)
	VerifyToken(ctx context.Context, in *TokenRequest, opts ...grpc.CallOption) (*AuthCodeToken, error)
	RefreshToken(ctx context.Context, in *TokenRequest, opts ...grpc.CallOption) (*AuthCodeToken, error)
	RevokeToken(ctx context.Context, in *TokenRequest, opts ...grpc.CallOption) (*AuthCodeToken, error)
	RevokeAll(ctx context.Context, in *User, opts ...grpc.CallOption) (*AuthCodeToken, error)
	CreateResetPasswordToken(ctx context.Context, in *User, opts ...grpc.CallOption) (*AuthCodeToken, error)
}

type tokenServiceClient struct {
	cc grpc.ClientConnInterface
}

func NewTokenServiceClient(cc grpc.ClientConnInterface) TokenServiceClient {
	return &tokenServiceClient{cc}
}

func (c *tokenServiceClient) NewAuthCodeToken(ctx context.Context, in *TokenRequest, opts ...grpc.CallOption) (*AuthCodeToken, error) {
	out := new(AuthCodeToken)
	err := c.cc.Invoke(ctx, "/token.TokenService/newAuthCodeToken", in, out, opts...)
	if err != nil {
		return nil, err
	}
	return out, nil
}

func (c *tokenServiceClient) NewAuthCodeGrant(ctx context.Context, in *TokenRequest, opts ...grpc.CallOption) (*wrappers.StringValue, error) {
	out := new(wrappers.StringValue)
	err := c.cc.Invoke(ctx, "/token.TokenService/newAuthCodeGrant", in, out, opts...)
	if err != nil {
		return nil, err
	}
	return out, nil
}

func (c *tokenServiceClient) VerifyAuthCodeToken(ctx context.Context, in *TokenRequest, opts ...grpc.CallOption) (*wrappers.BoolValue, error) {
	out := new(wrappers.BoolValue)
	err := c.cc.Invoke(ctx, "/token.TokenService/verifyAuthCodeToken", in, out, opts...)
	if err != nil {
		return nil, err
	}
	return out, nil
}

func (c *tokenServiceClient) NewAuthCodeRefreshToken(ctx context.Context, in *RefreshTokenRequest, opts ...grpc.CallOption) (*AuthCodeToken, error) {
	out := new(AuthCodeToken)
	err := c.cc.Invoke(ctx, "/token.TokenService/newAuthCodeRefreshToken", in, out, opts...)
	if err != nil {
		return nil, err
	}
	return out, nil
}

func (c *tokenServiceClient) AuthCodeRefreshTokenExists(ctx context.Context, in *RefreshTokenRequest, opts ...grpc.CallOption) (*wrappers.BoolValue, error) {
	out := new(wrappers.BoolValue)
	err := c.cc.Invoke(ctx, "/token.TokenService/authCodeRefreshTokenExists", in, out, opts...)
	if err != nil {
		return nil, err
	}
	return out, nil
}

func (c *tokenServiceClient) AddUserIdAuthCodeFlow(ctx context.Context, in *User, opts ...grpc.CallOption) (*wrappers.BoolValue, error) {
	out := new(wrappers.BoolValue)
	err := c.cc.Invoke(ctx, "/token.TokenService/addUserIdAuthCodeFlow", in, out, opts...)
	if err != nil {
		return nil, err
	}
	return out, nil
}

func (c *tokenServiceClient) GetUserIdfromAccesstoken(ctx context.Context, in *User, opts ...grpc.CallOption) (*wrappers.StringValue, error) {
	out := new(wrappers.StringValue)
	err := c.cc.Invoke(ctx, "/token.TokenService/getUserIdfromAccesstoken", in, out, opts...)
	if err != nil {
		return nil, err
	}
	return out, nil
}

func (c *tokenServiceClient) CreateToken(ctx context.Context, in *User, opts ...grpc.CallOption) (*AuthCodeToken, error) {
	out := new(AuthCodeToken)
	err := c.cc.Invoke(ctx, "/token.TokenService/createToken", in, out, opts...)
	if err != nil {
		return nil, err
	}
	return out, nil
}

func (c *tokenServiceClient) VerifyToken(ctx context.Context, in *TokenRequest, opts ...grpc.CallOption) (*AuthCodeToken, error) {
	out := new(AuthCodeToken)
	err := c.cc.Invoke(ctx, "/token.TokenService/verifyToken", in, out, opts...)
	if err != nil {
		return nil, err
	}
	return out, nil
}

func (c *tokenServiceClient) RefreshToken(ctx context.Context, in *TokenRequest, opts ...grpc.CallOption) (*AuthCodeToken, error) {
	out := new(AuthCodeToken)
	err := c.cc.Invoke(ctx, "/token.TokenService/refreshToken", in, out, opts...)
	if err != nil {
		return nil, err
	}
	return out, nil
}

func (c *tokenServiceClient) RevokeToken(ctx context.Context, in *TokenRequest, opts ...grpc.CallOption) (*AuthCodeToken, error) {
	out := new(AuthCodeToken)
	err := c.cc.Invoke(ctx, "/token.TokenService/revokeToken", in, out, opts...)
	if err != nil {
		return nil, err
	}
	return out, nil
}

func (c *tokenServiceClient) RevokeAll(ctx context.Context, in *User, opts ...grpc.CallOption) (*AuthCodeToken, error) {
	out := new(AuthCodeToken)
	err := c.cc.Invoke(ctx, "/token.TokenService/revokeAll", in, out, opts...)
	if err != nil {
		return nil, err
	}
	return out, nil
}

func (c *tokenServiceClient) CreateResetPasswordToken(ctx context.Context, in *User, opts ...grpc.CallOption) (*AuthCodeToken, error) {
	out := new(AuthCodeToken)
	err := c.cc.Invoke(ctx, "/token.TokenService/createResetPasswordToken", in, out, opts...)
	if err != nil {
		return nil, err
	}
	return out, nil
}

// TokenServiceServer is the server API for TokenService service.
// All implementations must embed UnimplementedTokenServiceServer
// for forward compatibility
type TokenServiceServer interface {
	NewAuthCodeToken(context.Context, *TokenRequest) (*AuthCodeToken, error)
	NewAuthCodeGrant(context.Context, *TokenRequest) (*wrappers.StringValue, error)
	VerifyAuthCodeToken(context.Context, *TokenRequest) (*wrappers.BoolValue, error)
	NewAuthCodeRefreshToken(context.Context, *RefreshTokenRequest) (*AuthCodeToken, error)
	AuthCodeRefreshTokenExists(context.Context, *RefreshTokenRequest) (*wrappers.BoolValue, error)
	AddUserIdAuthCodeFlow(context.Context, *User) (*wrappers.BoolValue, error)
	GetUserIdfromAccesstoken(context.Context, *User) (*wrappers.StringValue, error)
	CreateToken(context.Context, *User) (*AuthCodeToken, error)
	VerifyToken(context.Context, *TokenRequest) (*AuthCodeToken, error)
	RefreshToken(context.Context, *TokenRequest) (*AuthCodeToken, error)
	RevokeToken(context.Context, *TokenRequest) (*AuthCodeToken, error)
	RevokeAll(context.Context, *User) (*AuthCodeToken, error)
	CreateResetPasswordToken(context.Context, *User) (*AuthCodeToken, error)
	mustEmbedUnimplementedTokenServiceServer()
}

// UnimplementedTokenServiceServer must be embedded to have forward compatible implementations.
type UnimplementedTokenServiceServer struct {
}

func (UnimplementedTokenServiceServer) NewAuthCodeToken(context.Context, *TokenRequest) (*AuthCodeToken, error) {
	return nil, status.Errorf(codes.Unimplemented, "method NewAuthCodeToken not implemented")
}
func (UnimplementedTokenServiceServer) NewAuthCodeGrant(context.Context, *TokenRequest) (*wrappers.StringValue, error) {
	return nil, status.Errorf(codes.Unimplemented, "method NewAuthCodeGrant not implemented")
}
func (UnimplementedTokenServiceServer) VerifyAuthCodeToken(context.Context, *TokenRequest) (*wrappers.BoolValue, error) {
	return nil, status.Errorf(codes.Unimplemented, "method VerifyAuthCodeToken not implemented")
}
func (UnimplementedTokenServiceServer) NewAuthCodeRefreshToken(context.Context, *RefreshTokenRequest) (*AuthCodeToken, error) {
	return nil, status.Errorf(codes.Unimplemented, "method NewAuthCodeRefreshToken not implemented")
}
func (UnimplementedTokenServiceServer) AuthCodeRefreshTokenExists(context.Context, *RefreshTokenRequest) (*wrappers.BoolValue, error) {
	return nil, status.Errorf(codes.Unimplemented, "method AuthCodeRefreshTokenExists not implemented")
}
func (UnimplementedTokenServiceServer) AddUserIdAuthCodeFlow(context.Context, *User) (*wrappers.BoolValue, error) {
	return nil, status.Errorf(codes.Unimplemented, "method AddUserIdAuthCodeFlow not implemented")
}
func (UnimplementedTokenServiceServer) GetUserIdfromAccesstoken(context.Context, *User) (*wrappers.StringValue, error) {
	return nil, status.Errorf(codes.Unimplemented, "method GetUserIdfromAccesstoken not implemented")
}
func (UnimplementedTokenServiceServer) CreateToken(context.Context, *User) (*AuthCodeToken, error) {
	return nil, status.Errorf(codes.Unimplemented, "method CreateToken not implemented")
}
func (UnimplementedTokenServiceServer) VerifyToken(context.Context, *TokenRequest) (*AuthCodeToken, error) {
	return nil, status.Errorf(codes.Unimplemented, "method VerifyToken not implemented")
}
func (UnimplementedTokenServiceServer) RefreshToken(context.Context, *TokenRequest) (*AuthCodeToken, error) {
	return nil, status.Errorf(codes.Unimplemented, "method RefreshToken not implemented")
}
func (UnimplementedTokenServiceServer) RevokeToken(context.Context, *TokenRequest) (*AuthCodeToken, error) {
	return nil, status.Errorf(codes.Unimplemented, "method RevokeToken not implemented")
}
func (UnimplementedTokenServiceServer) RevokeAll(context.Context, *User) (*AuthCodeToken, error) {
	return nil, status.Errorf(codes.Unimplemented, "method RevokeAll not implemented")
}
func (UnimplementedTokenServiceServer) CreateResetPasswordToken(context.Context, *User) (*AuthCodeToken, error) {
	return nil, status.Errorf(codes.Unimplemented, "method CreateResetPasswordToken not implemented")
}
func (UnimplementedTokenServiceServer) mustEmbedUnimplementedTokenServiceServer() {}

// UnsafeTokenServiceServer may be embedded to opt out of forward compatibility for this service.
// Use of this interface is not recommended, as added methods to TokenServiceServer will
// result in compilation errors.
type UnsafeTokenServiceServer interface {
	mustEmbedUnimplementedTokenServiceServer()
}

func RegisterTokenServiceServer(s grpc.ServiceRegistrar, srv TokenServiceServer) {
	s.RegisterService(&TokenService_ServiceDesc, srv)
}

func _TokenService_NewAuthCodeToken_Handler(srv interface{}, ctx context.Context, dec func(interface{}) error, interceptor grpc.UnaryServerInterceptor) (interface{}, error) {
	in := new(TokenRequest)
	if err := dec(in); err != nil {
		return nil, err
	}
	if interceptor == nil {
		return srv.(TokenServiceServer).NewAuthCodeToken(ctx, in)
	}
	info := &grpc.UnaryServerInfo{
		Server:     srv,
		FullMethod: "/token.TokenService/newAuthCodeToken",
	}
	handler := func(ctx context.Context, req interface{}) (interface{}, error) {
		return srv.(TokenServiceServer).NewAuthCodeToken(ctx, req.(*TokenRequest))
	}
	return interceptor(ctx, in, info, handler)
}

func _TokenService_NewAuthCodeGrant_Handler(srv interface{}, ctx context.Context, dec func(interface{}) error, interceptor grpc.UnaryServerInterceptor) (interface{}, error) {
	in := new(TokenRequest)
	if err := dec(in); err != nil {
		return nil, err
	}
	if interceptor == nil {
		return srv.(TokenServiceServer).NewAuthCodeGrant(ctx, in)
	}
	info := &grpc.UnaryServerInfo{
		Server:     srv,
		FullMethod: "/token.TokenService/newAuthCodeGrant",
	}
	handler := func(ctx context.Context, req interface{}) (interface{}, error) {
		return srv.(TokenServiceServer).NewAuthCodeGrant(ctx, req.(*TokenRequest))
	}
	return interceptor(ctx, in, info, handler)
}

func _TokenService_VerifyAuthCodeToken_Handler(srv interface{}, ctx context.Context, dec func(interface{}) error, interceptor grpc.UnaryServerInterceptor) (interface{}, error) {
	in := new(TokenRequest)
	if err := dec(in); err != nil {
		return nil, err
	}
	if interceptor == nil {
		return srv.(TokenServiceServer).VerifyAuthCodeToken(ctx, in)
	}
	info := &grpc.UnaryServerInfo{
		Server:     srv,
		FullMethod: "/token.TokenService/verifyAuthCodeToken",
	}
	handler := func(ctx context.Context, req interface{}) (interface{}, error) {
		return srv.(TokenServiceServer).VerifyAuthCodeToken(ctx, req.(*TokenRequest))
	}
	return interceptor(ctx, in, info, handler)
}

func _TokenService_NewAuthCodeRefreshToken_Handler(srv interface{}, ctx context.Context, dec func(interface{}) error, interceptor grpc.UnaryServerInterceptor) (interface{}, error) {
	in := new(RefreshTokenRequest)
	if err := dec(in); err != nil {
		return nil, err
	}
	if interceptor == nil {
		return srv.(TokenServiceServer).NewAuthCodeRefreshToken(ctx, in)
	}
	info := &grpc.UnaryServerInfo{
		Server:     srv,
		FullMethod: "/token.TokenService/newAuthCodeRefreshToken",
	}
	handler := func(ctx context.Context, req interface{}) (interface{}, error) {
		return srv.(TokenServiceServer).NewAuthCodeRefreshToken(ctx, req.(*RefreshTokenRequest))
	}
	return interceptor(ctx, in, info, handler)
}

func _TokenService_AuthCodeRefreshTokenExists_Handler(srv interface{}, ctx context.Context, dec func(interface{}) error, interceptor grpc.UnaryServerInterceptor) (interface{}, error) {
	in := new(RefreshTokenRequest)
	if err := dec(in); err != nil {
		return nil, err
	}
	if interceptor == nil {
		return srv.(TokenServiceServer).AuthCodeRefreshTokenExists(ctx, in)
	}
	info := &grpc.UnaryServerInfo{
		Server:     srv,
		FullMethod: "/token.TokenService/authCodeRefreshTokenExists",
	}
	handler := func(ctx context.Context, req interface{}) (interface{}, error) {
		return srv.(TokenServiceServer).AuthCodeRefreshTokenExists(ctx, req.(*RefreshTokenRequest))
	}
	return interceptor(ctx, in, info, handler)
}

func _TokenService_AddUserIdAuthCodeFlow_Handler(srv interface{}, ctx context.Context, dec func(interface{}) error, interceptor grpc.UnaryServerInterceptor) (interface{}, error) {
	in := new(User)
	if err := dec(in); err != nil {
		return nil, err
	}
	if interceptor == nil {
		return srv.(TokenServiceServer).AddUserIdAuthCodeFlow(ctx, in)
	}
	info := &grpc.UnaryServerInfo{
		Server:     srv,
		FullMethod: "/token.TokenService/addUserIdAuthCodeFlow",
	}
	handler := func(ctx context.Context, req interface{}) (interface{}, error) {
		return srv.(TokenServiceServer).AddUserIdAuthCodeFlow(ctx, req.(*User))
	}
	return interceptor(ctx, in, info, handler)
}

func _TokenService_GetUserIdfromAccesstoken_Handler(srv interface{}, ctx context.Context, dec func(interface{}) error, interceptor grpc.UnaryServerInterceptor) (interface{}, error) {
	in := new(User)
	if err := dec(in); err != nil {
		return nil, err
	}
	if interceptor == nil {
		return srv.(TokenServiceServer).GetUserIdfromAccesstoken(ctx, in)
	}
	info := &grpc.UnaryServerInfo{
		Server:     srv,
		FullMethod: "/token.TokenService/getUserIdfromAccesstoken",
	}
	handler := func(ctx context.Context, req interface{}) (interface{}, error) {
		return srv.(TokenServiceServer).GetUserIdfromAccesstoken(ctx, req.(*User))
	}
	return interceptor(ctx, in, info, handler)
}

func _TokenService_CreateToken_Handler(srv interface{}, ctx context.Context, dec func(interface{}) error, interceptor grpc.UnaryServerInterceptor) (interface{}, error) {
	in := new(User)
	if err := dec(in); err != nil {
		return nil, err
	}
	if interceptor == nil {
		return srv.(TokenServiceServer).CreateToken(ctx, in)
	}
	info := &grpc.UnaryServerInfo{
		Server:     srv,
		FullMethod: "/token.TokenService/createToken",
	}
	handler := func(ctx context.Context, req interface{}) (interface{}, error) {
		return srv.(TokenServiceServer).CreateToken(ctx, req.(*User))
	}
	return interceptor(ctx, in, info, handler)
}

func _TokenService_VerifyToken_Handler(srv interface{}, ctx context.Context, dec func(interface{}) error, interceptor grpc.UnaryServerInterceptor) (interface{}, error) {
	in := new(TokenRequest)
	if err := dec(in); err != nil {
		return nil, err
	}
	if interceptor == nil {
		return srv.(TokenServiceServer).VerifyToken(ctx, in)
	}
	info := &grpc.UnaryServerInfo{
		Server:     srv,
		FullMethod: "/token.TokenService/verifyToken",
	}
	handler := func(ctx context.Context, req interface{}) (interface{}, error) {
		return srv.(TokenServiceServer).VerifyToken(ctx, req.(*TokenRequest))
	}
	return interceptor(ctx, in, info, handler)
}

func _TokenService_RefreshToken_Handler(srv interface{}, ctx context.Context, dec func(interface{}) error, interceptor grpc.UnaryServerInterceptor) (interface{}, error) {
	in := new(TokenRequest)
	if err := dec(in); err != nil {
		return nil, err
	}
	if interceptor == nil {
		return srv.(TokenServiceServer).RefreshToken(ctx, in)
	}
	info := &grpc.UnaryServerInfo{
		Server:     srv,
		FullMethod: "/token.TokenService/refreshToken",
	}
	handler := func(ctx context.Context, req interface{}) (interface{}, error) {
		return srv.(TokenServiceServer).RefreshToken(ctx, req.(*TokenRequest))
	}
	return interceptor(ctx, in, info, handler)
}

func _TokenService_RevokeToken_Handler(srv interface{}, ctx context.Context, dec func(interface{}) error, interceptor grpc.UnaryServerInterceptor) (interface{}, error) {
	in := new(TokenRequest)
	if err := dec(in); err != nil {
		return nil, err
	}
	if interceptor == nil {
		return srv.(TokenServiceServer).RevokeToken(ctx, in)
	}
	info := &grpc.UnaryServerInfo{
		Server:     srv,
		FullMethod: "/token.TokenService/revokeToken",
	}
	handler := func(ctx context.Context, req interface{}) (interface{}, error) {
		return srv.(TokenServiceServer).RevokeToken(ctx, req.(*TokenRequest))
	}
	return interceptor(ctx, in, info, handler)
}

func _TokenService_RevokeAll_Handler(srv interface{}, ctx context.Context, dec func(interface{}) error, interceptor grpc.UnaryServerInterceptor) (interface{}, error) {
	in := new(User)
	if err := dec(in); err != nil {
		return nil, err
	}
	if interceptor == nil {
		return srv.(TokenServiceServer).RevokeAll(ctx, in)
	}
	info := &grpc.UnaryServerInfo{
		Server:     srv,
		FullMethod: "/token.TokenService/revokeAll",
	}
	handler := func(ctx context.Context, req interface{}) (interface{}, error) {
		return srv.(TokenServiceServer).RevokeAll(ctx, req.(*User))
	}
	return interceptor(ctx, in, info, handler)
}

func _TokenService_CreateResetPasswordToken_Handler(srv interface{}, ctx context.Context, dec func(interface{}) error, interceptor grpc.UnaryServerInterceptor) (interface{}, error) {
	in := new(User)
	if err := dec(in); err != nil {
		return nil, err
	}
	if interceptor == nil {
		return srv.(TokenServiceServer).CreateResetPasswordToken(ctx, in)
	}
	info := &grpc.UnaryServerInfo{
		Server:     srv,
		FullMethod: "/token.TokenService/createResetPasswordToken",
	}
	handler := func(ctx context.Context, req interface{}) (interface{}, error) {
		return srv.(TokenServiceServer).CreateResetPasswordToken(ctx, req.(*User))
	}
	return interceptor(ctx, in, info, handler)
}

// TokenService_ServiceDesc is the grpc.ServiceDesc for TokenService service.
// It's only intended for direct use with grpc.RegisterService,
// and not to be introspected or modified (even as a copy)
var TokenService_ServiceDesc = grpc.ServiceDesc{
	ServiceName: "token.TokenService",
	HandlerType: (*TokenServiceServer)(nil),
	Methods: []grpc.MethodDesc{
		{
			MethodName: "newAuthCodeToken",
			Handler:    _TokenService_NewAuthCodeToken_Handler,
		},
		{
			MethodName: "newAuthCodeGrant",
			Handler:    _TokenService_NewAuthCodeGrant_Handler,
		},
		{
			MethodName: "verifyAuthCodeToken",
			Handler:    _TokenService_VerifyAuthCodeToken_Handler,
		},
		{
			MethodName: "newAuthCodeRefreshToken",
			Handler:    _TokenService_NewAuthCodeRefreshToken_Handler,
		},
		{
			MethodName: "authCodeRefreshTokenExists",
			Handler:    _TokenService_AuthCodeRefreshTokenExists_Handler,
		},
		{
			MethodName: "addUserIdAuthCodeFlow",
			Handler:    _TokenService_AddUserIdAuthCodeFlow_Handler,
		},
		{
			MethodName: "getUserIdfromAccesstoken",
			Handler:    _TokenService_GetUserIdfromAccesstoken_Handler,
		},
		{
			MethodName: "createToken",
			Handler:    _TokenService_CreateToken_Handler,
		},
		{
			MethodName: "verifyToken",
			Handler:    _TokenService_VerifyToken_Handler,
		},
		{
			MethodName: "refreshToken",
			Handler:    _TokenService_RefreshToken_Handler,
		},
		{
			MethodName: "revokeToken",
			Handler:    _TokenService_RevokeToken_Handler,
		},
		{
			MethodName: "revokeAll",
			Handler:    _TokenService_RevokeAll_Handler,
		},
		{
			MethodName: "createResetPasswordToken",
			Handler:    _TokenService_CreateResetPasswordToken_Handler,
		},
	},
	Streams:  []grpc.StreamDesc{},
	Metadata: "protos/token.proto",
}
