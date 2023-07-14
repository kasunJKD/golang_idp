package grpc_client

import (
	pb_login "idp-service/protos/login"
	pb_token "idp-service/protos/token"
	pb_user "idp-service/protos/user"

	pb_otp "bitbucket.org/project-99-games/otp_model"
	"google.golang.org/grpc"
	"google.golang.org/grpc/credentials/insecure"
	//"log"
)

type LoginServiceClient struct {
	Conn   *grpc.ClientConn
	Client pb_login.LoginServiceClient
}

type UserServiceClient struct {
	Conn   *grpc.ClientConn
	Client pb_user.UserServiceClient
}

type TokenServiceClient struct {
	Conn   *grpc.ClientConn
	Client pb_token.TokenServiceClient
}


type OTPServiceClient struct {
	Conn   *grpc.ClientConn
	Client pb_otp.CoreServiceClient
}

func NewLoginServiceClient(host, port string)(*LoginServiceClient, error){
	conn, err := grpc.Dial(host + port, grpc.WithTransportCredentials(insecure.NewCredentials()))
	if err != nil {
		return nil, err
	}

	client := pb_login.NewLoginServiceClient(conn)

	return &LoginServiceClient {
		Conn: conn,
		Client: client,
	}, nil
}

func NewOTPServiceClient(host, port string)(*OTPServiceClient, error){
	conn, err := grpc.Dial(host + port, grpc.WithTransportCredentials(insecure.NewCredentials()))
	if err != nil {
		return nil, err
	}

	client := pb_otp.NewCoreServiceClient(conn)

	return &OTPServiceClient {
		Conn: conn,
		Client: client,
	}, nil
}

func NewUserServiceClient(host, port string)(*UserServiceClient, error){
	conn, err := grpc.Dial(host + port, grpc.WithTransportCredentials(insecure.NewCredentials()))
	if err != nil {
		return nil, err
	}

	client := pb_user.NewUserServiceClient(conn)

	return &UserServiceClient {
		Conn: conn,
		Client: client,
	}, nil
}

func NewTokenServiceClient(host, port string)(*TokenServiceClient, error){
	conn, err := grpc.Dial(host + port, grpc.WithTransportCredentials(insecure.NewCredentials()))
	if err != nil {
		return nil, err
	}

	client := pb_token.NewTokenServiceClient(conn)

	return &TokenServiceClient {
		Conn: conn,
		Client: client,
	}, nil
}
