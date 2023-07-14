package crypto

import (
	"context"
	"crypto/hmac"
	"crypto/sha256"
	"encoding/hex"
	"math/rand"
	"os"
	"time"

	"idp-repository/pkg/grpc_client"
	pb "idp-repository/protos/login"

)

var Key = []byte("my_secret_key")

// Hashes the string using SHA-256
func hash(str string) string {
	token := hmac.New(sha256.New, Key)
	token.Write([]byte(str))
	// hasher := sha256.New()
	// hasher.Write([]byte(str))
	macSum := token.Sum(nil)
	return hex.EncodeToString(macSum)
}

const src = "abcdefghijklmnopqrstuvwxyz0123456789ABCDEFGHIJKLMNOPQRSTUVWXYZ123456789"
const srcS = "abcdefghijklmnopqrstuvwxyz0123456789ABCDEFGHIJKLMNOPQRSTUVWXYZ123456789-_"

// Generates a string of given length filled with random bytes
func generateNonce(n int) string {
	if n < 1 {
		return ""
	}

	b := make([]byte, n)
	srcLen := int64(len(src))

	for i := range b {
		b[i] = src[rand.Int63()%srcLen]
	}

	return string(b)
}

func init() {
	// Seeding the random package
	rand.Seed(time.Now().UnixNano())
}

func NewClientId() string {
	var client_id string
	var exists bool = true
	//check for duplicates
	for exists {
		client_id = generateNonce(40)
		//grpc client connection----------------------------
		login_client := grpc_client.ConnectLogin(os.Getenv("HOST"), os.Getenv("GRPC_PORT"))
		//create new client service
		rr := &pb.ClientReq{
			ClientId: client_id,
		}
		//calling ValidateClientId
		res, err := login_client.ValidateClientId(context.Background(), rr)

		if err != nil {
			panic(err)
		}

		exists = res.Value
	}
	return client_id
}

func NewClientSecret() string {
	secret := generateCryptoRandom(srcS, 32)
	return secret
}

func generateCryptoRandom(chars string, length int32) string {
	//reference https://earthly.dev/blog/cryptography-encryption-in-go/
	bytes := make([]byte, length)
	rand.Read(bytes)
	for index, element := range bytes {
		randomize := element%byte(len(chars))
		bytes[index] = chars[randomize]
	   }
	  
	return string(bytes)
}

