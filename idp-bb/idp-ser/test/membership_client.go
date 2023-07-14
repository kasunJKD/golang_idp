package main

import (
	"context"
	"fmt"
	pb "idp-service/protos/login"
	"log"

	"google.golang.org/grpc"
	"google.golang.org/grpc/credentials/insecure"
)

func mem_test_client() {
	conn, err := grpc.Dial("localhost:48081", grpc.WithTransportCredentials(insecure.NewCredentials()))
	if err != nil {
		log.Println(err)
	}

	client := pb.NewLoginServiceClient(conn)

	// Testing the functions
	//resp, err := client.SignUp(context.Background(), &pb.Request{Email: "kassy@gmail.com", Password: "pass", DisplayName: "kassy"})
	resp, err := client.PasswordSignIn(context.Background(), &pb.Request{Email: "ali@gmail.com", Password: "pass"})
	//resp, err := client.OtpLogin(context.Background(), &pb.OtpLoginRequest{UserId: "183cdd74-bbcc-43f1-aeb2-bb3a2474995b", OtpCode: 947141})
	//resp, err := client.GetLinkedAccountInfo(context.Background(), &pb.Request{UserId: "05f31fcf-d4c8-44fb-9d46-e528382b08ed", ProviderId: "facebook"})
	//resp, err := client.UnlinkAccount(context.Background(), &pb.Request{ProviderId: "Apple", FederatedId: "1"})

	if err != nil {
		log.Println(err)
	}
	fmt.Println(resp)

}
