package response

import "github.com/google/uuid"

type RoleResponse struct {
	ID   uuid.UUID `json:"id"`
	Role string    `json:"role"`
	Users []UserResponse `json:"users"`
}