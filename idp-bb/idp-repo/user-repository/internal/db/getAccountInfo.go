package db

import (
	"context"
	"fmt"

	pb "idp-repository/protos/user"
	//"database/sql"
	"log"
	"time"

	_ "github.com/lib/pq"
	timestamppb "google.golang.org/protobuf/types/known/timestamppb"
)

func (c DBConfig) GetAccountInfo(ctx context.Context, req *pb.Request) (*pb.Response, error) {

	email := req.GetEmail()
	log.Println(req)
	//sqlStatement := ("select u.localId, u.createdAt, u.updatedAt, u.emailVerified, i.email, i.displayName, i.photoUrl, p.federatedId, p.providerId from users u JOIN userInfo i ON u.localId = i.localId JOIN providerUserInfo p on i.localId = p.localId where p.providerId= $1 or p.federatedId = $2")
	sqlStatement := `select u.userId, u.createdAt, u.updatedAt, u.emailVerified, u.passwordHash, u.OTP,
				 i.displayName, i.firstName, i.lastName, i.photoUrl, i.gender, i.address, i.age, i.experience, i.playingTime, i.preferredPlatforms
				from users u JOIN userInfo i ON u.userId = i.userId
				where u.email= $1`

	var (
		userId        string
		createdAt     time.Time
		updatedAt     time.Time
		emailVerified bool
		passwordHash  string
		otpEnabled	  int32
		displayName   string
		firstName	  string
		lastName	  string
		photoUrl      string
		gender        string
		address		  string
		age      	  int32
		experience    string
		playingTime   int32
		preferredPlatforms string
	)

	err := c.roDb.QueryRow(sqlStatement, email).Scan(&userId, &createdAt, &updatedAt, &emailVerified, &passwordHash, &otpEnabled, &displayName, &firstName, &lastName, &photoUrl, &gender, &address, &age, &experience, &playingTime, &preferredPlatforms)

	if err != nil {
		log.Fatalln(err)
	}

	res := &pb.Response{
		Users: &pb.Users{
			UserId:        userId,
			Email:         email,
			EmailVerified: emailVerified,
			DisplayName:   displayName,
			PhotoUrl:      photoUrl,
			UpdatedAt:     timestamppb.New(updatedAt),
			CreatedAt:     timestamppb.New(createdAt),
			PasswordHash:  passwordHash,
			Gender: 	   gender,
			Address: 	   address,
			Age: 		   age,
			Experience:    experience,
			PlayingTime:   playingTime,
			PreferredPlatforms: preferredPlatforms,
			OtpEnabled: otpEnabled,
			//LinkedAccounts: &pb.LinkedAccounts {
			//	ProviderId:       providerId,
			//	FederatedId:      federatedId,
			//	Email:            idpEmail,
			//},
		},
		//OauthAccessToken: verify_res.GetOauthAccessToken(),
		FirstName:        firstName,
		LastName:         lastName,
		FullName:         fmt.Sprintf("%s %s", firstName, lastName),
		//ExpiresIn:        verify_res.GetExpiresIn(),
		//AccessToken:      verify_res.GetAccessToken(),
	}

	return res, err

}

func (c DBConfig) GetAccountInfoById(ctx context.Context, req *pb.Request) (*pb.Response, error) {

	userID := req.GetUserId()
	log.Println(req)
	//sqlStatement := ("select u.localId, u.createdAt, u.updatedAt, u.emailVerified, i.email, i.displayName, i.photoUrl, p.federatedId, p.providerId from users u JOIN userInfo i ON u.localId = i.localId JOIN providerUserInfo p on i.localId = p.localId where p.providerId= $1 or p.federatedId = $2")
	sqlStatement := `select u.userId, u.createdAt, u.updatedAt, u.emailVerified, u.OTP,
				 i.displayName, i.firstName, i.lastName, i.photoUrl, i.gender, i.address, i.age, i.experience, i.playingTime, i.preferredPlatforms, u.email
				from users u JOIN userInfo i ON u.userId = i.userId
				where u.userId= $1`

	var (
		userId        string
		createdAt     time.Time
		updatedAt     time.Time
		emailVerified bool
		otpEnabled	  int32
		displayName   string
		firstName	  string
		lastName	  string
		photoUrl      string
		gender        string
		address		  string
		age      	  int32
		experience    string
		playingTime   int32
		preferredPlatforms string
		email		string
	)

	err := c.roDb.QueryRow(sqlStatement, userID).Scan(&userId, &createdAt, &updatedAt, &emailVerified, &otpEnabled, &displayName, &firstName, &lastName, &photoUrl, &gender, &address, &age, &experience, &playingTime, &preferredPlatforms, &email)

	if err != nil {
		log.Fatalln(err)
	}

	res := &pb.Response{
		Users: &pb.Users{
			UserId:        userID,
			Email:         email,
			EmailVerified: emailVerified,
			DisplayName:   displayName,
			PhotoUrl:      photoUrl,
			UpdatedAt:     timestamppb.New(updatedAt),
			CreatedAt:     timestamppb.New(createdAt),
			Gender: 	   gender,
			Address: 	   address,
			Age: 		   age,
			Experience:    experience,
			PlayingTime:   playingTime,
			PreferredPlatforms: preferredPlatforms,
			OtpEnabled: otpEnabled,
			//LinkedAccounts: &pb.LinkedAccounts {
			//	ProviderId:       providerId,
			//	FederatedId:      federatedId,
			//	Email:            idpEmail,
			//},
		},
		//OauthAccessToken: verify_res.GetOauthAccessToken(),
		FirstName:        firstName,
		LastName:         lastName,
		FullName:         fmt.Sprintf("%s %s", firstName, lastName),
		//ExpiresIn:        verify_res.GetExpiresIn(),
		//AccessToken:      verify_res.GetAccessToken(),
	}

	return res, err

}
