package repository

import (
	"NotaBiz-backend/model/entity"
)

type UserRepository interface {
	// GetUsers(page, pageSize int, search string, companyId uuid.UUID) ()
	GetUserByEmail(email string) (user *entity.MasterUser, err error)
}