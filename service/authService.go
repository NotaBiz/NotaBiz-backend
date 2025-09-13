package service

import (
	// "NotaBiz-backend/model/request"
	"NotaBiz-backend/model/response"

	"github.com/markbates/goth"
)

type AuthService interface {
	GoogleCallback(googleUser goth.User) (res *response.LoginResponse, err error)
	// Register(registerRequest request.RegisterRequest) (res *response)
}
