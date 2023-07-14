package service

import (
	"context"
	"fmt"
	"google.golang.org/grpc/codes"
	"google.golang.org/grpc/status"
	"idp-service/env"
	"idp-service/login-service/internal/auth"
	"idp-service/pkg/grpc_client"
	"idp-service/pkg/logger"
	proto "idp-service/protos/login"
	pb_token "idp-service/protos/token"

	otp_model "bitbucket.org/project-99-games/otp_model"
	"google.golang.org/protobuf/types/known/wrapperspb"
)

//LoginService struct
type LoginService struct {
	proto.UnimplementedLoginServiceServer
	l *grpc_client.LoginServiceClient
	t *grpc_client.TokenServiceClient
	o *grpc_client.OTPServiceClient
	ll *logger.CustomLogger
	af auth.IAuthFunc
	env *env.Env
}

//NewLoginService create service
func NewLoginService(l *grpc_client.LoginServiceClient, t *grpc_client.TokenServiceClient, o *grpc_client.OTPServiceClient, ll *logger.CustomLogger,af auth.IAuthFunc, env *env.Env) proto.LoginServiceServer {
	return &LoginService{
		l: l,
		t: t,
		ll: ll,
		af: af,
		o: o,
		env: env,
	}
}

func (service *LoginService) SignUp(ctx context.Context, req *proto.Request) (res *proto.Response, err error) {
	service.ll.Infolog.Println("Service => Signup function called")

	res, err = service.af.SignUp(ctx, req)

	if err != nil {
		if service.env.Server.Mode == "debug" {
			service.ll.Debuglog.Printf("err %v", err)
		}
		service.ll.Errorlog.Println("Signup failed")
		return nil, err
	}

	return res, err
}

func (service *LoginService) PasswordSignIn(ctx context.Context, req *proto.Request) (res *proto.Response, err error) {
	service.ll.Infolog.Println("Service => PasswordSignIn function called")
	res, err = service.af.PasswordSignIn(ctx, req)

	if err != nil {
		if service.env.Server.Mode == "debug" {
			service.ll.Debuglog.Printf("err %v", err)
		}
		service.ll.Errorlog.Println("Password Sign In failed")
		return nil, err
	}

	// check of 2fa enabled
	if res.Users.OtpEnabled == 1 {
		createOtpRequest := &otp_model.CreateOtpCodeRequest{
			UserId : res.Users.UserId,
			ServiceId: "idp",
		}

		createOtpResponse, err := service.o.Client.CreateOtpCode(ctx, createOtpRequest)

		if err != nil {
			if service.env.Server.Mode == "debug" {
				service.ll.Debuglog.Printf("err %v", err)
			}
			service.ll.Errorlog.Println("OTP service ===> CreateOtpCode failed")
			return nil, err
		}

		response := &proto.Response{
			Users: &proto.Users{
				UserId:		  	  res.Users.GetUserId(),
				Email:			  res.Users.GetEmail(),
				DisplayName:      res.Users.GetDisplayName(),
				OtpEnabled: 	  res.Users.GetOtpEnabled(),
				PreferredMethod: createOtpResponse.GetPreferredMethod(),
				OtpCode: createOtpResponse.GetOtp(),
			},
		}

		return response, err
	}
	return res, err
}

func (service *LoginService) GetClientInfo(ctx context.Context, req *proto.ClientReq) (res *proto.ClientRes, err error) {
	service.ll.Infolog.Printf("Service => GetClientInfo function called\nCalling login client")
	res, err = service.l.Client.GetClientInfo(ctx, req)
	if err != nil {
		if service.env.Server.Mode == "debug" {
			service.ll.Debuglog.Printf("err %v", err)
		}
		service.ll.Errorlog.Println("Login repository ===> GetClientInfo failed")
		return nil, err
	}
	return res, err
}

