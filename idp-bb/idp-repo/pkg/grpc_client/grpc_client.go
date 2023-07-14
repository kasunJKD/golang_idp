package grpc_client

import (
	pb_login "idp-repository/protos/login"
	pb_token "idp-repository/protos/token"
	pb_user "idp-repository/protos/user"
	"google.golang.org/grpc"
	"google.golang.org/grpc/credentials/insecure"
	"log"
)

func ConnectUser(host, port string) pb_user.UserServiceClient {
	grpcConn, err := grpc.Dial(host + port, grpc.WithTransportCredentials(insecure.NewCredentials()))
	if err != nil {
		log.Println(err)
	}
	client := pb_user.NewUserServiceClient(grpcConn)
	log.Println("Starting grpc connection with user-repository")

	return client
}

func ConnectToken(host, port string) pb_token.TokenServiceClient {
	grpcConn, err := grpc.Dial(host + port, grpc.WithTransportCredentials(insecure.NewCredentials()))
	if err != nil {
		log.Println(err)
	}
	client := pb_token.NewTokenServiceClient(grpcConn)
	log.Println("Starting grpc connection with token-repository")

	return client
}

func ConnectLogin(host, port string) pb_login.LoginServiceClient {
	grpcConn, err := grpc.Dial(host + port, grpc.WithTransportCredentials(insecure.NewCredentials()))
	if err != nil {
		log.Println(err)
	}
	client := pb_login.NewLoginServiceClient(grpcConn)
	log.Println("Starting grpc connection with login-repository")

	return client
}