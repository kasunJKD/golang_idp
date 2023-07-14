package auth

import (
	"context"
	"idp-service/env"
	"idp-service/pkg/grpc_client"
	"idp-service/pkg/logger"

	pb "idp-service/protos/login"
)

//auth package auth function init
type AuthFunc struct {
	l *grpc_client.LoginServiceClient
	t *grpc_client.TokenServiceClient
	u *grpc_client.UserServiceClient
	ll *logger.CustomLogger
	env *env.Env
}

//auth package auth function init 
func NewAuthFunc(l *grpc_client.LoginServiceClient, t *grpc_client.TokenServiceClient, u *grpc_client.UserServiceClient, ll *logger.CustomLogger, env *env.Env) *AuthFunc {
	return &AuthFunc{
		l: l,
		t: t,
		u: u,
		ll: ll,
		env: env,
	}
}

type IAuthFunc interface {
	ForgotPassword(ctx context.Context, req *pb.Request) (*pb.ForgotPasswordResponse, error)
	MigrateAccount(ctx context.Context, req *pb.Request) (*pb.Response, error) 
	OtpLogin(ctx context.Context, req *pb.OtpLoginRequest) (*pb.Response, error)
	PasswordSignIn(ctx context.Context, req *pb.Request) (*pb.Response, error)
	SignUp(ctx context.Context, req *pb.Request) (*pb.Response, error)
}