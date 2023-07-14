package auth

import (
	"context"
	"google.golang.org/grpc/codes"
	"google.golang.org/grpc/status"
	pb "idp-service/protos/login"
	tokenpb "idp-service/protos/token"
	userpb "idp-service/protos/user"
)

func (a *AuthFunc) ForgotPassword(ctx context.Context, req *pb.Request) (*pb.ForgotPasswordResponse, error) {
	a.ll.Infolog.Printf("Auth =>  ForgotPassword function called")
	userExist, err := a.u.Client.CheckAuthUserExists(ctx, &userpb.Request{Email: req.Email})

	if err != nil {
		if a.env.Server.Mode == "debug" {
			a.ll.Debuglog.Printf("err %v", err)
		}
		a.ll.Errorlog.Println("User Service ==> CheckAuthUserExists failed")
		return nil, err
	}

	if !userExist.Value {
		return nil, status.Errorf(codes.Internal, "Account not exists; incorrect email")
	}

	userInfo, err := a.u.Client.GetAccountInfo(ctx, &userpb.Request{Email: req.Email})
	
	if err != nil {
		if a.env.Server.Mode == "debug" {
			a.ll.Debuglog.Printf("err %v", err)
		}
		a.ll.Errorlog.Println("User Service ==> GetAccountInfo failed")
		return nil, err
	}

	userToken, err := a.t.Client.CreateResetPasswordToken(ctx, &tokenpb.User{UserId: userInfo.Users.UserId})

	if err != nil {
		if a.env.Server.Mode == "debug" {
			a.ll.Debuglog.Printf("err %v", err)
		}
		a.ll.Errorlog.Println("Token Service ==> CreateResetPasswordToken failed")
		return nil, status.Errorf(codes.Internal, "cannot generate ResetPasswordToken")
	}

	resetUrl := "http://localhost:44201/resetPassword?token=" + userToken.GetAccessToken()

	res := &pb.ForgotPasswordResponse{
		Code: 0,
		Message: "resetUrl generated successfully",
		ResetUrl: resetUrl,
	}

	return res, nil

}
