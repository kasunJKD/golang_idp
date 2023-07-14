package db

import (
	"context"
	_ "github.com/lib/pq"

	pb "idp-repository/protos/user"
	//"database/sql"
	"log"
)

func (c DBConfig) LinkAccount(ctx context.Context, req *pb.Request) (*pb.Response, error) {

	userId := req.GetUserId()
	providerId := req.GetProviderId()
	federatedId := req.GetFederatedId()
	email := req.GetEmail()
	linkedUserId := req.GetLinkedUserId()

	//TODO: call the auth login and return the above user info

	log.Println(req)

	//Check if the idp account already linked to other p99 account
	accountExist, err := c.CheckIdpAccountLinked(&pb.Request{ProviderId: providerId, FederatedId: federatedId})

	if err != nil {
		c.log.Errorlog.Fatal(err)
	}

	if accountExist {
		c.log.Errorlog.Fatal("Account already linked to other p99 account")
	}

	sqlStatement := `INSERT INTO linkedAccounts(userId, providerId, federatedId, email, linkedUserId)
		VALUES ($1, $2, $3, $4, $5)`

	err = c.rwDb.QueryRow(sqlStatement, userId, providerId, federatedId, email, linkedUserId).Err()

	if err != nil {
		log.Fatalln(err)
	}

	res := &pb.Response{
		Users: &pb.Users{
			UserId:        userId,
			//Email:         email,
			//EmailVerified: emailVerified,
			//DisplayName:   displayName,
			//PhotoUrl:      photoUrl,
			//UpdatedAt:     timestamppb.New(updatedAt),
			//CreatedAt:     timestamppb.New(createdAt),
			//PasswordHash:  passwordHash,
			//Gender: 	   gender,
			//Address: 	   address,
			//Age: 		   age,
			//Experience:    experience,
			//PlayingTime:   playingTime,
			//PreferredPlatforms: preferredPlatforms,
			LinkedAccounts: &pb.LinkedAccounts {
				ProviderId:       providerId,
				FederatedId:      federatedId,
				Email:            email,
				LinkedUserId: 	  linkedUserId,
			},
		},
		//OauthAccessToken: verify_res.GetOauthAccessToken(),
		//FirstName:        firstName,
		//LastName:         lastName,
		//FullName:         fmt.Sprintf("%s %s", firstName, lastName),
		//ExpiresIn:        verify_res.GetExpiresIn(),
		//AccessToken:      verify_res.GetAccessToken(),
	}

	return res, err

}
