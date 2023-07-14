package service

import (
	"context"
	"google.golang.org/protobuf/types/known/wrapperspb"
	"idp-service/env"
	"idp-service/pkg/grpc_client"
	"idp-service/pkg/logger"
	proto "idp-service/protos/token"
	"idp-service/token-service/internal/jwt"
)

//LoginService struct
type TokenService struct {
	proto.UnimplementedTokenServiceServer
	t *grpc_client.TokenServiceClient
	ll *logger.CustomLogger
	j jwt.IjwtFunc
	env *env.Env
}

//NewTokenService create service
func NewTokenService(t *grpc_client.TokenServiceClient, ll *logger.CustomLogger, j jwt.IjwtFunc, env *env.Env) proto.TokenServiceServer {
	return &TokenService{
		t: t,
		ll: ll,
		j: j,
		env: env,
	}
}

func (service *TokenService) NewAuthCodeToken(ctx context.Context, req *proto.TokenRequest) (res *proto.AuthCodeToken, err error) {
	service.ll.Infolog.Printf("Service =>  NewAuthCodeToken function called\nCalling token client")
	res, err = service.t.Client.NewAuthCodeToken(ctx, req)
	if err != nil {
		if service.env.Server.Mode == "debug" {
			service.ll.Debuglog.Printf("err %v", err)
		}
		service.ll.Errorlog.Println("Token repository ===> NewAuthCodeToken failed")
		return nil, err
	}
	return res, err
}

func (service *TokenService) NewAuthCodeGrant(ctx context.Context, req *proto.TokenRequest) (res *wrapperspb.StringValue, err error) {
	service.ll.Infolog.Printf("Service =>  NewAuthCodeGrant function called\nCalling token client")
	res, err = service.t.Client.NewAuthCodeGrant(ctx, req)
	if err != nil {
		if service.env.Server.Mode == "debug" {
			service.ll.Debuglog.Printf("err %v", err)
		}
		service.ll.Errorlog.Println("Token repository ===> NewAuthCodeGrant failed")
		return nil, err
	}
	return res, err
}

func (service *TokenService) VerifyAuthCodeToken(ctx context.Context, req *proto.TokenRequest) (res *wrapperspb.BoolValue, err error) {
	service.ll.Infolog.Printf("Service =>  VerifyAuthCodeToken function called\nCalling token client")
	res, err = service.t.Client.VerifyAuthCodeToken(ctx, req)
	if err != nil {
		if service.env.Server.Mode == "debug" {
			service.ll.Debuglog.Printf("err %v", err)
		}
		service.ll.Errorlog.Println("Token repository ===> VerifyAuthCodeToken failed")
		return nil, err
	}
	return res, err
}

func (service *TokenService) NewAuthCodeRefreshToken(ctx context.Context, req *proto.RefreshTokenRequest) (res *proto.AuthCodeToken, err error) {
	service.ll.Infolog.Printf("Service =>  NewAuthCodeRefreshToken function called\nCalling token client")
	res, err = service.t.Client.NewAuthCodeRefreshToken(ctx, req)
	if err != nil {
		if service.env.Server.Mode == "debug" {
			service.ll.Debuglog.Printf("err %v", err)
		}
		service.ll.Errorlog.Println("Token repository ===> NewAuthCodeRefreshToken failed")
		return nil, err
	}
	return res, err
}

func (service *TokenService) AuthCodeRefreshTokenExists(ctx context.Context, req *proto.RefreshTokenRequest) (res *wrapperspb.BoolValue, err error) {
	service.ll.Infolog.Printf("Service =>  AuthCodeRefreshTokenExists function called\nCalling token client")
	res, err = service.t.Client.AuthCodeRefreshTokenExists(ctx, req)
	if err != nil {
		if service.env.Server.Mode == "debug" {
			service.ll.Debuglog.Printf("err %v", err)
		}
		service.ll.Errorlog.Println("Token repository ===> AuthCodeRefreshTokenExists failed")
		return nil, err
	}
	return res, err
}

