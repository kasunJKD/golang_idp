package main

import (
	"context"
	"idp-repository/env"
	"idp-repository/pkg/logger"
	"idp-repository/pkg/postgres"
	pb "idp-repository/protos/user"
	db "idp-repository/user-repository/internal/db"
	"log"
	"net"

	grpc "google.golang.org/grpc"
	"google.golang.org/protobuf/types/known/wrapperspb"
)

//var dbConn *sql.DB

//var (
//	grpcPort = getEnv("GRPC_PORT", ":8081")
//	token_grpcPort = getEnv("TOKEN_GRPC_PORT", ":8082")
//	token_host = getEnv("TOKEN_HOST", "token-repository")
//	mem_host = getEnv("MEM_HOST", "membership-db")
//	mem_port = getEnv("MEM_PORT", "5432")
//	user = getEnv("DB_USER", "postgres")
//	password = getEnv("PASSWORD", "9221")
//	mem_dbname = getEnv("MEM_DBNAME", "membership")
//)
//
//func getEnv(key, fallback string) string {
//	value, found := os.LookupEnv(key)
//	if found {
//		return value
//	}
//	log.Println("Key not found: ", key)
//	os.Setenv(key, fallback)
//	return fallback
//}

type UserServiceServer struct {
	pb.UnimplementedUserServiceServer
	Idb db.IMemDB
	log *logger.CustomLogger
}

func NewUserServiceServer(Idb db.IMemDB, log *logger.CustomLogger) *UserServiceServer {
	return &UserServiceServer{Idb: Idb, log: log}
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

	GrpcMicroservice := NewUserServiceServer(memDB, customlogger)

	//Register DataService
	pb.RegisterUserServiceServer(grpcserver, GrpcMicroservice)
	log.Println("Starting grpc connection on port " + env.Server.GrpcPort)

	//startServing requests
	//go grpcserver.Serve(listen)
	grpcserver.Serve(listen)
}

func (s *UserServiceServer) GetLinkedAccountInfo(ctx context.Context, req *pb.Request) (*pb.Response, error) {
	// md, _ := metadata.FromIncomingContext(ctx)
	// token := md.Get("authorization")[0]
	// //set Headers
	// header := metadata.New(map[string]string{"authorization": token})
	// c := metadata.NewOutgoingContext(ctx, header)

	// client := grpc_client.ConnectToken(token_host, token_grpcPort)
	// _, err := client.VerifyToken(c, &pb_token.TokenRequest{Token: req.AccessToken})
	// if err != nil {
	// 	return nil, err
	// }

	data, err := s.Idb.GetLinkedAccountInfo(ctx, req)
	return data, err
}

func (s *UserServiceServer) UnlinkAccount(ctx context.Context, req *pb.Request) (*pb.Response, error) {
	// md, _ := metadata.FromIncomingContext(ctx)
	// token := md.Get("authorization")[0]
	// //set Headers
	// header := metadata.New(map[string]string{"authorization": token})
	// c := metadata.NewOutgoingContext(ctx, header)

	// client := grpc_client.ConnectToken(token_host, token_grpcPort)
	// _, err := client.VerifyToken(c, &pb_token.TokenRequest{Token: req.AccessToken})
	// if err != nil {
	// 	return nil, err
	// }

	data, err := s.Idb.UnlinkAccount(ctx, req)
	return data, err
}

func (s *UserServiceServer) GetUserInfoById(ctx context.Context, req *pb.Request) (*pb.Response, error) {
	// md, _ := metadata.FromIncomingContext(ctx)
	// token := md.Get("authorization")[0]
	// //set Headers
	// header := metadata.New(map[string]string{"authorization": token})
	// c := metadata.NewOutgoingContext(ctx, header)

	// client := grpc_client.ConnectToken(token_host, token_grpcPort)
	// _, err := client.VerifyToken(c, &pb_token.TokenRequest{Token: req.AccessToken})
	// if err != nil {
	// 	return nil, err
	// }

	data, err := s.Idb.GetAccountInfoById(ctx, req)
	return data, err
}

func (s *UserServiceServer) CheckAuthUserExists(ctx context.Context, req *pb.Request) (res *wrapperspb.BoolValue, err error) {
	data, err := s.Idb.CheckAuthUserExists(req)
	response := &wrapperspb.BoolValue{
		Value: data,
	}

	return response, err
}

func (s *UserServiceServer) GetAccountInfo(ctx context.Context, req *pb.Request) (*pb.Response, error) {
	data, err := s.Idb.GetAccountInfo(ctx, req)
	return data, err
}

func (s *UserServiceServer) CreateNewUser(ctx context.Context, req *pb.Request) (*pb.Response, error) {
	data, err := s.Idb.CreateNewUser(req)
	return data, err
}

func (s *UserServiceServer) SetAccountInfo(ctx context.Context, req *pb.Request) (*pb.Response, error) {
	// md, _ := metadata.FromIncomingContext(ctx)
	// token := md.Get("authorization")[0]
	// //set Headers
	// header := metadata.New(map[string]string{"authorization": token})
	// c := metadata.NewOutgoingContext(ctx, header)

	// client := grpc_client.ConnectToken(token_host, token_grpcPort)
	// _, err := client.VerifyToken(c, &pb_token.TokenRequest{Token: req.AccessToken})
	// if err != nil {
	// 	return nil, err
	// }

	data, err := s.Idb.SetAccountInfo(ctx, req)
	if err == nil {
		data, err = s.Idb.GetAccountInfoById(ctx, req)

	}
	return data, err
}

func (s *UserServiceServer) CheckIdpAccountLinked(ctx context.Context, req *pb.Request) (res *wrapperspb.BoolValue, err error) {
	data, err := s.Idb.CheckIdpAccountLinked(req)
	response := &wrapperspb.BoolValue{
		Value: data,
	}

	return response, err
}

func (s *UserServiceServer) MigrateAccount(ctx context.Context, req *pb.Request) (*pb.Response, error) {
	data, err := s.Idb.MigrateAccount(req)
	return data, err
}