package jwt

import (
	"context"
	"fmt"
	"github.com/golang-jwt/jwt/v4"
	"google.golang.org/grpc/codes"
	"google.golang.org/grpc/metadata"
	"google.golang.org/grpc/status"
	"idp-service/env"
	"idp-service/pkg/grpc_client"
	"idp-service/pkg/logger"
	pb "idp-service/protos/token"
	"time"
)

//TokenService struct
type jwtfunc struct {
	t *grpc_client.TokenServiceClient
	ll *logger.CustomLogger
	env *env.Env
}

//NewTokenService create service
func NewJWTFunc(t *grpc_client.TokenServiceClient, ll *logger.CustomLogger, env *env.Env) *jwtfunc {
	return &jwtfunc{
		t: t,
		ll: ll,
		env: env,
	}
}

type IjwtFunc interface {
	CreateToken(req *pb.User) (*pb.AuthCodeToken, error)
	VerifyToken(ctx context.Context, req *pb.TokenRequest) (*pb.AuthCodeToken, error)
	RefreshToken(ctx context.Context, req *pb.TokenRequest) (*pb.AuthCodeToken, error) 
	RevokeToken(ctx context.Context, req *pb.TokenRequest) (*pb.AuthCodeToken, error)
	RevokeAll(ctx context.Context, req *pb.User) (*pb.AuthCodeToken, error)
	CreateResetPasswordToken(ctx context.Context, req *pb.User) (*pb.AuthCodeToken, error)
}

var jwtKey = []byte("my_secret_key")

type Claims struct {
	LocalId string `json:"localId"`
	//jwt.StandardClaims // 'StandardClaims' is deprecated
	Issuer string `json:"iss,omitempty"`
	Subject string `json:"sub,omitempty"`
	Audience string `json:"aud,omitempty"`
	ExpiresAt int64 `json:"exp,omitempty"`
	NotBefore int64 `json:"nbf,omitempty"`
	IssuedAt int64 `json:"iat,omitempty"`
	Id string `json:"jti,omitempty"`
}

func (c Claims) Valid() error {
	//TODO implement me
	panic("implement me")
}

//type MD map[string][]string

type Metadata struct {
	MdMainAccessToken string
}

func (j *jwtfunc) CreateToken(req *pb.User) (*pb.AuthCodeToken, error) {
	j.ll.Infolog.Printf("jwt =>  CreateToken function called")
	expirationTime := time.Now().Add(5 * time.Minute)
	issuedTime := time.Now()

	//Create new accessToken
	claims := &Claims{
		LocalId: req.UserId,
		Issuer: "Ali",
		Subject: "Test",
		Audience: "P99",
		ExpiresAt: expirationTime.Unix(),
		NotBefore: 0,
		IssuedAt: issuedTime.Unix(),
		Id: "",
	}

	token := jwt.NewWithClaims(jwt.SigningMethodHS256, claims)

	tokenString, err := token.SignedString(jwtKey)

	if err != nil {
		if j.env.Server.Mode == "debug" {
			j.ll.Debuglog.Println(err)
		}
		j.ll.Errorlog.Fatal("InternalServerError")
	}

	//Create new refreshToken
	expirationTime = time.Now().Add(7 * 24 * time.Hour)

	claims = &Claims{
		LocalId: req.UserId,
		Issuer: "Ali",
		Subject: "Test",
		Audience: "P99",
		ExpiresAt: expirationTime.Unix(),
		NotBefore: 0,
		IssuedAt: issuedTime.Unix(),
		Id: "",
	}

	token = jwt.NewWithClaims(jwt.SigningMethodHS256, claims)

	refreshToken, err := token.SignedString(jwtKey)

	if err != nil {
		if j.env.Server.Mode == "debug" {
			j.ll.Debuglog.Println(err)
		}
		j.ll.Errorlog.Fatal("InternalServerError")
	}

	//store token in redis by the repo layer
	_, err = j.t.Client.CreateToken(context.Background(), &pb.User{UserId: req.UserId, AccessToken: tokenString, RefreshToken: refreshToken})
	if err != nil {
		j.ll.Errorlog.Println("create tokens failed")
		if j.env.Server.Mode == "debug" {
			j.ll.Debuglog.Println(err)
		}
		return nil, fmt.Errorf("create tokens failed")
	}

	expirationTime = time.Now().Add(5 * time.Minute)

	res := &pb.AuthCodeToken{
		AccessToken: tokenString,
		ExpiresIn:        int32(expirationTime.Unix()),
		RefreshToken: 	  refreshToken,
	}

	return res, err
}

