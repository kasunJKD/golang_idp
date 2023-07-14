package env

import (
	"log"
	"os"
)

// App env struct
type Env struct {
	Server   Server
	Postgres Postgres
	Redis Redis
}

type Server struct {
	AppVersion        string
	GrpcPort          string
	Host              string
	Mode              string
	// ReadTimeout       time.Duration
	// MaxConnectionIdle time.Duration
	// Timeout           time.Duration
	// MaxConnectionAge  time.Duration
	// Time              time.Duration
}

// Postgresql
type Postgres struct {
	Host     string
	Port     string
	User     string
	Password string
	Dbname   string
	SSLMode  bool
	Driver string
	
}

// Redis
type Redis struct {
	Host     string
	Port     string
	Password string
}

// GetEnv returns a Env struct
func GetEnv() (*Env, error) {
	env := &Env{}

	env.Server.AppVersion = "1.0.0"
	env.Server.GrpcPort = getEnv("GRPC_PORT", ":81")
	env.Server.Host = getEnv("HOST", "localhost")
	env.Server.Mode = getEnv("MODE", "Development")

	env.Postgres.Host = getEnv("DB_HOST", "localhost")
	env.Postgres.Port = getEnv("DB_PORT", "5432")
	env.Postgres.User = getEnv("DB_USER", "user")
	env.Postgres.Password = getEnv("DB_PASS", "pass")
	env.Postgres.Dbname = getEnv("DB_NAME", "membership")
	env.Postgres.SSLMode = false
	env.Postgres.Driver = "pgx"

	env.Redis.Host = getEnv("REDIS_HOST", "localhost")
	env.Redis.Port = getEnv("REDIS_PORT", "6379")
	env.Redis.Password = getEnv("REDIS_PASS", "pass")

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
