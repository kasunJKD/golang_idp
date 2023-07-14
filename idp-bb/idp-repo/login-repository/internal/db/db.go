package db

import (
	"context"
	"database/sql"
	"idp-repository/pkg/logger"
	pb "idp-repository/protos/login"
)

type DBConfig struct {
	roDb *sql.DB
	rwDb	*sql.DB
	log *logger.CustomLogger
}

func NewDBInstiate(roDb *sql.DB, rwDb	*sql.DB, log *logger.CustomLogger) *DBConfig {
	return &DBConfig{roDb: roDb, rwDb: rwDb,log: log}
}

type IMemDB interface {
	Check2faEnabled(req *pb.OtpRequest) (int32, error)
	Update2faStatus (req *pb.OtpRequest) (int32, error)
	ChangePassword(req *pb.Passwordreq) (*pb.Status, error)
	checkPassword (req *pb.Passwordreq) (bool, error)
	CheckClientidExists(req *pb.ClientReq) (bool, error)
	CheckClientNameExists(req *pb.ClientReq) (bool, error)
	CreateOauthClient(req *pb.ClientReq) (*pb.ClientRes, error)
	GetClientInfo(ctx context.Context, req *pb.ClientReq) (*pb.ClientRes, error)
	GetClientsByUserId(ctx context.Context, req *pb.ClientReq) (*pb.GetClientsByUserIdResponse, error)
	ResetPassword(ctx context.Context, req *pb.Request) (*pb.Status, error)
}