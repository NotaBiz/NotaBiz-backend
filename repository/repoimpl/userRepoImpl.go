package repoimpl

import (
	"NotaBiz-backend/config"
	"NotaBiz-backend/model/entity"
)

type UserRepository struct{}

// GetUser implements repository.UserRepository.
func (u *UserRepository) GetUser(identifier string) (user *entity.MasterUser, err error) {
	err = config.DB.Debug().Preload("Role").Preload("Company.Subscription").Where("(email = ? OR phone_number = ?) AND is_verified = ?", identifier, identifier, true).First(&user).Error
	if err != nil {
		return nil, err
	}
	return user, nil
}

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
