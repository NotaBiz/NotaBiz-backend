package service

import (
	"NotaBiz-backend/model/request"
	"NotaBiz-backend/model/response"

	"github.com/markbates/goth"
)

type AuthService interface {
	GoogleCallback(googleUser goth.User) (res *response.LoginResponse, err error)
	RegisterOwner(registerRequest request.RegisterOwnerRequest) (err error)
	VerifyOTP(req request.VerifyOTPRequest) (user *response.UserResponse, err error)
	Login(req request.LoginRequest) (res *response.LoginResponse, err error)
}