func (j *jwtfunc) VerifyToken(ctx context.Context, req *pb.TokenRequest) (*pb.AuthCodeToken, error) {
	j.ll.Infolog.Printf("jwt =>  VerifyToken function called")
	//isValid := true
	AccessToken := req.Token
	claims := jwt.MapClaims{}

	md, ok := metadata.FromIncomingContext(ctx)
	if !ok {
		j.ll.Debuglog.Println("metadata is not provided")
	}

	found := false
	values := md["authorization"]
	// check the rest request header and grpc metadata
	if len(values) > 0 {
		AccessToken = values[0]
		found = true
	} else {
		//check the grpc metadata
		for key := range claims {
			if key == "authorization"  {
				AccessToken = md.Get("authorization")[0]
				found = true
				break
			}
		}
	}
	if !found {
		j.ll.Debuglog.Println("authorization token is not provided")
	}

	_, err:= jwt.ParseWithClaims(AccessToken, claims, func(t *jwt.Token) (interface{}, error) {
		return jwtKey, nil
	})

	//fmt.Println(token)
	//log.Println(AccessToken)

	if req.Code == "ResetPasswordToken" {
		if claims["sub"] != req.Code {
			return nil, fmt.Errorf("token is invalid")
		}
	}

	if err != nil {
		// Refresh Token
		//if err.Error() == "Token is expired" {
		//	res, err := CreateToken(req)
		//	return res, err
		//}
		j.ll.Errorlog.Println("Token is invalid")
		//isValid = false
	}


	if err == nil {
		//Check in redis if the accessToken is revoked
		_, err = j.t.Client.VerifyToken(context.Background(), &pb.TokenRequest{Token: AccessToken})
		if err != nil {
			return nil, err
		}
	}

	// Identity Authentication
	//if req.LocalId != "" {
	//	if req.LocalId != claims["localId"] {
	//		fmt.Printf("Token is invalid; req.LocalId: %v, id not equal to jwtLocalId: %v\n", req.LocalId, claims["localId"])
	//		err = fmt.Errorf("identity authentication failed")
	//		//isValid = false
	//	}
	//}

	userId := fmt.Sprintf("%v", claims["localId"])

	//exp := fmt.Sprintf("%v", claims["exp"])

	res := &pb.AuthCodeToken{
		AccessToken: AccessToken,
		UserId: userId,
	}

	return res, err
}

