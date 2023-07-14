package redis

import (
	"idp-repository/env"
	"idp-repository/pkg/logger"
	rd "idp-repository/pkg/redis"
	pb "idp-repository/protos/token"
	"log"
	"os"
	"testing"
)

func TestRedisFunctions(t *testing.T) {
	os.Setenv("REDIS_PORT", "46379")
	os.Setenv("REDIS_PASS", "9221")

	customLogger := logger.NewCustomLogger()
	env , err := env.GetEnv()
	if err != nil {
		customLogger.Errorlog.Fatalf("Loading env: %v", err)
	}
	redisConn := rd.Connect(env)
	defer redisConn.Close()

	c := RedisConfig{rd: redisConn, log: customLogger, env: nil}

	log.Println("Testing token redis functions ------------>")

	//Test case #1 (check set token)
	log.Println("Test case #1: set token ------------>")
	req := &pb.User{
		UserId: "TestUserId",
		AccessToken : "TestAccessToken",
		RefreshToken: "TestRefreshToken",
	}

	err = c.SetToken(req)

	if err != nil {
		t.Fatalf("error test case #1: %v", err)
	} else {
		log.Println("test case #1 is passed")
	}

	//Test case #2 (check get token)
	log.Println("Test case #2: get token ------------>")
	AccessToken := "TestAccessToken"

	_, err = c.GetToken(AccessToken)

	if err != nil {
		t.Fatalf("error test case #2: %v", err)
	} else {
		log.Println("test case #2 is passed")
	}

	//Test case #3 (check delete token)
	log.Println("Test case #3: delete token ------------>")

	err = c.DeleteToken(AccessToken)

	if err != nil {
		t.Fatalf("error test case #3: %v", err)
	} else {
		log.Println("test case #3 is passed")
	}

	log.Println("All test cases are finished.")

}