package request

import "github.com/google/uuid"

type RegisterOwnerRequest struct {
	Name               string    `json:"name" validate:"required"`
	Email              string    `json:"email,omitempty" validate:"omitempty,email"`
	PhoneNumber        string    `json:"phoneNumber,omitempty" validate:"omitempty,e164"`
	Password           string    `json:"password" validate:"required,min=8"`
	Company            string    `json:"company" validate:"required"`
	CompanyDescription string    `json:"companyDescription,omitempty"`
	CompanyAddress     string    `json:"companyAddress,omitempty"`
	SubscriptionID     uuid.UUID `json:"subscriptionID" validate:"required"`
	Method             string    `json:"method" validate:"required,oneof=email whatsapp"`
}

type VerifyOTPRequest struct {
	Email       string `json:"email,omitempty"`
	PhoneNumber string `json:"phoneNumber,omitempty"`
	OTP         string `json:"otp"`
	Method      string `json:"method"`
}

type LoginRequest struct {
	Identifier string `json:"identifier" validate:"required"`
	Password   string `json:"password" validate:"required,min=8"`
}
