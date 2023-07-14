package db

import (
	"context"

	pb "idp-repository/protos/login"
	//"database/sql"
	"log"
	"time"

	_ "github.com/lib/pq"
	timestamppb "google.golang.org/protobuf/types/known/timestamppb"
)

func (c DBConfig) GetClientInfo(ctx context.Context, req *pb.ClientReq) (*pb.ClientRes, error) {

	ClientID := req.GetClientId()
	log.Println(req)
	sqlStatement := `select distinct clientId, clientName, clientSecret, projectId, userId, redirectUrl, createdAt, updatedAt, active
				from clients 
				where clientId= $1`

	var (
		clientId        string
		clientName      string
		clientSecret    string
		projectId  		string
		userId  		string
		redirectUrl  	string
		createdAt     time.Time
		updatedAt     time.Time
		active bool
	)

	err := c.roDb.QueryRow(sqlStatement, ClientID).Scan(&clientId, &clientName, &clientSecret, &projectId, &userId, &redirectUrl, &createdAt, &updatedAt, &active)

	if err != nil {
		log.Fatalln(err)
	}

	res := &pb.ClientRes{
		ClientId: clientId,
		ClientName: clientName,
		ClientSecret: clientSecret,
		ProjectId: projectId,
		UserId: userId,
		RedirectUrl: redirectUrl,
		CreatedAt:     timestamppb.New(createdAt),
		UpdatedAt:     timestamppb.New(updatedAt),
		Active: active,
	}

	return res, err

}
