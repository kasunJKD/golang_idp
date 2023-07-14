package jwt

import (
	"crypto/hmac"
	"crypto/sha256"
	"encoding/hex"
	"math/rand"
	"time"

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


