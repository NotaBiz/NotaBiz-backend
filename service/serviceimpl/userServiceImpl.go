package serviceimpl

import (
	"NotaBiz-backend/model/entity"
	"NotaBiz-backend/model/response"
	"NotaBiz-backend/repository"
	"NotaBiz-backend/repository/repoimpl"
)

type UserService struct{}

var userRepo repository.UserRepository = repoimpl.NewUserRepository()

func (UserService) GetUserByEmail(email string) (res *response.UserResponse, err error) {
	user, err := userRepo.GetUserByEmail(email)
	if err != nil {
		return nil, err
	}
	return userToUserResponse(user), nil
}

func userToUserResponse(user *entity.MasterUser) *response.UserResponse {
	company := user.Company
	subs := company.Subscription
	role := user.Role
	return &response.UserResponse{
		ID:          user.ID,
		Name:        user.Name,
		Email:       user.Email,
		PhoneNumber: user.PhoneNumber,
		CompanyID:   user.CompanyID,
		Company: response.CompanyResponse{
			ID:             company.ID,
			Company:        company.Company,
			Description:    company.Description,
			Address:        company.Address,
			SubscriptionID: subs.ID,
			Subscription: response.SubscriptionResponse{
				ID:           subs.ID,
				Subscription: string(subs.Subscription),
			},
			CreatedAt: company.CreatedAt,
		},
		RoleID: user.RoleID,
		Role: response.RoleResponse{
			ID:   role.ID,
			Role: string(role.Role),
		},
		CreatedAt: user.CreatedAt,
		UpdatedAt: user.UpdatedAt,
	}
}
