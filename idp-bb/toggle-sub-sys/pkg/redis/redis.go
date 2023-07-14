package redis

import (
	"context"
	"github.com/go-redis/redis/v8"
	"log"
	"os"
)

type redisConfig struct {
	RedisClient *redis.Client
}

func Connect() *redis.Client {

	log.Printf("Connecting to Redis\n")

	redisConn := redis.NewClient(&redis.Options{
		Addr: os.Getenv("REDIS_HOST") + os.Getenv("REDIS_PORT"),
		Password: os.Getenv("PASSWORD"),
		DB: 1,
	})

	_, err := redisConn.Ping(context.Background()).Result()

	if err != nil {
		log.Fatal("Error connecting to redis: ", err)
	}

	return redisConn

}

func CreateClient() *redis.Client {
	rdb := redis.NewClient(&redis.Options{
		Addr: os.Getenv("REDIS_HOST") + os.Getenv("REDIS_PORT"),
		Password: os.Getenv("PASSWORD"),
		DB: 1,
	})

	return rdb
}
