package service

import (
	"NotaBiz-backend/model/response"
)

type UserService interface {
	GetUsers(query string) (res []response.UserResponse, paging response.Paging, err error)
	GetUserByEmail(email string) (res response.UserResponse, err error)
}
