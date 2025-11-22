package repoimpl

import (
	"NotaBiz-backend/config"
	"NotaBiz-backend/model/entity"
	"fmt"

	"github.com/google/uuid"
)

type AuthRepository struct{}

// VerifyUser implements repository.AuthRepository.
func (a *AuthRepository) VerifyUser(email string, phoneNumber string, method string) (user entity.MasterUser, err error) {
	if method == "email" {
		if err := config.DB.Preload("Role").Preload("Company.Subscription").Where("email = ?", email).First(&user).Error; err != nil {
			return user, err
		}
	} else {
		if err := config.DB.Preload("Role").Preload("Company.Subscription").Where("phone_number = ?", phoneNumber).First(&user).Error; err != nil {
			return user, err
		}
	}
	user.IsVerified = true
	if err := config.DB.Save(&user).Error; err != nil {
		return user, err
	}
	return user, nil
}

// Register implements repository.AuthRepository.
func (a *AuthRepository) RegisterOwner(user entity.MasterUser, company entity.MasterCompany) (err error) {
	var existUser entity.MasterUser
	if user.Email != nil {
		if err := config.DB.Where("email = ?", *user.Email).Find(&existUser).Error; err != nil {
			return err
		}
	} else {
		if err := config.DB.Where("phone_number = ?", *user.PhoneNumber).Find(&existUser).Error; err != nil {
			return err
		}
	}
	if existUser.ID != uuid.Nil && existUser.IsVerified {
		return fmt.Errorf("user already exists")
	}

	var subs entity.MasterSubscription
	if err := config.DB.Select("id").Where("subscription = ?", entity.SubscriptionFree).First(&subs).Error; err != nil {
		return err
	}
	var role entity.MasterRole
	if err := config.DB.Select("id").Where("role = ?", entity.RoleOwner).First(&role).Error; err != nil {
		return err
	}

	userId := uuid.New()
	tx := config.DB.Begin()
	defer func() {
		if r := recover(); r != nil {
			tx.Rollback()
		}
	}()

	if existUser.ID == uuid.Nil {
		company.CreatedBy = &userId
		company.SubscriptionID = subs.ID
		if err := tx.Create(&company).Error; err != nil {
			tx.Rollback()
			return err
		}

		user.CompanyID = company.ID
		user.RoleID = role.ID
		user.ID = userId
		if err := tx.Create(&user).Error; err != nil {
			tx.Rollback()
			return err
		}
	} else {
		if err := tx.Model(entity.MasterCompany{}).Where("id = ?", existUser.CompanyID).
		Updates(map[string]interface{}{
			"company":     company.Company,
			"description": company.Description,
			"address":     company.Address,
			"updated_by":  existUser.ID,
		}).Error; err != nil {
			return err
		}

		if err := tx.Model(&existUser).Updates(map[string]interface{}{
			"name":     user.Name,
			"password": user.Password,
		}).Error; err != nil {
			return err
		}
	}

	tx.Commit()
	return err
}

func NewAuthRepository() *AuthRepository {
	return &AuthRepository{}
}
