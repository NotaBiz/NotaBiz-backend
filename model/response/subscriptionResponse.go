package response

import (
	"github.com/google/uuid"
)

type SubscriptionResponse struct {
	ID           uuid.UUID         `json:"id"`
	Subscription string            `json:"subscription"`
	Companies    []CompanyResponse `json:"companies"`
}
