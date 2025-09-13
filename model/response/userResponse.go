package response

import (
	"time"

	"github.com/google/uuid"
)

type UserResponse struct {
	ID          uuid.UUID       `json:"id"`
	Name        string          `json:"name"`
	Email       string          `json:"email"`
	Username    string          `json:"username"`
	PhoneNumber string          `json:"phoneNumber"`
	CompanyID   uuid.UUID       `json:"companyId"`
	Company     CompanyResponse `json:"company"`
	RoleID      uuid.UUID       `json:"roleID"`
	Role        RoleResponse    `json:"role"`
	CreatedAt   time.Time       `json:"createdAt"`
	UpdatedAt   time.Time       `json:"updatedAt"`
}