func (service *LoginService) CreateOauthClient(ctx context.Context, req *proto.ClientReq) (res *proto.ClientRes, err error) {
	service.ll.Infolog.Printf("Service => CreateOauthClient function called\nCalling login client")
	res, err = service.l.Client.CreateOauthClient(ctx, req)
	if err != nil {
		if service.env.Server.Mode == "debug" {
			service.ll.Debuglog.Printf("err %v", err)
		}
		service.ll.Errorlog.Println("Login repository ===> CreateOauthClient failed")
		return nil, err
	}
	return res, err
}

func (service *LoginService) Logout(ctx context.Context, req *proto.Request) (res *wrapperspb.BoolValue, err error) {
	service.ll.Infolog.Printf("Service => Logout function called\nCalling token client")
	_, err = service.t.Client.RevokeToken(ctx, &pb_token.TokenRequest{Token: req.AccessToken})
	if err != nil {
		if service.env.Server.Mode == "debug" {
			service.ll.Debuglog.Printf("err %v", err)
		}
		service.ll.Errorlog.Println("Logout failed")
		return wrapperspb.Bool(false), err
	}
	return wrapperspb.Bool(true), err
}

func (service *LoginService) ChangePassword(ctx context.Context, req *proto.Passwordreq) (res *proto.Status, err error) {
	service.ll.Infolog.Printf("Service => ChangePassword function called\nCalling login client")
	res, err = service.l.Client.ChangePassword(ctx, req)
	if err != nil {
		if service.env.Server.Mode == "debug" {
			service.ll.Debuglog.Printf("err %v", err)
		}
		service.ll.Errorlog.Println("Login repository ===> ChangePassword failed")
		return nil, err
	}
	return res, err
}

func (service *LoginService) GetClientsByUserId(ctx context.Context, req *proto.ClientReq) (res *proto.GetClientsByUserIdResponse, err error) {
	service.ll.Infolog.Printf("Service => GetClientsByUserId function called\nCalling login client")
	res, err = service.l.Client.GetClientsByUserId(ctx, req)
	if err != nil {
		if service.env.Server.Mode == "debug" {
			service.ll.Debuglog.Printf("err %v", err)
		}
		service.ll.Errorlog.Println("Login repository ===> GetClientsByUserId failed")
		return nil, err
	}
	return res, err
}

func (service *LoginService) ForgotPassword(ctx context.Context, req *proto.Request) (res *proto.ForgotPasswordResponse, err error) {
	service.ll.Infolog.Printf("Service => ForgotPassword function called\nCalling login client")
	res, err = service.af.ForgotPassword(ctx, req)
	if err != nil {
		if service.env.Server.Mode == "debug" {
			service.ll.Debuglog.Printf("err %v", err)
		}
		service.ll.Errorlog.Println("ForgotPassword failed")
		return nil, err
	}
	return res, err
}

func (service *LoginService) ResetPassword(ctx context.Context, req *proto.Request) (res *proto.Status, err error) {
	service.ll.Infolog.Printf("Service =>  ResetPassword function called\nCalling login client")
	token_res, err := service.t.Client.VerifyToken(ctx, &pb_token.TokenRequest{Code: "ResetPasswordToken",Token: req.AccessToken})
	if err != nil {
		service.ll.Errorlog.Println("Token service ===> ResetPassword failed")
		return nil, err
	}

	service.ll.Infolog.Printf("Service =>  ResetPassword function called\nCalling login client")
	res, err = service.l.Client.ResetPassword(ctx, &proto.Request{UserId: token_res.UserId,Password: req.Password})
	if err != nil {
		if service.env.Server.Mode == "debug" {
			service.ll.Debuglog.Printf("err %v", err)
		}
		service.ll.Errorlog.Println("Login repository ===> ResetPassword failed")
		return nil, err
	}

	_, token_err := service.t.Client.RevokeToken(ctx, &pb_token.TokenRequest{Token: req.AccessToken})
	if token_err != nil {
		service.ll.Errorlog.Println("Token service ===> Remove ResetPasswordToken failed")
	}

	return res, err
}

