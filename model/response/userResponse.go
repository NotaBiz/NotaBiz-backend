package response

import (
	"time"

	"github.com/google/uuid"
)

type UserResponse struct {
	ID          uuid.UUID       `json:"id"`
	Name        *string         `json:"name,omitempty"`
	Email       *string          `json:"email"`
	PhoneNumber *string         `json:"phoneNumber,omitempty"`
	CompanyID   uuid.UUID       `json:"companyId"`
	Company     CompanyResponse `json:"company"`
	RoleID      uuid.UUID       `json:"roleID"`
	Role        RoleResponse    `json:"role"`
	CreatedAt   time.Time       `json:"createdAt"`
	UpdatedAt   time.Time       `json:"updatedAt"`
}
