package db

import (
	"context"
	_ "github.com/lib/pq"

	pb "idp-repository/protos/user"
	//"database/sql"
	"log"
)

func (c DBConfig) UnlinkAccount(ctx context.Context, req *pb.Request) (*pb.Response, error) {

	providerId := req.GetProviderId()
	federatedId := req.GetFederatedId()

	log.Println(req)

	sqlStatement := `DELETE FROM linkedAccounts l
		WHERE l.providerId = $1 AND l.federatedId = $2`

	err := c.rwDb.QueryRow(sqlStatement, providerId, federatedId).Err()

	if err != nil {
		log.Fatalln(err)
	}
	c.log.Infolog.Println("Account is unlinked successfully")

	res := &pb.Response{
		Users: &pb.Users{
			//UserId:        userId,
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
			//LinkedAccounts: &pb.LinkedAccounts {
			//	ProviderId:       providerId,
			//	FederatedId:      federatedId,
			//	Email:            email,
			//},
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
