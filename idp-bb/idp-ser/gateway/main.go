package main

import (
	"context"
	"idp-service/env"

	"google.golang.org/grpc/credentials/insecure"

	//"crypto/tls"
	"fmt"
	"idp-service/gateway/third_party"
	"idp-service/pkg/grpc_client"
	"idp-service/pkg/logger"
	"io/fs"
	"io/ioutil"
	"mime"
	"net/http"
	"os"
	"strings"

	"github.com/grpc-ecosystem/grpc-gateway/v2/runtime"
	"github.com/rs/cors"
	"google.golang.org/grpc"

	//"google.golang.org/grpc/credentials"
	"idp-service/gateway/handlers"

	"google.golang.org/grpc/grpclog"

	//"google.golang.org/grpc/credentials/insecure"
	//"membership/insecure"
	pblogin "idp-service/protos/login"
	pbuser "idp-service/protos/user"
)

// Run runs the gRPC-Gateway, dialling the provided address.
func main() {
	customlogger := logger.NewCustomLogger()

	env , err := env.GetEnv()
	if err != nil {
		customlogger.Errorlog.Fatalf("Loading env: %v", err)
	}

	// Adds gRPC internal logs. This is quite verbose, so adjust as desired!
	log := grpclog.NewLoggerV2(os.Stdout, ioutil.Discard, ioutil.Discard)
	grpclog.SetLoggerV2(log)

	gwmux := runtime.NewServeMux()
	err = pblogin.RegisterLoginServiceHandlerFromEndpoint(context.Background(), gwmux, env.Server.LoginService + env.Server.GrpcPort, []grpc.DialOption{grpc.WithBlock(), grpc.WithTransportCredentials(insecure.NewCredentials())})
	if err != nil {
		fmt.Errorf("failed to register login gateway: %w", err)
	}
	err = pbuser.RegisterUserServiceHandlerFromEndpoint(context.Background(), gwmux, env.Server.UserService + env.Server.GrpcPort, []grpc.DialOption{grpc.WithBlock(), grpc.WithTransportCredentials(insecure.NewCredentials())})
	if err != nil {
		fmt.Errorf("failed to register user gateway: %w", err)
	}

	oa := getOpenAPIHandler()

	c := cors.New(cors.Options {
		AllowedOrigins: []string{"*"},
		AllowedMethods: []string{"POST", "GET"},
		AllowCredentials: true,
		AllowedHeaders:[]string{"*"},
		Debug: true,
	})

	//Handler dependecies 
	tokenClient, err := grpc_client.NewTokenServiceClient(env.Server.TokenService, env.Server.GrpcPort)
	if err != nil {
		customlogger.Errorlog.Println("token client creation failed")
	}

	defer tokenClient.Conn.Close()
	
	loginClient, err := grpc_client.NewLoginServiceClient(env.Server.LoginService, env.Server.GrpcPort)
	if err != nil {
		customlogger.Errorlog.Println("login client creation failed")
	}

	defer loginClient.Conn.Close()

	userClient, err := grpc_client.NewUserServiceClient(env.Server.UserService, env.Server.GrpcPort)
	if err != nil {
		customlogger.Errorlog.Println("user client creation failed")
	}

	defer userClient.Conn.Close()

	newHandler := handlers.NewHandlerFunc(loginClient, tokenClient, userClient, customlogger, env)

	tt := http.NewServeMux()
	//tt.HandleFunc("/membership/authorize", handlers.HandleAuth)
	tt.HandleFunc("/membership/response", newHandler.HandleResponse)
	tt.HandleFunc("/membership/token", newHandler.HandleToken)
	//tt.HandleFunc("/membership/signin/oauth", handlers.HandleOauthSignIn)
	tt.HandleFunc("/membership/signin/response", newHandler.HandleOauthSignInResponse)
	tt.HandleFunc("/membership/oauth2/userinfo", newHandler.HandleUserInfo)
	tt.HandleFunc("/membership/otpLogin/response", newHandler.HandleOauthOtpLoginResponse)

	port := os.Getenv("PORT")
	if port == "" {
		port =  env.Server.HttpPort[1:]
	}
	gatewayAddr := ":" + port
	handler := http.HandlerFunc(func(w http.ResponseWriter, r *http.Request) {
		isSwagger := true
		if strings.HasPrefix(r.URL.Path, "/api") {
			gwmux.ServeHTTP(w, r)
			isSwagger = false
		}
		if strings.HasPrefix(r.URL.Path, "/membership") {
			 tt.ServeHTTP(w, r)
			 isSwagger = false
		}
		if isSwagger {
			oa.ServeHTTP(w, r)
		}
	})
	http.ListenAndServe(gatewayAddr, c.Handler(handler))
}


// getOpenAPIHandler serves an OpenAPI UI.
// Adapted from https://github.com/philips/grpc-gateway-example/blob/a269bcb5931ca92be0ceae6130ac27ae89582ecc/cmd/serve.go#L63
func getOpenAPIHandler() http.Handler {
	mime.AddExtensionType(".svg", "image/svg+xml")
	// Use subdirectory in embedded files
	subFS, err := fs.Sub(third_party.OpenAPI, "OpenAPI")
	if err != nil {
		panic("couldn't create sub filesystem: " + err.Error())
	}
	return http.FileServer(http.FS(subFS))
}