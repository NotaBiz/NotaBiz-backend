package entity

import (
	"time"

	"github.com/google/uuid"
	"gorm.io/gorm"
	"gorm.io/plugin/soft_delete"
)

type MasterCompany struct {
	ID             uuid.UUID             `gorm:"type:uuid;primaryKey"`
	Company        string                `gorm:"type:varchar(255);not null"`
	Description    *string               `gorm:"type:text"`
	Address        *string               `gorm:"type:text"`
	SubscriptionID uuid.UUID             `gorm:"type:uuid;not null"`
	Subscription   MasterSubscription    `gorm:"foreignKey:SubscriptionID"`
	CreatedAt      time.Time             `gorm:"default:now();type:timestamp"`
	CreatedBy      *uuid.UUID            `gorm:"type:uuid"`
	UpdatedAt      time.Time             `gorm:"default:now();type:timestamp"`
	UpdatedBy      *uuid.UUID            `gorm:"type:uuid"`
	DeletedAt      soft_delete.DeletedAt `gorm:"softDelete:nano;default:0"`
	Users          []MasterUser          `gorm:"foreignKey:CompanyID"`
	Transactions   []Transaction         `gorm:"foreignKey:CompanyID"`
}

func (c *MasterCompany) BeforeCreate(tx *gorm.DB) error {
	if c.ID == uuid.Nil {
		c.ID = uuid.New()
	}
	return nil
}
