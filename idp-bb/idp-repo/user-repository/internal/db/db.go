package db

import (
	"context"
	"database/sql"
	"idp-repository/pkg/logger"
	pb "idp-repository/protos/user"
)

type DBConfig struct {
	roDb *sql.DB
	rwDb	*sql.DB
	log *logger.CustomLogger
}

func NewDBInstiate(roDB *sql.DB, rwDb	*sql.DB, log *logger.CustomLogger) *DBConfig {
	return &DBConfig{roDb: roDB, rwDb: rwDb,log: log}
}

type IMemDB interface {
	CheckIdpAccountLinked(req *pb.Request) (bool, error)
	CheckIdpProviderLinked(req *pb.Request) (bool, error)
	CheckAuthUserExists(req *pb.Request) (bool, error)
	CreateNewUser(req *pb.Request) (*pb.Response, error)
	GetAccountInfo(ctx context.Context, req *pb.Request) (*pb.Response, error)
	GetAccountInfoById(ctx context.Context, req *pb.Request) (*pb.Response, error)
	GetLinkedAccountInfo(ctx context.Context, req *pb.Request) (*pb.Response, error)
	LinkAccount(ctx context.Context, req *pb.Request) (*pb.Response, error)
	MigrateAccount(req *pb.Request) (*pb.Response, error)
	SetAccountInfo(ctx context.Context, req *pb.Request) (*pb.Response, error)
	UnlinkAccount(ctx context.Context, req *pb.Request) (*pb.Response, error)
}
