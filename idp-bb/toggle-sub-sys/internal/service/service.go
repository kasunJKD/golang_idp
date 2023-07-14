package service

import (
	"context"
	"login-throttles/internal/redis"
	proto "login-throttles/protos/core"
)

type CoreService struct {
	proto.UnimplementedCoreServiceServer
}

//NewCoreService create service
func NewCoreService() proto.CoreServiceServer {
	return &CoreService{}
}

func (service *CoreService) CheckStatus(ctx context.Context, req *proto.Request) (res *proto.Response, err error) {
	res, err = redis.CheckStatus(req)

	return res, err
}

func (service *CoreService) UpdateStatus(ctx context.Context, req *proto.Request) (res *proto.Response, err error) {
	res, err = redis.UpdateStatus(req)

	return res, err
}

func (service *CoreService) CheckLimitCounter(ctx context.Context, req *proto.Request) (res *proto.Response, err error) {
	res, err = redis.CheckLimitCounter(req)

	return res, err
}

func (service *CoreService) UpdateLimitCounter(ctx context.Context, req *proto.Request) (res *proto.Response, err error) {
	res, err = redis.UpdateLimitCounter(req)

	return res, err
}
