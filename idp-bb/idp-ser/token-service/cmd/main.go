package main

import (
	"context"
	"idp-service/env"
	"idp-service/pkg/logger"
	cmd "idp-service/token-service/server"
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
		customlogger.Infolog.Fatal("run server err: ", err)
		os.Exit(1)
	}
}