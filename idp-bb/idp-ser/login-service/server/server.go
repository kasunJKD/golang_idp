package server

import (
	"context"
	"idp-service/env"
	"idp-service/login-service/internal/auth"
	"idp-service/login-service/internal/service"
	"idp-service/pkg/grpc_client"
	"idp-service/pkg/logger"
	pb "idp-service/protos/login"
	"net"
	"os"
	"os/signal"

	"google.golang.org/grpc"
)

// Server
type Server struct {
	ll         *logger.CustomLogger
	env		   *env.Env
}

// Server constructor
func NewServer(ll *logger.CustomLogger, env *env.Env) *Server {
	return &Server{ll: ll, env: env}
}

//RunServer run server
func (s *Server) RunServer(ctx context.Context, port string) error {
	listen, err := net.Listen("tcp", port)
	if err != nil {
		return err
	}

	//Client connection intialisation===================> start
	tokenClient, err := grpc_client.NewTokenServiceClient(s.env.Server.TokenService, s.env.Server.GrpcPort)
	if err != nil {
		if s.env.Server.Mode == "debug" {
			s.ll.Debuglog.Printf("token client err %v", err)
		}
		s.ll.Errorlog.Println("token client creation failed")
		return err
	}
	s.ll.Infolog.Println("Token RPC Client Connection established")

	defer tokenClient.Conn.Close()
	
	loginClient, err := grpc_client.NewLoginServiceClient(s.env.Server.RepoHost, s.env.Server.GrpcPort)
	if err != nil {
		if s.env.Server.Mode == "debug" {
			s.ll.Debuglog.Printf("login client err %v", err)
		}
		s.ll.Errorlog.Println("login client creation failed")
		return err
	}
	s.ll.Infolog.Println("Login RPC Client Connection established")

	defer loginClient.Conn.Close()

	userClient, err := grpc_client.NewUserServiceClient(s.env.Server.UserService, s.env.Server.GrpcPort)
	if err != nil {
		if s.env.Server.Mode == "debug" {
			s.ll.Debuglog.Printf("user client err %v", err)
		}
		s.ll.Errorlog.Println("user client creation failed")
		return err
	}
	s.ll.Infolog.Println("User RPC Client Connection established")

	defer userClient.Conn.Close()

	otpClient, err := grpc_client.NewOTPServiceClient(s.env.Server.OtpService, s.env.Server.GrpcPort)
	if err != nil {
		if s.env.Server.Mode == "debug" {
			s.ll.Debuglog.Printf("otp client err %v", err)
		}
		s.ll.Errorlog.Println("user client creation failed")
		return err
	}
	s.ll.Infolog.Println("OTP RPC Client Connection established")

	defer userClient.Conn.Close()
	//Client connection intialisation===================> end

	iAuth := auth.NewAuthFunc(loginClient, tokenClient, userClient, s.ll, s.env)
	serv := service.NewLoginService(loginClient, tokenClient, otpClient, s.ll, iAuth, s.env)
	server := grpc.NewServer()
	pb.RegisterLoginServiceServer(server, serv)

	c := make(chan os.Signal, 1)
	signal.Notify(c, os.Interrupt)

	go func() {
		for range c {
			s.ll.Errorlog.Println("shutting down gRPC server...")

			server.GracefulStop()

			<-ctx.Done()
		}
	}()

	s.ll.Infolog.Println("Start Login service" + port + " ...")
	return server.Serve(listen)
}