func (service *TokenService) AddUserIdAuthCodeFlow(ctx context.Context, req *proto.User) (res *wrapperspb.BoolValue, err error) {
	service.ll.Infolog.Printf("Service =>  AddUserIdAuthCodeFlow function called\nCalling token client")
	res, err = service.t.Client.AddUserIdAuthCodeFlow(ctx, req)
	if err != nil {
		if service.env.Server.Mode == "debug" {
			service.ll.Debuglog.Printf("err %v", err)
		}
		service.ll.Errorlog.Println("Token repository ===> AddUserIdAuthCodeFlow failed")
		return nil, err
	}
	return res, err
}

func (service *TokenService) GetUserIdfromAccesstoken(ctx context.Context, req *proto.User) (res *wrapperspb.StringValue, err error) {
	service.ll.Infolog.Printf("Service =>  GetUserIdfromAccesstoken function called\nCalling token client")
	res, err = service.t.Client.GetUserIdfromAccesstoken(ctx, req)
	if err != nil {
		if service.env.Server.Mode == "debug" {
			service.ll.Debuglog.Printf("err %v", err)
		}
		service.ll.Errorlog.Println("Token repository ===> GetUserIdfromAccesstoken failed")
		return nil, err
	}
	return res, err
}

func (service *TokenService) CreateToken(ctx context.Context, req *proto.User) (res *proto.AuthCodeToken, err error) {
	service.ll.Infolog.Printf("Service => CreateToken function called")
	res, err = service.j.CreateToken(req)
	if err != nil {
		if service.env.Server.Mode == "debug" {
			service.ll.Debuglog.Printf("err %v", err)
		}
		service.ll.Errorlog.Println("jwt ===> CreateToken failed")
		return nil, err
	}
	return res, err
}

func (service *TokenService) VerifyToken(ctx context.Context, req *proto.TokenRequest) (res *proto.AuthCodeToken, err error) {
	service.ll.Infolog.Printf("Service => VerifyToken function called")
	res, err = service.j.VerifyToken(ctx, req)
	if err != nil {
		if service.env.Server.Mode == "debug" {
			service.ll.Debuglog.Printf("err %v", err)
		}
		service.ll.Errorlog.Println("jwt ===> VerifyToken failed")
		return nil, err
	}
	return res, err
}

func (service *TokenService) RefreshToken(ctx context.Context, req *proto.TokenRequest) (res *proto.AuthCodeToken, err error) {
	service.ll.Infolog.Printf("Service => RefreshToken function called")
	res, err = service.j.RefreshToken(ctx, req)
	if err != nil {
		if service.env.Server.Mode == "debug" {
			service.ll.Debuglog.Printf("err %v", err)
		}
		service.ll.Errorlog.Println("jwt ===> RefreshToken failed")
		return nil, err
	}
	return res, err
}

func (service *TokenService) RevokeToken(ctx context.Context, req *proto.TokenRequest) (res *proto.AuthCodeToken, err error) {
	service.ll.Infolog.Printf("Service =>  RevokeToken function called")
	res, err = service.j.RevokeToken(ctx, req)
	if err != nil {
		if service.env.Server.Mode == "debug" {
			service.ll.Debuglog.Printf("err %v", err)
		}
		service.ll.Errorlog.Println("jwt ===> RevokeToken failed")
		return nil, err
	}
	return res, err
}

func (service *TokenService) RevokeAll(ctx context.Context, req *proto.User) (res *proto.AuthCodeToken, err error) {
	service.ll.Infolog.Printf("Service =>  RevokeAll function called")
	res, err = service.j.RevokeAll(ctx, req)
	if err != nil {
		if service.env.Server.Mode == "debug" {
			service.ll.Debuglog.Printf("err %v", err)
		}
		service.ll.Errorlog.Println("jwt ===> RevokeAll failed")
		return nil, err
	}
	return res, err
}

func (service *TokenService) CreateResetPasswordToken(ctx context.Context, req *proto.User) (res *proto.AuthCodeToken, err error) {
	service.ll.Infolog.Printf("Service =>  CreateResetPasswordToken function called")
	res, err = service.j.CreateResetPasswordToken(ctx, req)
	if err != nil {
		if service.env.Server.Mode == "debug" {
			service.ll.Debuglog.Printf("err %v", err)
		}
		service.ll.Errorlog.Println("jwt ===> CreateResetPasswordToken failed")
		return nil, err
	}
	return res, err
}