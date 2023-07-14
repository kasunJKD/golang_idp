package server

import (
	"context"
	"github.com/go-redis/redis/v8"
	"google.golang.org/grpc"
	"log"
	redisFunctions "login-throttles/internal/redis"
	"login-throttles/internal/service"
	rd "login-throttles/pkg/redis"
	pb "login-throttles/protos/core"
	"net"
	"os"
	"os/signal"
)

//RunServer run server
func RunServer(ctx context.Context, port string) error {
	listen, err := net.Listen("tcp", port)
	if err != nil {
		return err
	}
	serv := service.NewCoreService()

	server := grpc.NewServer()
	pb.RegisterCoreServiceServer(server, serv)

	c := make(chan os.Signal, 1)
	signal.Notify(c, os.Interrupt)

	//connecting to redis
	var redisConn *redis.Client
	redisConn = rd.Connect()
	redisFunctions.RedisConn = redisConn
	defer redisConn.Close()

	go func() {
		for range c {
			log.Println("shutting down gRPC server...")

			server.GracefulStop()

			<-ctx.Done()
		}
	}()

	log.Println("Start service otp" + port + " ...")
	return server.Serve(listen)
}
