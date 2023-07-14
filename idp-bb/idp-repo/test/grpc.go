package main

import (
	pb "idp-repository/protos/login"
	//pbu "idp-repository/protos/user"
	"context"
	"fmt"
	"google.golang.org/grpc"
	"google.golang.org/grpc/credentials/insecure"
	"log"
)

func testgrpc() {
	opts := []grpc.DialOption{
       //grpc.WithInitialWindowSize(65536),
       //grpc.WithInitialConnWindowSize(largerWindowSize),
	   grpc.WithTransportCredentials(insecure.NewCredentials()),
    }
	conn, err := grpc.Dial("localhost:8083", opts...)
	//conn, err := grpc.Dial("localhost:8081", opts...)
	if err != nil {
		log.Println(err)
	}

	client := pb.NewLoginServiceClient(conn)
	//client := pbu.NewUserServiceClient(conn)
	
	//md := metadata.New(map[string]string{"authorization": "eyJhbGciOiJIUzI1NiIsInR5cCI6IkpXVCJ9.eyJsb2NhbElkIjoiYmJlNTE2ZjYtM2UzOS0xMWViLWI4OTctMDg2MjY2MGNjYmQ0IiwidXNlclR5cGUiOiJhZG1pbiIsImlzcyI6IkFsaSIsInN1YiI6IlRlc3QiLCJhdWQiOiJQOTkiLCJleHAiOjE2NjI2NjQyOTEsIm5iZiI6MCwiaWF0IjoxNjYyNjYzOTkxLCJqdGkiOiIifQ.U1hxNrPfKdfo7hBJWbUgrSQHjJjpdyR2lEYtjfUY_ko"})
	//ctx := metadata.NewOutgoingContext(context.Background(), md)

	// Testing the functions
	//resp, err := client.SignUp(context.Background(), &pb.Request{Email: "kas@gmail", Password: "test", DisplayName: "test"})
	//resp, err := client.CreateAuthUrl(context.Background(), &pb.AuthUrlRequest{ProviderId: "google", ClientId: "299401868668-0kfkdqnk6kggsutt5nkp4nq9m2q9ouis.apps.googleusercontent.com"})

	//resp, err := client.CreateNewUser(context.Background(), &pbu.Request{Email: "k@kgasdasd", DisplayName: "kasu"})
	resp, err := client.OtpLogin(context.Background(), &pb.OtpLoginRequest{UserId: "asdasd"})	
	if err != nil {
		log.Println(err)
	}
	fmt.Println(resp)

}
