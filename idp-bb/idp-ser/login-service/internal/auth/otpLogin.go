package auth

import (
	"context"
	"google.golang.org/grpc/codes"
	"google.golang.org/grpc/status"
	pb "idp-service/protos/login"
	tokenpb "idp-service/protos/token"
	userpb "idp-service/protos/user"
)

func (a *AuthFunc) OtpLogin (ctx context.Context, req *pb.OtpLoginRequest) (*pb.Response, error) {
	a.ll.Infolog.Printf("Auth =>  OtpLogin function called")
	userInfo, err := a.u.Client.GetUserInfoById(ctx, &userpb.Request{UserId: req.UserId})

	if err != nil {
		if a.env.Server.Mode == "debug" {
			a.ll.Debuglog.Printf("err %v", err)
		}
		a.ll.Errorlog.Println("User Service ==> GetUserInfoById failed")
		return nil, err
	}

	userToken, err := a.t.Client.CreateToken(ctx, &tokenpb.User{UserId: req.UserId})

	if err != nil {
		if a.env.Server.Mode == "debug" {
			a.ll.Debuglog.Printf("err %v", err)
		}
		a.ll.Errorlog.Println("Token Service ==> CreateToken failed")
		return nil, status.Errorf(codes.Internal, "cannot generate access token")
	}

	//return response
	res := &pb.Response{
		Users: &pb.Users{
			UserId:		  	  req.UserId,
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
			PreferredMethod:  req.OtpMethod,
		},
		FirstName:        userInfo.GetFirstName(),
		LastName:         userInfo.GetLastName(),
		FullName:         userInfo.GetFullName(),
		OauthAccessToken: userToken.GetAccessToken(),
		RefreshToken: 	  userToken.GetRefreshToken(),
	}

	return res, err
}