package env

import (
	"log"
	"os"
)

// App env struct
type Env struct {
	Server   Server
}

type Server struct {
	AppVersion        string
	GrpcPort          string
	HttpPort          string
	Host          	  string
	RepoHost          string
	TokenService      string
	UserService       string
	LoginService       string
	OtpService       string
	Mode              string
	// ReadTimeout       time.Duration
	// MaxConnectionIdle time.Duration
	// Timeout           time.Duration
	// MaxConnectionAge  time.Duration
	// Time              time.Duration
}

// GetEnv returns a Env struct
func GetEnv() (*Env, error) {
	env := &Env{}

	env.Server.AppVersion = "1.0.0"
	env.Server.GrpcPort = getEnv("GRPC_PORT", ":81")
	env.Server.HttpPort = getEnv("HTTP_PORT", ":80")
	env.Server.Host = getEnv("HOST", "localhost")
	env.Server.RepoHost = getEnv("REPO_HOST", "localhost")
	env.Server.TokenService = getEnv("TOKEN_SERVICE", "localhost")
	env.Server.UserService = getEnv("USER_SERVICE", "localhost")
	env.Server.LoginService = getEnv("LOGIN_SERVICE", "localhost")
	env.Server.OtpService = getEnv("OTP_SERVICE", "localhost")
	env.Server.Mode = getEnv("MODE", "Development")


	return env, nil
}

func getEnv(key, fallback string) string {
	value, found := os.LookupEnv(key)
	if found {
		return value
	}
	log.Println("Key not found: ", key)
	os.Setenv(key, fallback)
	return fallback
}
