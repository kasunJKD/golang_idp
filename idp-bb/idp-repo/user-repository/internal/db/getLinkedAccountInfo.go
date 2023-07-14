package db

import (
	"context"
	_ "github.com/lib/pq"
	pb "idp-repository/protos/user"
	//"database/sql"
	"log"
)

func (c DBConfig) GetLinkedAccountInfo(ctx context.Context, req *pb.Request) (*pb.Response, error) {

	userId := req.GetUserId()
	providerId := req.GetProviderId()

	log.Println(req)

	sqlStatement := `select l.federatedId, l.email, l.linkedUserId
				from linkedAccounts l
				where l.userId = $1 AND l.providerId = $2`

	var (
		federatedId string
		email   	string
		linkedUserId string
	)

	err := c.roDb.QueryRow(sqlStatement, userId, providerId).Scan(&federatedId, &email, &linkedUserId)

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
