package redis

import (
	"github.com/go-redis/redis/v8"
	"idp-repository/env"
	"idp-repository/pkg/logger"
	pb "idp-repository/protos/token"
)

type RedisConfig struct {
	rd *redis.Client
	log *logger.CustomLogger
	env *env.Env
}

func NewRedisInstiate(rd *redis.Client, log *logger.CustomLogger, env *env.Env) *RedisConfig {
	return &RedisConfig{rd: rd, log: log, env: env}
}

type IAuthRedis interface {
	SetToken(req *pb.User) error
	GetToken(key string) (string, error)
	DeleteToken(token string) error
}