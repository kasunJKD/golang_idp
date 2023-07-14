package main

import (
	"context"
	"fmt"
	"google.golang.org/grpc"
	"google.golang.org/grpc/credentials/insecure"
	pb "idp-repository/protos/token"
	"log"
)

func test_token() {
	conn, err := grpc.Dial("localhost:8082", grpc.WithTransportCredentials(insecure.NewCredentials()))
	if err != nil {
		log.Println(err)
	}

	client := pb.NewTokenServiceClient(conn)

	//md := metadata.New(map[string]string{"authorization": "eyJhbGciOiJIUzI1NiIsInR5cCI6IkpXVCJ9.eyJsb2NhbElkIjoiYmJlNTE2ZjYtM2UzOS0xMWViLWI4OTctMDg2MjY2MGNjYmQ0IiwidXNlclR5cGUiOiJhZG1pbiIsImlzcyI6IkFsaSIsInN1YiI6IlRlc3QiLCJhdWQiOiJQOTkiLCJleHAiOjE2NjI2NjQyOTEsIm5iZiI6MCwiaWF0IjoxNjYyNjYzOTkxLCJqdGkiOiIifQ.U1hxNrPfKdfo7hBJWbUgrSQHjJjpdyR2lEYtjfUY_ko"})
	//ctx := metadata.NewOutgoingContext(context.Background(), md)

	// Testing the functions
	resp, err := client.CreateToken(context.Background(), &pb.User{UserId: "b826df6a-e137-4e3c-a2a5-b947e0a229ce"})
	//resp, err := client.VerifyToken(context.Background(), &pb.TokenRequest{Token: "eyJhbGciOiJIUzI1NiIsInR5cCI6IkpXVCJ9.eyJsb2NhbElkIjoiYjgyNmRmNmEtZTEzNy00ZTNjLWEyYTUtYjk0N2UwYTIyOWNlIiwiaXNzIjoiQWxpIiwic3ViIjoiVGVzdCIsImF1ZCI6IlA5OSIsImV4cCI6MTY2ODYyNzQwMiwiaWF0IjoxNjY4NjI3MTAyfQ.6dj-U1d5j4l73L3HyQjUigtsi8uK5QrwBIckYAkBaq8"})
	//resp, err := client.RefreshToken(context.Background(), &pb.TokenRequest{RefreshToken: "eyJhbGciOiJIUzI1NiIsInR5cCI6IkpXVCJ9.eyJsb2NhbElkIjoiYjgyNmRmNmEtZTEzNy00ZTNjLWEyYTUtYjk0N2UwYTIyOWNlIiwiaXNzIjoiQWxpIiwic3ViIjoiVGVzdCIsImF1ZCI6IlA5OSIsImV4cCI6MTY2OTExNjMwNiwiaWF0IjoxNjY4NTExNTA2fQ._gai42T6jjauKonqzj4r0vnYtYhEZRU5CWiiJSdZt7M"})
	//resp, err := client.RevokeToken(context.Background(), &pb.TokenRequest{Token: "aaa"})
	//resp, err := client.RevokeAll(context.Background(), &pb.User{UserId: "k"})
	if err != nil {
		log.Println(err)
	}
	fmt.Println(resp)

}
