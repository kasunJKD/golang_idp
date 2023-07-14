package main

import (
	"context"
	"fmt"
	"google.golang.org/grpc"
	"google.golang.org/grpc/credentials/insecure"
	"google.golang.org/grpc/metadata"
	pb "idp-repository/protos/user"
	"log"
)

func test_user() {
	conn, err := grpc.Dial("localhost:8081", grpc.WithTransportCredentials(insecure.NewCredentials()))
	if err != nil {
		log.Println(err)
	}

	client := pb.NewUserServiceClient(conn)

	md := metadata.New(map[string]string{"authorization": "eyJhbGciOiJIUzI1NiIsInR5cCI6IkpXVCJ9.eyJsb2NhbElkIjoiYmJlNTE2ZjYtM2UzOS0xMWViLWI4OTctMDg2MjY2MGNjYmQ0IiwidXNlclR5cGUiOiJhZG1pbiIsImlzcyI6IkFsaSIsInN1YiI6IlRlc3QiLCJhdWQiOiJQOTkiLCJleHAiOjE2NjI2NjQyOTEsIm5iZiI6MCwiaWF0IjoxNjYyNjYzOTkxLCJqdGkiOiIifQ.U1hxNrPfKdfo7hBJWbUgrSQHjJjpdyR2lEYtjfUY_ko"})
	ctx := metadata.NewOutgoingContext(context.Background(), md)

	// Testing the functions
	resp, err := client.SetAccountInfo(ctx, &pb.Request{LocalId: "b826df6a-e137-4e3c-a2a5-b947e0a229ce", ProviderId: "password", FederatedId: "ali@gmail.com", DisplayName: "Bob", FirstName: "", LastName: "", PhotoUrl: ""})
	if err != nil {
		log.Println(err)
	}
	fmt.Println(resp)

}
