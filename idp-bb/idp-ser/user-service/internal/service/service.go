package service

import (
	"context"
	"google.golang.org/grpc/metadata"
	"google.golang.org/protobuf/types/known/wrapperspb"
	"idp-service/env"
	"idp-service/pkg/grpc_client"
	"idp-service/pkg/logger"
	pb_token "idp-service/protos/token"
	proto "idp-service/protos/user"
)

//LoginService struct
type UserService struct {
	proto.UnimplementedUserServiceServer
	u *grpc_client.UserServiceClient
	t *grpc_client.TokenServiceClient
	ll *logger.CustomLogger
	env *env.Env
}

//NewLoginService create service
func NewUserService(u *grpc_client.UserServiceClient, t *grpc_client.TokenServiceClient, ll *logger.CustomLogger, env *env.Env) proto.UserServiceServer {
	return &UserService{
		u: u,
		t:t,
		ll:ll,
		env: env,
	}
}

func (service *UserService) GetLinkedAccountInfo(ctx context.Context, req *proto.Request) (res *proto.Response, err error) {
	service.ll.Infolog.Printf("Service => GetLinkedAccountInfo function called")

	md, _ := metadata.FromIncomingContext(ctx)
    token := md.Get("authorization")[0]
	//set Headers
	header := metadata.New(map[string]string{"authorization": token})
	c := metadata.NewOutgoingContext(ctx, header)

	service.ll.Infolog.Printf("Verifying token")
	_, err = service.t.Client.VerifyToken(c, &pb_token.TokenRequest{Token: req.AccessToken})
	if err != nil {
		if service.env.Server.Mode == "debug" {
			service.ll.Debuglog.Printf("err %v", err)
		}
		service.ll.Errorlog.Println("Token Service ===> VerifyToken failed")
		return nil, err
	}
	//Dial repository
	//return response from repository
	res, err = service.u.Client.GetLinkedAccountInfo(ctx, req)
	if err != nil {
		if service.env.Server.Mode == "debug" {
			service.ll.Debuglog.Printf("err %v", err)
		}
		service.ll.Errorlog.Println("User Repository ===> GetLinkedAccountInfo failed")
		return nil, err
	}
	return res, err
}

func (service *UserService) GetUserInfoById(ctx context.Context, req *proto.Request) (res *proto.Response, err error) {
	service.ll.Infolog.Printf("Service => GetUserInfoById function called")
	// md, _ := metadata.FromIncomingContext(ctx)
    // token := md.Get("authorization")[0]
	// //set Headers
	// header := metadata.New(map[string]string{"authorization": token})
	// c := metadata.NewOutgoingContext(ctx, header)
	//Dial repository
	//return response from repository
	res, err = service.u.Client.GetUserInfoById(ctx, req)
	if err != nil {
		if service.env.Server.Mode == "debug" {
			service.ll.Debuglog.Printf("err %v", err)
		}
		service.ll.Errorlog.Println("User Repository ===> GetUserInfoById failed")
		return nil, err
	}
	return res, err
}

func (service *UserService) UnlinkAccount(ctx context.Context, req *proto.Request) (res *proto.Response, err error) {
	service.ll.Infolog.Printf("Service => UnlinkAccount function called")
	md, _ := metadata.FromIncomingContext(ctx)
    token := md.Get("authorization")[0]
	//set Headers
	header := metadata.New(map[string]string{"authorization": token})
	c := metadata.NewOutgoingContext(ctx, header)

	_, err = service.t.Client.VerifyToken(c, &pb_token.TokenRequest{Token: req.AccessToken})
	if err != nil {
		if service.env.Server.Mode == "debug" {
			service.ll.Debuglog.Printf("err %v", err)
		}
		service.ll.Errorlog.Println("Token Service ===> VerifyToken failed")
		return nil, err
	}

	//Dial repository
	//return response from repository
	res, err = service.u.Client.UnlinkAccount(ctx, req)
	if err != nil {
		if service.env.Server.Mode == "debug" {
			service.ll.Debuglog.Printf("err %v", err)
		}
		service.ll.Errorlog.Println("User Repository ===> UnlinkAccount failed")
		return nil, err
	}
	return res, err
}

func (service *UserService) SetAccountInfo(ctx context.Context, req *proto.Request) (res *proto.Response, err error) {
	service.ll.Infolog.Printf("Service => SetAccountInfo function called")

	md, _ := metadata.FromIncomingContext(ctx)
    token := md.Get("authorization")[0]
	//set Headers
	header := metadata.New(map[string]string{"authorization": token})
	c := metadata.NewOutgoingContext(ctx, header)

	_, err = service.t.Client.VerifyToken(c, &pb_token.TokenRequest{Token: req.AccessToken})
	if err != nil {
		if service.env.Server.Mode == "debug" {
			service.ll.Debuglog.Printf("err %v", err)
		}
		service.ll.Errorlog.Println("Token Service ===> VerifyToken failed")
		return nil, err
	}

	//Dial repository
	//return response from repository
	res, err = service.u.Client.SetAccountInfo(ctx, req)
	if err != nil {
		if service.env.Server.Mode == "debug" {
			service.ll.Debuglog.Printf("err %v", err)
		}
		service.ll.Errorlog.Println("User Repository ===> SetAccountInfo failed")
		return nil, err
	}
	return res, err
}

func (service *UserService) CheckAuthUserExists(ctx context.Context, req *proto.Request) (res *wrapperspb.BoolValue, err error) {
	service.ll.Infolog.Printf("Service => CheckAuthUserExists function called")
	data, err := service.u.Client.CheckAuthUserExists(ctx, req)
	if err != nil {
		if service.env.Server.Mode == "debug" {
			service.ll.Debuglog.Printf("err %v", err)
		}
		service.ll.Errorlog.Println("User Repository ===> CheckAuthUserExists failed")
		return nil, err
	}
	response := &wrapperspb.BoolValue{
		Value: data.Value,
	}

	return response, err
}

func (service *UserService)  CreateNewUser(ctx context.Context, req *proto.Request) (*proto.Response, error) {
	service.ll.Infolog.Printf("Service => CreateNewUser function called")
	data, err := service.u.Client.CreateNewUser(ctx, req)
	if err != nil {
		if service.env.Server.Mode == "debug" {
			service.ll.Debuglog.Printf("err %v", err)
		}
		service.ll.Errorlog.Println("User Repository ===> CreateNewUser failed")
		return nil, err
	}
	return data, err
}

func (service *UserService) GetAccountInfo(ctx context.Context, req *proto.Request) (res *proto.Response, err error) {
	service.ll.Infolog.Printf("Service => GetAccountInfo function called")
	// md, _ := metadata.FromIncomingContext(ctx)
	// token := md.Get("authorization")[0]
	// //set Headers
	// header := metadata.New(map[string]string{"authorization": token})
	// c := metadata.NewOutgoingContext(ctx, header)
	//Dial repository
	//return response from repository
	res, err = service.u.Client.GetAccountInfo(ctx, req)
	if err != nil {
		if service.env.Server.Mode == "debug" {
			service.ll.Debuglog.Printf("err %v", err)
		}
		service.ll.Errorlog.Println("User Repository ===> GetAccountInfo failed")
		return nil, err
	}
	return res, err
}
