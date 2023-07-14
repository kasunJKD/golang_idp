package redis

import (
	"context"
	"fmt"
	pb "idp-repository/protos/token"
	"time"
)

func (c RedisConfig) SetToken(req *pb.User) error {
	redisExp := 7 * 24 * time.Hour

	//Store the accessToken in redis
	redisCmd := c.rd.Set(context.Background(), req.GetAccessToken(), req.GetRefreshToken(), redisExp)

	if redisCmd.Err() != nil {
		if c.env.Server.Mode == "debug" {
			c.log.Debuglog.Println(redisCmd.Err())
		}
		c.log.Errorlog.Fatal("Error adding the token to redis")
	}

	//Store the refreshToken in redis
	redisCmd = c.rd.Set(context.Background(), req.GetRefreshToken(), req.GetUserId(), redisExp)

	if redisCmd.Err() != nil {
		if c.env.Server.Mode == "debug" {
			c.log.Debuglog.Println(redisCmd.Err())
		}
		c.log.Errorlog.Fatal("Error adding the token to redis")
	}

	//Get the list of the all access & refresh tokens for the user
	val, redis_err := c.rd.Get(context.Background(), req.GetUserId()).Result()

	if redis_err != nil || val == "" {
		//Create a list for the user's access & refresh tokens
		val = req.GetAccessToken() + ", " + req.GetRefreshToken()
		redisCmd = c.rd.Set(context.Background(), req.GetUserId(), val, 0)
	} else {
		//Add the new accessToken & refreshTokens to the user's tokens list
		val += ", " + req.GetAccessToken() + ", " + req.GetRefreshToken()
		redisCmd = c.rd.Set(context.Background(), req.GetUserId(), val, 0)
	}
	if redisCmd.Err() != nil {
		if c.env.Server.Mode == "debug" {
			c.log.Debuglog.Println(redisCmd.Err())
		}
		c.log.Errorlog.Fatal("Error adding the token to redis")
	}
	return nil
}

func (c RedisConfig) GetToken(key string) (string, error) {
	//Check in redis if the token is revoked
	val, redis_err := c.rd.Get(context.Background(), key).Result()

	if redis_err != nil || val == "" {
		c.log.Errorlog.Println("token is revoked; not found in redis")
		return "", fmt.Errorf("token is revoked")
	}
	return val, nil
}

func (c RedisConfig) DeleteToken(token string) error {
	//revoke token in redis
	redisCmd := c.rd.Del(context.Background(), token)

	if redisCmd.Err() != nil {
		c.log.Errorlog.Println("Error deleting the token from redis")
		if c.env.Server.Mode == "debug" {
			c.log.Debuglog.Println(redisCmd.Err())
		}

		return fmt.Errorf("error deleting the token from redis")
	}
	c.log.Infolog.Printf("token: \"%s\" is revoked", token)

	return nil
}
