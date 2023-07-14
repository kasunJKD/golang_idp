package auth

import (
	"context"
	//"p99system/internal/auth/jwt"
	pb "idp-service/protos/login"
	tokenpb "idp-service/protos/token"
	userpb "idp-service/protos/user"

	"idp-service/login-service/internal/utils"

	"google.golang.org/grpc/codes"
	"google.golang.org/grpc/status"
)

func (a *AuthFunc) PasswordSignIn (ctx context.Context, req *pb.Request) (*pb.Response, error) {
	a.ll.Infolog.Printf("Auth =>  PasswordSignIn function called")
	userExist, err := a.u.Client.CheckAuthUserExists(ctx, &userpb.Request{Email: req.Email})

	if err != nil {
		if a.env.Server.Mode == "debug" {
			a.ll.Debuglog.Printf("err %v", err)
		}
		a.ll.Errorlog.Println("User Service ==> CheckAuthUserExists failed")
		return nil, err
	}

	if userExist.Value {
		//getUser info
		userInfo, err := a.u.Client.GetAccountInfo(ctx, &userpb.Request{Email: req.Email})
		
		if err != nil {
			if a.env.Server.Mode == "debug" {
				a.ll.Debuglog.Printf("err %v", err)
			}
			a.ll.Errorlog.Println("User Service ==> GetAccountInfo failed")
			return nil, err
		}

		//check request password with stored password
		err = utils.CheckPassword(req.Password, userInfo.Users.PasswordHash)
		if err != nil {
			if a.env.Server.Mode == "debug" {
				a.ll.Debuglog.Printf("err %v", err)
			}
			a.ll.Errorlog.Println("Util func ==> CheckPassword failed")
			return nil, status.Errorf(codes.Internal, "Password incorrect")
		}

		var accessToken, refreshToken string

		//check if 2fa disabled
		if userInfo.Users.OtpEnabled == 0 {
			// Create user access token from id
			userToken, err := a.t.Client.CreateToken(ctx, &tokenpb.User{UserId: userInfo.Users.UserId})

			if err != nil {
				if a.env.Server.Mode == "debug" {
					a.ll.Debuglog.Printf("err %v", err)
				}
				a.ll.Errorlog.Println("Token service ==> CreateToken failed")
				return nil, status.Errorf(codes.Internal, "cannot generate access token")
			}

			accessToken = userToken.GetAccessToken()
			refreshToken = userToken.GetRefreshToken()
		}

		//return response
		res := &pb.Response{
			Users: &pb.Users{
				UserId:		  	  userInfo.Users.GetUserId(),
				Email:			  userInfo.Users.GetEmail(),
				EmailVerified:    userInfo.Users.GetEmailVerified(),
				DisplayName:      userInfo.Users.GetDisplayName(),
				PhotoUrl:         userInfo.Users.GetPhotoUrl(),
				UpdatedAt:        userInfo.Users.GetUpdatedAt(),
				CreatedAt:        userInfo.Users.GetCreatedAt(),
				Gender: 		  userInfo.Users.GetGender(),
				Address: 		  userInfo.Users.GetAddress(),
				Age: 			  userInfo.Users.GetAge(),
				Experience: 	  userInfo.Users.GetExperience(),
				PlayingTime: 	  userInfo.Users.GetPlayingTime(),
				PreferredPlatforms: userInfo.Users.GetPreferredPlatforms(),
				OtpEnabled: 	  userInfo.Users.GetOtpEnabled(),
				//LinkedAccounts: &pb.LinkedAccounts {
				//	ProviderId:       "providerId",
				//	FederatedId:      "federatedId",
				//	Email:            "email",
				//},
			},
			FirstName:        userInfo.GetFirstName(),
			LastName:         userInfo.GetLastName(),
			FullName:         userInfo.GetFullName(),
			OauthAccessToken: accessToken,
			//ExpiresIn:        userToken.GetExpiresIn(),
			//IsNewUser:		  userInfo.GetIsNewUser(),
			RefreshToken: 	  refreshToken,
		}

		return res, err

	}

	return nil, status.Errorf(codes.Internal, "User not valid")

}