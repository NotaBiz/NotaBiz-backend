package repoimpl

import (
	"NotaBiz-backend/config"
	"NotaBiz-backend/model/entity"
)

type UserRepository struct{}

func NewUserRepository() *UserRepository {
	return &UserRepository{}
}

func (UserRepository) GetUserByEmail(email string) (user *entity.MasterUser, err error) {
	err = config.DB.Debug().Preload("Role").Preload("Company.Subscription").Where("email = ? AND is_verified = ?", email, true).First(&user).Error
	if err != nil {
		return nil, err
	}
	return user, nil
}
