package main

import (
	"context"
	"idp-repository/env"
	"idp-repository/pkg/logger"
	rd "idp-repository/pkg/redis"
	pb "idp-repository/protos/token"
	"idp-repository/token-repository/internal/jwt"
	"idp-repository/token-repository/internal/redis"
	"log"
	"net"
	"strings"

	grpc "google.golang.org/grpc"
	"google.golang.org/protobuf/types/known/wrapperspb"
)


//var (
//	grpcPort = getEnv("GRPC_PORT", ":8082")
//	redisHost = getEnv("REDIS_HOST", "membership-redis")
//	redisPort = getEnv("REDIS_PORT", ":6379")
//	password = getEnv("PASSWORD", "9221")
//	mem_host = getEnv("MEM_HOST", "membership-db")
//	mem_port = getEnv("MEM_PORT", "5432")
//	user = getEnv("DB_USER", "postgres")
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

type TokenServiceServer struct {
	pb.UnimplementedTokenServiceServer
	Ird redis.IAuthRedis
	log *logger.CustomLogger
}

func NewTokenServiceServer(Ird redis.IAuthRedis, log *logger.CustomLogger) *TokenServiceServer {
	return &TokenServiceServer{Ird: Ird, log: log}
}

func main() {
	customLogger := logger.NewCustomLogger()

	customLogger.Infolog.Println("Welcome to the server")

	env , err := env.GetEnv()
	if err != nil {
		customLogger.Errorlog.Fatalf("Loading env: %v", err)
	}

	//connecting to redis
	redisConn := rd.Connect(env)

	defer redisConn.Close()

	//initialize redis
	authRedis := redis.NewRedisInstiate(redisConn, customLogger, env)

	//start listening for grpc
	listen, err := net.Listen("tcp", env.Server.GrpcPort)
	if err != nil {
		log.Fatal(err)
	}

	grpcserver := grpc.NewServer()

	GrpcMicroservice := NewTokenServiceServer(authRedis, customLogger)

	//Register DataService
	pb.RegisterTokenServiceServer(grpcserver, GrpcMicroservice)
	log.Println("Starting grpc connection on port " + env.Server.GrpcPort)

	//startServing requests
	//go grpcserver.Serve(listen)
	grpcserver.Serve(listen)

}

func (s *TokenServiceServer) NewAuthCodeToken(ctx context.Context, req *pb.TokenRequest) (*pb.AuthCodeToken, error) {
	data, err := jwt.NewAuthCodeToken(req.GetCode(), req.GetRefreshToken(), req.GetRedirectURI())
	return data, err
}

func (s *TokenServiceServer) NewAuthCodeGrant(ctx context.Context, req *pb.TokenRequest) (*wrapperspb.StringValue, error) {
	data := jwt.NewAuthCodeGrant(ctx, req.GetRedirectURI())
	response := &wrapperspb.StringValue{
		Value: data,
	}
	return response, nil
}

func (s *TokenServiceServer) VerifyAuthCodeToken(ctx context.Context, req *pb.TokenRequest) (*wrapperspb.BoolValue, error) {
	data := jwt.VerifyAuthCodeToken(req.GetToken())
	response := &wrapperspb.BoolValue{
		Value: data,
	}
	return response, nil
}

func (s *TokenServiceServer) NewAuthCodeRefreshToken(ctx context.Context, req *pb.RefreshTokenRequest) (*pb.AuthCodeToken, error) {
	data, err := jwt.NewAuthCodeRefreshToken(ctx, req.GetRefreshToken())
	return data, err
}

func (s *TokenServiceServer) AuthCodeRefreshTokenExists(ctx context.Context, req *pb.RefreshTokenRequest) (*wrapperspb.BoolValue, error) {
	data := jwt.AuthCodeRefreshTokenExists(req.GetRefreshToken(), req.GetInvalidateIfFound())
	response := &wrapperspb.BoolValue{
		Value: data,
	}
	return response, nil
}


func (s *TokenServiceServer) AddUserIdAuthCodeFlow(ctx context.Context, req *pb.User) (*wrapperspb.BoolValue, error) {
	data, err := jwt.AddUserIdAuthCodeFlow(req.GetUserId())
	response := &wrapperspb.BoolValue{
		Value: data,
	}
	return response, err
}

func (s *TokenServiceServer) GetUserIdfromAccesstoken(ctx context.Context, req *pb.User) (*wrapperspb.StringValue, error) {
	data, err := jwt.GetUserIdfromAccesstoken(req.GetAccessToken())
	response := &wrapperspb.StringValue{
		Value: data,
	}
	return response, err
}

func (s *TokenServiceServer) CreateToken(ctx context.Context, req *pb.User) (*pb.AuthCodeToken, error) {
	err := s.Ird.SetToken(req)
	return &pb.AuthCodeToken{}, err
}

func (s *TokenServiceServer) VerifyToken(ctx context.Context, req *pb.TokenRequest) (*pb.AuthCodeToken, error) {
	_, err := s.Ird.GetToken(req.Token)
	return &pb.AuthCodeToken{}, err
}

func (s *TokenServiceServer) RefreshToken(ctx context.Context, req *pb.TokenRequest) (*pb.AuthCodeToken, error) {
	_, err := s.Ird.GetToken(req.RefreshToken)
	if err == nil {
		//Revoke the current refreshToken as we will create and return a new one
		s.log.Infolog.Println("revoke the old refreshToken..")
		err = s.Ird.DeleteToken(req.RefreshToken)
	}
	return &pb.AuthCodeToken{}, err
}

func (s *TokenServiceServer) RevokeToken(ctx context.Context, req *pb.TokenRequest) (*pb.AuthCodeToken, error) {
	//Get the refresh token from redis to revoke it as well
	refreshToken, err := s.Ird.GetToken(req.Token)
	if err != nil {
		s.log.Infolog.Println("revoke refreshToken..")
		err = s.Ird.DeleteToken(refreshToken)
	}
	s.log.Infolog.Println("revoke accessToken..")
	err = s.Ird.DeleteToken(req.Token)
	return &pb.AuthCodeToken{}, err
}

func (s *TokenServiceServer) RevokeAll(ctx context.Context, req *pb.User) (*pb.AuthCodeToken, error) {
	val, err := s.Ird.GetToken(req.UserId)
	if err == nil {
		s.log.Infolog.Println("No Active tokens found for the user..")
	} else {
		tokens_list := strings.Split(val, ", ")

		s.log.Infolog.Printf("User's: \"%s\" list of active tokens: %s\n", req.UserId, tokens_list)

		for _, token := range tokens_list {
			err = s.Ird.DeleteToken(token)
		}
		err = s.Ird.DeleteToken(req.UserId)
	}

	return &pb.AuthCodeToken{}, err
}

func (s *TokenServiceServer) CreateResetPasswordToken(ctx context.Context, req *pb.User) (*pb.AuthCodeToken, error) {
	err := s.Ird.SetToken(req)
	return &pb.AuthCodeToken{}, err
}