func (service *LoginService) EnableTwoFactorAuthentication(ctx context.Context, req *proto.EnableTwoFactorAuthenticationRequest) (res *proto.EnableTwoFactorAuthenticationResponse, err error) {
	service.ll.Infolog.Printf("Service =>  ResetPassword function called")
	//update idp 2fa to enable 
	statusUpdate := &proto.OtpRequest{
		UserId: req.UserId,
		OtpEnabled: req.OtpEnabled,
	}
	_, err = service.l.Client.UpdateotpStatus(ctx ,statusUpdate)

	if err != nil {
		if service.env.Server.Mode == "debug" {
			service.ll.Debuglog.Printf("err %v", err)
		}
		service.ll.Errorlog.Println("Login repository ===> UpdateotpStatus failed")
		return nil, err
	}

	//connect to otp enable twofactor auth and insert user preference data
	enabletwofaRequest := &otp_model.EnableTwoFARequest{
		UserId : req.UserId,
		ServiceId: "idp",
		PreferredMethod: req.OtpMethod,
		Email: req.OtpMethodData,
		PhoneNumber: req.OtpMethodData,
	}

	_, err = service.o.Client.EnableTwoFA(ctx, enabletwofaRequest)

	if err != nil {
		if service.env.Server.Mode == "debug" {
			service.ll.Debuglog.Printf("err %v", err)
		}
		service.ll.Errorlog.Println("OTP service ===> EnableTwoFA failed")
		return nil, err
	}

	response := &proto.EnableTwoFactorAuthenticationResponse{
		UserId: req.UserId,
	}

	return response, err
}

func (service *LoginService) OtpLogin(ctx context.Context, req *proto.OtpLoginRequest) (res *proto.Response, err error) {
	service.ll.Infolog.Printf("Service =>  OtpLogin function called")
	validated := false

	if req.OtpMethod == "totp" {
		totpRequest := &otp_model.VerifyTotpCodeRequest{
			UserId : req.UserId,
			ServiceId: "idp",
			Code: req.OtpCode,
		}
		resp, err := service.o.Client.VerifyTotpCode(ctx, totpRequest)
		if err != nil {
			if service.env.Server.Mode == "debug" {
				service.ll.Debuglog.Printf("err %v", err)
			}
			service.ll.Errorlog.Println("OTP service ===> VerifyTotpCode failed")
			return nil, err
		}		
		validated = resp.Validated
	} else {
		otpRequest := &otp_model.VerifyOtpCodeRequest{
			UserId : req.UserId,
			ServiceId: "idp",
			Otp: fmt.Sprint(req.OtpCode),
		}
		resp, err := service.o.Client.VerifyOtpCode(ctx, otpRequest)
		if err != nil {
			if service.env.Server.Mode == "debug" {
				service.ll.Debuglog.Printf("err %v", err)
			}
			service.ll.Errorlog.Println("OTP service ===> VerifyTotpCode failed")
			return nil, err
		}
		validated = resp.Validated
	}

	if !validated {
		service.ll.Errorlog.Println("otpLogin failed; otp code in not valid")
		return nil, status.Errorf(codes.Internal, "otp not valid")
	}

	response, err := service.l.Client.OtpLogin(ctx ,req)

	if err != nil {
		if service.env.Server.Mode == "debug" {
			service.ll.Debuglog.Printf("err %v", err)
		}
		service.ll.Errorlog.Println("OTP service ===> OtpLogin failed")
		return nil, err
	}

	return response, err
}

func (service *LoginService) MigrateAccount(ctx context.Context, req *proto.Request) (*proto.Response, error) {
	service.ll.Infolog.Printf("Service =>  MigrateAccount function called")
	data, err := service.af.MigrateAccount(ctx, req)
	if err != nil {
		if service.env.Server.Mode == "debug" {
			service.ll.Debuglog.Printf("err %v", err)
		}
		service.ll.Errorlog.Println("MigrateAccount failed")
		return nil, err
	}
	return data, err
}
