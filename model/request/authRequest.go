package request

import "github.com/google/uuid"

type RegisterRequest struct {
	Email       string    `json:"email,omitempty" validate:"omitempty,email"`
	PhoneNumber string    `json:"phoneNumber,omitempty" validate:"omitempty,e164"`
	Password    string    `json:"password" validate:"required,min=8"`
	CompanyID   uuid.UUID `json:"companyID" validate:"required"`
}
