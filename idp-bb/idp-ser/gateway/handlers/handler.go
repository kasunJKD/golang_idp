package handlers

import (
	"idp-service/env"
	"idp-service/pkg/grpc_client"
	"idp-service/pkg/logger"
)

//auth package auth function init
type HandlerFunc struct {
	l *grpc_client.LoginServiceClient
	t *grpc_client.TokenServiceClient
	u *grpc_client.UserServiceClient
	ll *logger.CustomLogger
	env *env.Env
}

//auth package auth function init 
func NewHandlerFunc(l *grpc_client.LoginServiceClient, t *grpc_client.TokenServiceClient, u *grpc_client.UserServiceClient, ll *logger.CustomLogger, env *env.Env) *HandlerFunc {
	return &HandlerFunc{
		l: l,
		t: t,
		u: u,
		ll: ll,
		env: env,
	}
}