func (j *jwtfunc) RefreshToken(ctx context.Context, req *pb.TokenRequest) (*pb.AuthCodeToken, error) {
	j.ll.Infolog.Printf("jwt =>  RefreshToken function called")
	//_, verify_err := VerifyToken(ctx, req)
	//
	//if verify_err == nil {

	refreshToken := req.RefreshToken

	if refreshToken == "" {
		j.ll.Errorlog.Println("refreshToken is not provided")
		return nil, fmt.Errorf("refreshToken is not provided")
	}

	claims := jwt.MapClaims{}

	_, err := jwt.ParseWithClaims(refreshToken, claims, func(t *jwt.Token) (interface{}, error) {
		return jwtKey, nil
	})

	//fmt.Println(token)
	//log.Println(refreshToken)

	if err != nil {
		j.ll.Errorlog.Println("Invalid RefreshToken")
		return nil, fmt.Errorf("Invalid RefreshToken")
	}

	//Check in redis if the refreshToken is revoked
	if err == nil {
		_, err = j.t.Client.RefreshToken(context.Background(), &pb.TokenRequest{RefreshToken: refreshToken})
		if err != nil {
			return nil, err
		}
	}

	localId := fmt.Sprintf("%v", claims["localId"])

	// Create user new access token from id
	userToken, err := j.CreateToken(&pb.User{UserId: localId})

	if err != nil {
		j.ll.Errorlog.Println("cannot generate new access token")
		return nil, status.Errorf(codes.Internal, "cannot generate new access token")
	}

	res := &pb.AuthCodeToken{
		AccessToken: userToken.GetAccessToken(),
		ExpiresIn: userToken.GetExpiresIn(),
		RefreshToken: userToken.GetRefreshToken(),
	}
	return res, err

	//}
	//res := &pb.Response{
	//	//IsValid: 		  verify_res.IsValid,
	//}
	//
	//return res, verify_err
}

func (j *jwtfunc) RevokeToken(ctx context.Context, req *pb.TokenRequest) (*pb.AuthCodeToken, error) {
	j.ll.Infolog.Printf("jwt =>  RevokeToken function called")
	if req.Token == "" {
		j.ll.Errorlog.Println("accessToken is not provided")
		return nil, fmt.Errorf("accessToken is not provided")
	}

	_, err := j.t.Client.RevokeToken(context.Background(), &pb.TokenRequest{Token: req.Token})
	if err != nil {
		return nil, err
	}

	j.ll.Infolog.Println("accessToken is revoked")

	res := &pb.AuthCodeToken{
		AccessToken: "revoked",
		RefreshToken: "revoked",
	}

	return res, nil
}

func (j *jwtfunc) RevokeAll(ctx context.Context, req *pb.User) (*pb.AuthCodeToken, error) {
	j.ll.Infolog.Printf("jwt =>  RevokeAll function called")
	_, err := j.t.Client.RevokeAll(context.Background(), &pb.User{UserId: req.UserId})
	if err != nil {
		return nil, err
	}

	j.ll.Infolog.Println("All user's Access and Refresh Tokens are revoked")
	res := &pb.AuthCodeToken{
		UserId: req.UserId,
		AccessToken: "revoked",
		RefreshToken: "revoked",
	}

	return res, nil
}

func (j *jwtfunc) CreateResetPasswordToken(ctx context.Context, req *pb.User) (*pb.AuthCodeToken, error) {
	j.ll.Infolog.Printf("jwt =>  CreateResetPasswordToken function called")

	expirationTime := time.Now().Add(1 * time.Hour)
	issuedTime := time.Now()

	claims := &Claims{
		LocalId: req.UserId,
		Issuer: "Ali",
		Subject: "ResetPasswordToken",
		Audience: "P99",
		ExpiresAt: expirationTime.Unix(),
		NotBefore: 0,
		IssuedAt: issuedTime.Unix(),
		Id: "",
	}

	token := jwt.NewWithClaims(jwt.SigningMethodHS256, claims)

	tokenString, err := token.SignedString(jwtKey)

	if err != nil {
		j.ll.Errorlog.Fatal("InternalServerError", err)
	}

	//store token in redis by the repo layer
	_, err = j.t.Client.CreateResetPasswordToken(context.Background(), &pb.User{AccessToken: tokenString, RefreshToken: "ResetPasswordToken"})
	if err != nil {
		j.ll.Errorlog.Println("create tokens failed")
		if j.env.Server.Mode == "debug" {
			j.ll.Debuglog.Println(err)
		}
		return nil, fmt.Errorf("create tokens failed")
	}

	res := &pb.AuthCodeToken{
		AccessToken: tokenString,
		ExpiresIn:        int32(expirationTime.Unix()),
	}

	return res, err
}
