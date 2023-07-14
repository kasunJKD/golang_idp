package main

import (
	"context"
	"log"
	cmd "login-throttles/server"
	"os"
)

func main() {
	if err := cmd.RunServer(context.Background(), os.Getenv("GRPC_PORT")); err != nil {
		log.Fatal("run server err: ", err)
		os.Exit(1)
	}
}
