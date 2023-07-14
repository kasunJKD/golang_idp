package auth

import (
	"context"
	pb "idp-service/protos/login"
	tokenpb "idp-service/protos/token"
	userpb "idp-service/protos/user"

	"google.golang.org/grpc/codes"
	"google.golang.org/grpc/status"
)

func (a *AuthFunc) MigrateAccount (ctx context.Context, req *pb.Request) (*pb.Response, error) {
	a.ll.Infolog.Printf("Auth =>  MigrateAccount function called")
	//check user exists -- returns true/false
	p99UserExist, err := a.u.Client.CheckAuthUserExists(ctx, &userpb.Request{Email: req.Email})

	if err != nil {
		if a.env.Server.Mode == "debug" {
			a.ll.Debuglog.Printf("err %v", err)
		}
		a.ll.Errorlog.Println("User Service ==> CheckAuthUserExists failed")
		return nil, err
	}

	if p99UserExist.Value {
		return nil, status.Errorf(codes.Internal, "User already exists")
	}

	//Check if the idp account already linked to other p99 account, and is so unlink the account
	idpAccountExist, err := a.u.Client.CheckIdpAccountLinked(ctx, &userpb.Request{ProviderId: req.ProviderId, FederatedId: req.FederatedId})

	if err != nil {
		if a.env.Server.Mode == "debug" {
			a.ll.Debuglog.Printf("err %v", err)
		}
		a.ll.Errorlog.Println("User Service ==> CheckIdpAccountLinked failed")
		return nil, err
	}

	if idpAccountExist.Value {
		a.ll.Infolog.Println("Account already linked to other p99 account, account will be unlinked..")

		_, err = a.u.Client.UnlinkAccount(ctx, &userpb.Request{ProviderId: req.ProviderId, FederatedId: req.FederatedId})

		if err != nil {
			if a.env.Server.Mode == "debug" {
				a.ll.Debuglog.Printf("err %v", err)
			}
			a.ll.Errorlog.Println("User Service ==> UnlinkAccount failed")
			return nil, err
		}
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
		LinkedUserId: req.LinkedUserId,
	}
	migrateAccountResponse, err := a.u.Client.MigrateAccount(ctx, ureq)

	if err != nil {
		if a.env.Server.Mode == "debug" {
			a.ll.Debuglog.Printf("err %v", err)
		}
		a.ll.Errorlog.Println("User Service ==> MigrateAccount failed")
		return nil, err
	}

	userToken, err := a.t.Client.CreateToken(ctx, &tokenpb.User{UserId: migrateAccountResponse.Users.UserId})

	if err != nil {
		if a.env.Server.Mode == "debug" {
			a.ll.Debuglog.Printf("err %v", err)
		}
		a.ll.Errorlog.Println("Token Service ==> CreateToken failed")
		return nil, status.Errorf(codes.Internal, "cannot generate access token")
	}

	res := &pb.Response{
		Users: &pb.Users{
			UserId:		  		migrateAccountResponse.Users.GetUserId(),
			Email:			  	migrateAccountResponse.Users.GetEmail(),
			EmailVerified:    	migrateAccountResponse.Users.GetEmailVerified(),
			DisplayName:      	migrateAccountResponse.Users.GetDisplayName(),
			PhotoUrl:         	migrateAccountResponse.Users.GetPhotoUrl(),
			UpdatedAt:        	migrateAccountResponse.Users.GetUpdatedAt(),
			CreatedAt:        	migrateAccountResponse.Users.GetCreatedAt(),
			Gender: 			migrateAccountResponse.Users.GetGender(),
			Address: 			migrateAccountResponse.Users.GetAddress(),
			Age: 			 	migrateAccountResponse.Users.GetAge(),
			Experience: 	 	migrateAccountResponse.Users.GetExperience(),
			PlayingTime: 	 	migrateAccountResponse.Users.GetPlayingTime(),
			PreferredPlatforms: migrateAccountResponse.Users.GetPreferredPlatforms(),
		},
		FirstName:        migrateAccountResponse.GetFirstName(),
		LastName:         migrateAccountResponse.GetLastName(),
		FullName:         migrateAccountResponse.GetFullName(),
		OauthAccessToken: userToken.GetAccessToken(),
		//ExpiresIn:        userToken.GetExpiresIn(),
		IsNewUser:		  migrateAccountResponse.GetIsNewUser(),
		RefreshToken: 	  userToken.GetRefreshToken(),
	}

	return res, err

}