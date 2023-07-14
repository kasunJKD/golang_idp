package db

import (
	"context"
	_ "github.com/lib/pq"
	"google.golang.org/grpc/codes"
	"google.golang.org/grpc/status"
	"google.golang.org/protobuf/types/known/timestamppb"
	pb "idp-repository/protos/login"
	"log"
	"time"
)

func (c DBConfig) GetClientsByUserId(ctx context.Context, req *pb.ClientReq) (*pb.GetClientsByUserIdResponse, error) {

	clients := make([]*pb.ClientRes, 0)
	user_Id := req.GetUserId()
	log.Println(req)

	sqlStatement := `select C.clientId, C.clientName, C.clientSecret, C.redirectUrl, C.createdAt, C.updatedAt, C.active
				from Clients C
				where C.userId = $1`

	var (
		clientId string
		clientName string
		clientSecret string
		redirectUrl string
		createdAt time.Time
		updatedAt time.Time
		active bool
	)

	sqlRows, err := c.roDb.Query(sqlStatement, user_Id)
	if err != nil {
		return nil, status.Errorf(codes.Internal, err.Error())
	}
	defer sqlRows.Close()

	for sqlRows.Next() {
		client := new(pb.ClientRes)
		if err := sqlRows.Scan(&clientId, &clientName, &clientSecret, &redirectUrl, &createdAt, &updatedAt, &active); err != nil {
			return nil, status.Errorf(codes.Internal, err.Error())
		}
		client = &pb.ClientRes {
			ClientId: clientId,
			ClientName: clientName,
			ClientSecret: clientSecret,
			RedirectUrl: redirectUrl,
			CreatedAt: timestamppb.New(createdAt),
			UpdatedAt: timestamppb.New(updatedAt),
			Active: active,
		}
		clients = append(clients, client)
	}

	res := &pb.GetClientsByUserIdResponse{
		Clients: clients,
	}

	return res, err

}
