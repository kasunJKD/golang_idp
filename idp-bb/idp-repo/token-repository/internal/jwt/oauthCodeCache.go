package jwt

import (
	"encoding/json"
	"fmt"
	"log"
	"strconv"
	"context"
	"time"

	"idp-repository/token-repository/internal/utils"
	"idp-repository/pkg/redis"
	pb "idp-repository/protos/token"

)

const (
	// Redis HSET which holds the issued tokens
	authCodeTokensSet = "OA2B_AC_Tokens"

	// Redis HSET which holds the issued grants until a token request is made.
	authCodeGrantSet = "OA2B_AC_Grants"

	// AuthCodeFlowID is prepended to a refresh token issued by the Authorization Code flow
	AuthCodeFlowID = "AUTHCODE"

	//auth userID
	AuthCodeFlowUserId = "UserIdAuthCodeFlow"
)

// AuthCodeToken represents a token issued by the Authorization Code flow
// https://tools.ietf.org/html/rfc6749#section-4.1.3
// type AuthCodeToken struct {
// 	AccessToken  string `json:"access_token"`
// 	RefreshToken string `json:"refresh_token"`
// 	ExpiresIn    int    `json:"expires_in"`
// }

// Holds the meta data of an access token
type authCodeTokenMeta struct {
	AuthGrant    string    `json:"auth_grant"`
	CreationTime time.Time `json:"creation_time"`
	Nonce        string    `json:"nonce"`
}

// Holds the token as well as its metadata.
// It is the internal representation of the token inside the Redis cache.
type internalAuthCodeToken struct {
	Token pb.AuthCodeToken     `json:"token"`
	Meta  authCodeTokenMeta `json:"meta"`
}

// NewAuthCodeToken issues new access tokens for the Authorization Code flow.
// It searches for 'code' in the Redis cache and throws errors if not found.
// If found, it checks if it has crossed is expiry limit which is 10 minutes.
// If crossed, an error is thrown.
// Else a new token is generated and returned.
// Refer RFC 6749 Section 4.1.2 (https://tools.ietf.org/html/rfc6749#section-4.1.2)
func NewAuthCodeToken(code, refreshToken, redirectURI string) (*pb.AuthCodeToken, error) {
	// First check if such an authorization grant has been issued
	rediscon := redis.CreateClient(0)
	defer rediscon.Close()
	
	value := code + ":" + redirectURI
	reply := rediscon.HExists(context.Background(), authCodeGrantSet, value)
	reply_value, err := utils.Bool2int(reply.Val())
	if err != nil {
		log.Println("NewAuthCodeToken: " + err.Error())
		return nil, err
	}

	// If 'value' is not found in the Redis cache, there are the following possibilites:
	// - A token was already issued on this authorization grant and must be revoked.
	// - It has expired and was removed by housekeep().
	// - It was never issued.
	// - the redirect URI is wrong
	if reply_value == 0 {
		return nil, fmt.Errorf("recycled/expired/invalid authorization grant or wrong redirect_uri")
	}

	// If found, check if it has expired since housekeeping runs only every 5 minutes
	//intTime, err := redis.HGet(conn.Do("HGET", authCodeGrantSet, value))
	intTime, err := strconv.ParseInt(rediscon.HGet(context.Background(), authCodeGrantSet, value).Val(),10,64)
	if err != nil {
		log.Println("NewAuthCodeToken: " + err.Error())
		return nil, err
	}

	issueTime := time.Unix(intTime, 0)
	if time.Since(issueTime) >= 10*time.Minute {
		return nil, fmt.Errorf("expired authorization grant")
	}

	// If not expired, remove it from the Redis cache since
	// we're about to issue a token for it.
	go removeAuthCodeGrant(code, redirectURI)

	var token *pb.AuthCodeToken
	var meta *authCodeTokenMeta

	// Generates a new key if a duplicate is encountered
	for reply_value == 1 {
		token, meta = generateAuthCodeToken(code)

		// Replace newly-generated refresh token with function parameter 'refreshToken'
		// if it is of length 72 since SHA-256 generates a string of length 64 and we
		// prepend it with a flow identifier of length 8. (AUTHCODE)
		if len(refreshToken) == 72 {
			token.RefreshToken = refreshToken
		}

		reply = rediscon.HExists(context.Background(), authCodeTokensSet, token.AccessToken)
		reply_value, err = utils.Bool2int(reply.Val())
		if err != nil {
			log.Println(err)
			return nil, err
		}
	}

	jsonBytes, err := json.Marshal(internalAuthCodeToken{Token: *token, Meta: *meta})
	if err != nil {
		panic(err)
	}

	_ = rediscon.HSet(context.Background(),authCodeTokensSet, token.AccessToken, string(jsonBytes))
	if err != nil {
		return nil, err
	}

	return token, nil
}

func removeAuthCodeGrant(code, redirectURI string) {
	rediscon := redis.CreateClient(0)
	defer rediscon.Close()
	
	redisCmd := rediscon.HDel(context.Background(), authCodeGrantSet, code+":"+redirectURI)
	if redisCmd.Err() != nil {
		log.Println(redisCmd.Err())
	}
}

