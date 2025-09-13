package response

import (
	"time"

	"github.com/google/uuid"
)

type CompanyResponse struct {
	ID uuid.UUID `json:"id"`
	Company string `json:"company"`
	Description *string `json:"description"`
	Address *string `json:"address"`
	SubscriptionID uuid.UUID `json:"subscriptionID"`
	Subscription SubscriptionResponse `json:"subscription"`
	CreatedAt time.Time `json:"createdAt"`
	Users []UserResponse `json:"users"`
	// Transactions []TransactionResponse `json:"transactions"`
}