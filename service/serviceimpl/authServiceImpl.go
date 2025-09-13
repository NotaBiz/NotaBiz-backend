package serviceimpl

import (
	"NotaBiz-backend/middleware"
	"NotaBiz-backend/model/entity"
	"NotaBiz-backend/model/response"

	"github.com/markbates/goth"
)

type AuthService struct{}

func NewAuthService() *AuthService {
	return &AuthService{}
}

// var authRepo repository.AuthRepository = repoimpl.

// GoogleCallback implements service.AuthService.
func (AuthService) GoogleCallback(googleUser goth.User) (res *response.LoginResponse, err error) {
	user, err := UserService{}.GetUserByEmail(googleUser.Email)
	if err != nil {
		return nil, err
	}
	token, err := middleware.GenerateTokenJwt(user.ID, user.Email, user.PhoneNumber, entity.RoleName(user.Role.Role), entity.SubscriptionName(user.Company.Subscription.Subscription))
	if err != nil {
		return nil, err
	}

	return toLoginResponse(user, token), nil
}

func toLoginResponse(user *response.UserResponse, token string) *response.LoginResponse {
	return &response.LoginResponse{
		Token: token,
		User:  *user,
	}
}
