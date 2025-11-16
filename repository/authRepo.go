package repository

import "NotaBiz-backend/model/entity"

type AuthRepository interface {
	RegisterOwner(user entity.MasterUser, company  entity.MasterCompany) (err error)
	VerifyUser(email, phoneNumber, method string) (user entity.MasterUser, err error)
}