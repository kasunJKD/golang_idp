package server

import (
	"context"
	"idp-service/env"
	"idp-service/pkg/grpc_client"
	"idp-service/pkg/logger"
	pb "idp-service/protos/user"
	"idp-service/user-service/internal/service"
	"log"
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

	userClient, err := grpc_client.NewUserServiceClient(s.env.Server.RepoHost, s.env.Server.GrpcPort)
	if err != nil {
		log.Println("user client creation failed")
		return err
	}

	defer userClient.Conn.Close()

	tokenClient, err := grpc_client.NewTokenServiceClient(s.env.Server.TokenService, s.env.Server.GrpcPort)
	if err != nil {
		log.Println("token client creation failed")
		return err
	}

	defer tokenClient.Conn.Close()

	serv := service.NewUserService(userClient,tokenClient, s.ll, s.env)
	
	server := grpc.NewServer()
	pb.RegisterUserServiceServer(server, serv)

	c := make(chan os.Signal, 1)
	signal.Notify(c, os.Interrupt)

	go func() {
		for range c {
			log.Println("shutting down gRPC server...")

			server.GracefulStop()

			<-ctx.Done()
		}
	}()

	log.Println("Start service Login" + port + " ...")
	return server.Serve(listen)
}