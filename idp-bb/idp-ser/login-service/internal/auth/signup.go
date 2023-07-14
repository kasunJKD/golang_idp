package auth

import (
	"context"
	//"p99system/internal/auth/jwt"
	pb "idp-service/protos/login"
	tokenpb "idp-service/protos/token"
	userpb "idp-service/protos/user"

	"google.golang.org/grpc/codes"
	"google.golang.org/grpc/status"
)

func (a *AuthFunc) SignUp (ctx context.Context, req *pb.Request) (*pb.Response, error) {
	//check user exists -- returns true/false
	a.ll.Infolog.Printf("Auth =>  SignUp function called")
	userExist, err := a.u.Client.CheckAuthUserExists(ctx, &userpb.Request{Email: req.Email})

	if err != nil {
		if a.env.Server.Mode == "debug" {
			a.ll.Debuglog.Printf("err %v", err)
		}
		a.ll.Errorlog.Println("User Service ==> CheckAuthUserExists failed")
		return nil, err
	}

	if userExist.Value {
		return nil, status.Errorf(codes.Internal, "User already exists")
	}

	//create new user
	ureq := &userpb.Request{
		UserId : req.UserId,
		ProviderId : req.ProviderId,
		AccessToken : req.AccessToken,
		EmailVerified : req.EmailVerified,
		Email : req.Email,
		OauthAccessToken : req.OauthAccessToken,
		FirstName : req.FirstName,
		LastName: req.LastName,
		FullName: req.FullName,
		DisplayName : req.DisplayName,
		PhotoUrl : req.PhotoUrl,
		ExpiresIn : req.ExpiresIn,
		FederatedId : req.FederatedId,
		LocalId : req.LocalId,
		RefreshToken : req.RefreshToken,
		Password: req.Password,
		Gender: req.Gender,
		Address: req.Address,
		Age : req.Age,
		Experience : req.Experience,
		PlayingTime : req.PlayingTime,
		PreferredPlatforms : req.PreferredPlatforms,
	}
	userCreateResponse, err := a.u.Client.CreateNewUser(ctx, ureq)

	if err != nil {
		if a.env.Server.Mode == "debug" {
			a.ll.Debuglog.Printf("err %v", err)
		}
		a.ll.Errorlog.Println("User Service ==> CreateNewUser failed")
		return nil, err
	}

	userToken, err := a.t.Client.CreateToken(ctx, &tokenpb.User{UserId: userCreateResponse.Users.UserId})

	if err != nil {
		if a.env.Server.Mode == "debug" {
			a.ll.Debuglog.Printf("err %v", err)
		}
		a.ll.Errorlog.Println("Token Service ==> CreateToken failed")
		return nil, status.Errorf(codes.Internal, "cannot generate access token")
	}

	res := &pb.Response{
		Users: &pb.Users{
			UserId:		  		userCreateResponse.Users.GetUserId(),
			Email:			  	userCreateResponse.Users.GetEmail(),
			EmailVerified:    	userCreateResponse.Users.GetEmailVerified(),
			DisplayName:      	userCreateResponse.Users.GetDisplayName(),
			PhotoUrl:         	userCreateResponse.Users.GetPhotoUrl(),
			UpdatedAt:        	userCreateResponse.Users.GetUpdatedAt(),
			CreatedAt:        	userCreateResponse.Users.GetCreatedAt(),
			Gender: 			userCreateResponse.Users.GetGender(),
			Address: 			userCreateResponse.Users.GetAddress(),
			Age: 			 	userCreateResponse.Users.GetAge(),
			Experience: 	 	userCreateResponse.Users.GetExperience(),
			PlayingTime: 	 	userCreateResponse.Users.GetPlayingTime(),
			PreferredPlatforms: userCreateResponse.Users.GetPreferredPlatforms(),
		},
		FirstName:        userCreateResponse.GetFirstName(),
		LastName:         userCreateResponse.GetLastName(),
		FullName:         userCreateResponse.GetFullName(),
		OauthAccessToken: userToken.GetAccessToken(),
		//ExpiresIn:        userToken.GetExpiresIn(),
		IsNewUser:		  userCreateResponse.GetIsNewUser(),
		RefreshToken: 	  userToken.GetRefreshToken(),
	}

	return res, err

}