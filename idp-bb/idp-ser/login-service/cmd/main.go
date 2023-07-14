package main

import (
	"context"
	"idp-service/env"
	cmd "idp-service/login-service/server"
	"idp-service/pkg/logger"
	"os"
)

func main() {
	customlogger := logger.NewCustomLogger()

	env , err := env.GetEnv()
	if err != nil {
		customlogger.Errorlog.Fatalf("Loading env: %v", err)
	}

	s := cmd.NewServer(customlogger, env)
	if err := s.RunServer(context.Background(), env.Server.GrpcPort); err != nil {
		customlogger.Errorlog.Fatal("run server err: ", err)
		os.Exit(1)
	}
}