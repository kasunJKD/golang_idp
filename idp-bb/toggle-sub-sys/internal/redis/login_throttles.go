package redis

import (
	"context"
	"fmt"
	"github.com/go-redis/redis/v8"
	"log"
	pb "login-throttles/protos/core"
	"strconv"
	"time"
)

var RedisConn *redis.Client

func CheckStatus(req *pb.Request) (res *pb.Response, err error) {
	//check status in redis
	key := "status" + ":" + req.Email
	val, redis_err := RedisConn.Get(context.Background(), key).Result()

	if redis_err != nil || val == "" {
		fmt.Println("active account; status not found in redis")
		return &pb.Response{
				Status: "active",
			}, nil
	}

	log.Println("inactive account; status is found in redis: ", val)

	return &pb.Response{
		Status: val,
	}, nil
}

func UpdateStatus(req *pb.Request) (res *pb.Response, err error) {
	//update status in redis
	key := "status" + ":" + req.Email

	if req.Status == "active"{
		redisCmd := RedisConn.Del(context.Background() ,key)
		if redisCmd.Err() != nil {
			log.Fatal("Error deleting the status from redis: ", redisCmd.Err())
			return nil, redisCmd.Err()
		}
	} else {
		var redisExp time.Duration = 0
		if req.Status == "locked" {
			redisExp = 1 * time.Hour
		}
		redisCmd := RedisConn.Set(context.Background(), key, req.Status, redisExp)
		if redisCmd.Err() != nil {
			log.Fatal("Error updating the status to redis: ", redisCmd.Err())
			return nil, redisCmd.Err()
		}
	}
	log.Println("account status updated in redis: ", req.Status)

	return &pb.Response{
		Status: "status updated",
	}, nil
}

func CheckLimitCounter(req *pb.Request) (res *pb.Response, err error) {
	//check limit counter in redis
	key := "passLogin" + ":" + req.Email
	if req.LoginType == "otpLogin"{
		key = "otpLogin" + ":" + req.UserId
	}
	val, redis_err := RedisConn.Get(context.Background(), key).Result()

	if redis_err != nil || val == "" {
		fmt.Println("no failed attempts recently; counter not found in redis")
		return &pb.Response{
			Counter: 0,
		}, nil
	}

	log.Println("number of failed attempts recently: ", val)

	intVal, err := strconv.Atoi(val)
	if err != nil {
		fmt.Println("Error during counter conversion: ", err)
		return nil, err
	}

	return &pb.Response{
		Counter: int32(intVal),
	}, nil
}

func UpdateLimitCounter(req *pb.Request) (res *pb.Response, err error) {
	//check status in redis
	redisExp := 1 * time.Hour
	key := "passLogin" + ":" + req.Email
	if req.LoginType == "otpLogin"{
		key = "otpLogin" + ":" + req.UserId
	}
	val, redis_err := RedisConn.Get(context.Background(), key).Result()

	if redis_err != nil || val == "" {
		fmt.Println("no failed attempts recently; counter not found in redis")
		redisCmd := RedisConn.Set(context.Background(), key, 1, redisExp)
		if redisCmd.Err() != nil {
			log.Fatal("Error updating the counter to redis: ", redisCmd.Err())
			return nil, redisCmd.Err()
		}
		return &pb.Response{
			Counter: 1,
		}, nil
	}


	intVal, err := strconv.Atoi(val)

	counter := intVal+1

	log.Println("this is the failed attempt number: ", counter)

	if counter >= 5 && req.LoginType == "passLogin" {
		log.Println("account status will be updated to locked for an hour")
		UpdateStatus(&pb.Request{Email: req.Email, Status: "locked"})
	}

	return &pb.Response{
		Counter: int32(counter),
	}, nil
}
