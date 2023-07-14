package db

import (
	"google.golang.org/protobuf/types/known/timestamppb"
	"time"

	pb "idp-repository/protos/login"
	"log"

	_ "github.com/lib/pq"
	"google.golang.org/grpc/codes"
	"google.golang.org/grpc/status"
	"idp-repository/login-repository/internal/crypto"
)

func (c DBConfig) CreateOauthClient(req *pb.ClientReq) (*pb.ClientRes, error) {
	clientId := crypto.NewClientId()
	clientSecret := crypto.NewClientSecret()
	redirectUri := req.GetRedirectUrl()
	clientName := req.GetClientName()
	userId := req.GetUserId()
	projectId := req.GetProjectId()

	//check name exists
	b, _ := c.CheckClientNameExists(req)
	if b {
		return nil, status.Errorf(codes.AlreadyExists, "client name already exists")
	}

	sqlStatement := `
		INSERT INTO clients(id, clientId, clientName, clientSecret, projectId, userId, redirectUrl, createdAt, updatedAt, active)
		VALUES (gen_random_uuid(), $1, $2, $3, $4, $5, $6, (select current_timestamp at time zone ('utc')), (select current_timestamp at time zone ('utc')), TRUE)
		RETURNING createdAt;
	`

	var (
		createdAt time.Time
	)

	err := c.rwDb.QueryRow(sqlStatement,
		clientId,
		clientName,
		clientSecret,
		projectId,
		userId,
		redirectUri).Scan(&createdAt)

	if err != nil {
		log.Fatalln(err)
		return nil, status.Errorf(codes.Internal, "cannot create client")
	}

	res := &pb.ClientRes{
		ClientId: clientId,
		ProjectId: projectId,
		UserId: userId,
		RedirectUrl: redirectUri,
		CreatedAt: timestamppb.New(createdAt),
		Active: true,
	}

	return res, nil
}