// Generates access and refresh tokens.
// Access token is a hex-encoded string of the SHA-256 hash of the concatenation of
// the code, time of creation and a nonce.
// Refresh token starts with the flow identifier "AUTHCODE" followed by a hex-encoded string of
// the SHA-256 hash of the concatenation of time of creation and the same nonce as above.
func generateAuthCodeToken(code string) (*pb.AuthCodeToken, *authCodeTokenMeta) {
	nonce := generateNonce(16)
	creationTime := time.Now()

	rediscon := redis.CreateClient(0)
	defer rediscon.Close()
	redisCmd := rediscon.Get(context.Background(), AuthCodeFlowUserId)
	// {
  //   "audience": "my_audience",
  //   "email": "my_email",
  //   "expires_in": 0,
  //   "issued_to": "my_issued_to",
  //   "scope": "my_scope",
  //   "user_id": "my_user_id",
  //   "verified_email": false
  // }

	accessToken := hash(fmt.Sprintf("code:%s,\ncreationTime:%s,\n nonce:%s,\n userid:%s", code, creationTime, nonce, redisCmd.Val()))
	refreshToken := hash(fmt.Sprintf("creationTime:%s,\n nonce:%s", creationTime, nonce))

	rediscon2 := redis.CreateClient(0)
	defer rediscon2.Close()
	cmdset := rediscon2.Set(context.Background(), accessToken, redisCmd.Val(), 0)
	if cmdset.Err() != nil {
		log.Println(cmdset.Err())
	}
	cmddel := rediscon2.Del(context.Background(), AuthCodeFlowUserId)
	if cmddel.Err() != nil {
		log.Println(cmddel.Err())
	}

	return &pb.AuthCodeToken{
			AccessToken:  accessToken,
			RefreshToken: refreshToken,
			ExpiresIn:    3600,
		}, &authCodeTokenMeta{
			AuthGrant:    code,
			CreationTime: creationTime,
			Nonce:        nonce,
		}
}

// NewAuthCodeGrant generates a new authorization grant and adds it to a Redis cache set.
// This function takes the redirect URI as an argument, since RFC 6749 requires the same URI
// to be used in the token request as was used in the authorization grant request, if any.
// Thus, we store it along with the authorization grant in order for us to verify it against
// the one sent in the token request.
// Refer: https://tools.ietf.org/html/rfc6749#section-4.1.3
func NewAuthCodeGrant(ctx context.Context, redirectURI string) string {
	var code string
	var reply int64 = 0
	// In case we get a duplicate value, we iterate until we get a unique one.
	rediscon := redis.CreateClient(0)
	defer rediscon.Close()

	for reply == 0 {
		code = generateNonce(20)
		value := code + ":" + redirectURI
		cmd := rediscon.HSet(ctx, authCodeGrantSet, value, time.Now().Unix())
		reply = cmd.Val()
		if cmd.Err() != nil {
			log.Println(cmd.Err())
		}
	}

	return code
}

// VerifyAuthCodeToken checks if the token exists in the Redis cache.
// Returns true if token found, false otherwise.
func VerifyAuthCodeToken(token string) bool {
	rediscon := redis.CreateClient(0)
	defer rediscon.Close()
	cmd := rediscon.HGet(context.Background(),authCodeTokensSet, token)
	return cmd.Err() == nil
}

// NewAuthCodeRefreshToken returns new token for the previously issued refresh token
// The refresh token is kept intact and can be used for future requests.
func NewAuthCodeRefreshToken(ctx context.Context, refreshToken string) (*pb.AuthCodeToken, error) {
	code := NewAuthCodeGrant(ctx, "")
	token, err := NewAuthCodeToken(code, refreshToken, "")
	if err != nil {
		return nil, err
	}

	return token, nil
}

func invalidateAuthCodeToken(accessToken string) {
	rediscon := redis.CreateClient(0)
	defer rediscon.Close()
	cmd := rediscon.HDel(context.Background(),authCodeTokensSet, accessToken)
	if cmd.Err() != nil {
		log.Println(cmd.Err())
	}
}

// AuthCodeRefreshTokenExists checks if the refresh token exists in the Redis cache
// and returns the appropriate boolean value.
// Params:
// refreshToken: the token to look for in the cache
// invalidateIfFound: if true, the token is invalidated if found
func AuthCodeRefreshTokenExists(refreshToken string, invalidateIfFound bool) bool {
	rediscon := redis.CreateClient(0)
	defer rediscon.Close()
	var token internalAuthCodeToken
	cmd := rediscon.HGetAll(context.Background(),authCodeTokensSet)
	items, err := utils.ByteSlices(cmd.Val(), cmd.Err())	
	if err != nil {
		log.Println(cmd.Err())
	}

	for i := 1; i < len(items); i += 2 {
		err := json.Unmarshal(items[i], &token)
		if err != nil {
			log.Println(err)
			break
		}

		if refreshToken == token.Token.RefreshToken {
			if invalidateIfFound {
				invalidateAuthCodeToken(token.Token.AccessToken)
			}

			return true
		}
	}

	return false
}

func AddUserIdAuthCodeFlow (userid string)(bool, error) {
	rediscon := redis.CreateClient(0)
	defer rediscon.Close()

	cmd := rediscon.Set(context.Background(), AuthCodeFlowUserId, userid, 0)
	if cmd.Err() != nil {
		log.Println(cmd.Err())
	}

	return cmd.Err() == nil, cmd.Err()
}

func GetUserIdfromAccesstoken (accesstoken string)(string, error) {
	rediscon := redis.CreateClient(0)
	defer rediscon.Close()

	cmd := rediscon.Get(context.Background(), accesstoken)
	if cmd.Err() != nil {
		log.Println(cmd.Err())
	}

	return cmd.Val(), nil
}


