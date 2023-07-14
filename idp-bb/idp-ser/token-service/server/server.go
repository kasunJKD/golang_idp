package server

import (
	"context"
	"idp-service/env"
	"idp-service/pkg/grpc_client"
	"idp-service/pkg/logger"
	pb "idp-service/protos/token"
	"idp-service/token-service/internal/jwt"
	"idp-service/token-service/internal/service"
	"net"
	"os"
	"os/signal"

	"google.golang.org/grpc"
)

type Server struct {
	ll         *logger.CustomLogger
	env 	   *env.Env
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

	tokenClient, err := grpc_client.NewTokenServiceClient(s.env.Server.RepoHost, s.env.Server.GrpcPort)
	if err != nil {
		if s.env.Server.Mode == "debug" {
			s.ll.Debuglog.Printf("token client err %v", err)
		}
		s.ll.Errorlog.Println("token client creation failed")
		return err
	}

	defer tokenClient.Conn.Close()

	jwtInterface := jwt.NewJWTFunc(tokenClient, s.ll, s.env)
	serv := service.NewTokenService(tokenClient, s.ll, jwtInterface, s.env)
	
	server := grpc.NewServer()
	pb.RegisterTokenServiceServer(server, serv)

	c := make(chan os.Signal, 1)
	signal.Notify(c, os.Interrupt)

	go func() {
		for range c {
			s.ll.Infolog.Println("shutting down gRPC server...")

			server.GracefulStop()

			<-ctx.Done()
		}
	}()

	s.ll.Infolog.Println("Start service token" + port + " ...")
	return server.Serve(listen)
}