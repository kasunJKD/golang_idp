package main

import (
	"context"
	"idp-repository/env"
	db "idp-repository/login-repository/internal/db"
	"idp-repository/pkg/logger"
	"idp-repository/pkg/postgres"
	pb "idp-repository/protos/login"
	"log"
	"net"

	grpc "google.golang.org/grpc"
	"google.golang.org/protobuf/types/known/wrapperspb"
)

//var (
//	grpcPort = getEnv("GRPC_PORT", ":8083")
//	user_grpcPort = getEnv("USER_GRPC_PORT", ":8081")
//	user_host = getEnv("USER_HOST", "user-repository")
//	token_grpcPort = getEnv("TOKEN_GRPC_PORT", ":8082")
//	token_host = getEnv("TOKEN_HOST", "token-repository")
//	mem_host = getEnv("MEM_HOST", "membership-db")
//	mem_port = getEnv("MEM_PORT", "5432")
//	user = getEnv("DB_USER", "postgres")
//	password = getEnv("PASSWORD", "9221")
//	mem_dbname = getEnv("MEM_DBNAME", "membership")
//)

//func getEnv(key, fallback string) string {
//	value, found := os.LookupEnv(key)
//	if found {
//		return value
//	}
//	log.Println("Key not found: ", key)
//	os.Setenv(key, fallback)
//	return fallback
//}

type LoginServiceServer struct {
	pb.UnimplementedLoginServiceServer
	Idb db.IMemDB
	log *logger.CustomLogger
}

func NewLoginServiceServer(Idb db.IMemDB, log *logger.CustomLogger) *LoginServiceServer {
	return &LoginServiceServer{Idb: Idb, log: log}
}

func main() {
	customlogger := logger.NewCustomLogger()

	customlogger.Infolog.Println("Welcome to the server")

	env , err := env.GetEnv()
	if err != nil {
		customlogger.Errorlog.Fatalf("Loading env: %v", err)
	}

	//connecting to Database RW user
	rwDB, err := postgres.Connect(env)
	if err != nil {
		customlogger.Errorlog.Fatal("RW Postgresql init: ", err)
	}
	defer rwDB.Close()

	//connecting to Database RO user
	roDB, err := postgres.Connect(env)
	if err != nil {
		customlogger.Errorlog.Fatal("RO Postgresql init: ", err)
	}
	defer roDB.Close()


	//initialize database
	memDB := db.NewDBInstiate(roDB, rwDB, customlogger)

	//start listening for grpc
	listen, err := net.Listen("tcp", env.Server.GrpcPort)
	if err != nil {
		log.Fatal(err)
	}

	grpcserver := grpc.NewServer()

	GrpcMicroservice := NewLoginServiceServer(memDB, customlogger)


	//Register DataService
	pb.RegisterLoginServiceServer(grpcserver, GrpcMicroservice)
	log.Println("Starting grpc connection on port " + env.Server.GrpcPort)

	//startServing requests
	grpcserver.Serve(listen)


}

func (s *LoginServiceServer) GetClientInfo(ctx context.Context, req *pb.ClientReq) (*pb.ClientRes, error) {
	data, err := s.Idb.GetClientInfo(ctx, req)
	return data, err
}

func (s *LoginServiceServer) ValidateClientId(ctx context.Context, req *pb.ClientReq) (res *wrapperspb.BoolValue, err error) {
	data, err := s.Idb.CheckClientidExists(req)
	response := &wrapperspb.BoolValue{
		Value: data,
	}

	return response, err
}

func (s *LoginServiceServer) CreateOauthClient(ctx context.Context, req *pb.ClientReq) (*pb.ClientRes, error) {
	data, err := s.Idb.CreateOauthClient(req)
	return data, err
}

func (s *LoginServiceServer) ChangePassword (ctx context.Context, req *pb.Passwordreq) (*pb.Status, error) {
	data, err := s.Idb.ChangePassword(req)
	return data, err
}

func (s *LoginServiceServer) GetClientsByUserId(ctx context.Context, req *pb.ClientReq) (*pb.GetClientsByUserIdResponse, error) {
	data, err := s.Idb.GetClientsByUserId(ctx, req)
	return data, err
}

func (s *LoginServiceServer) ResetPassword (ctx context.Context, req *pb.Request) (*pb.Status, error) {
	data, err := s.Idb.ResetPassword(ctx, req)
	return data, err
}

func (s *LoginServiceServer) CheckotpEnabled (ctx context.Context, req *pb.OtpRequest) (res *pb.OtpResponse, err error) {
	enabled, err := s.Idb.Check2faEnabled(req)
	res = &pb.OtpResponse{
		OtpEnabled: enabled,
	}
	return res, err
}

func (s *LoginServiceServer) UpdateotpStatus (ctx context.Context, req *pb.OtpRequest) (res *pb.OtpResponse, err error) {
	status, err := s.Idb.Update2faStatus(req)
	res = &pb.OtpResponse{
		OtpEnabled: status,
	}
	return res, err
}