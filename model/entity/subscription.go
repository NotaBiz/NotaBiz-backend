package entity

import (
	"time"

	"github.com/google/uuid"
	"gorm.io/gorm"
	"gorm.io/plugin/soft_delete"
)

type SubscriptionName string

const (
	SubscriptionFree    SubscriptionName = "free"
	SubscriptionBasic   SubscriptionName = "basic"
	SubscriptionAdvance SubscriptionName = "advance"
)

type MasterSubscription struct {
	ID           uuid.UUID             `gorm:"type:uuid;primaryKey"`
	Subscription SubscriptionName      `gorm:"type:varchar(20);unique;not null"`
	CreatedAt    time.Time             `gorm:"default:now();type:timestamp"`
	UpdatedAt    time.Time             `gorm:"default:now();type:timestamp"`
	DeletedAt    soft_delete.DeletedAt `gorm:"softDelete:datetime"`
	Companies    []MasterCompany       `gorm:"foreignKey:SubscriptionID"`
}

func (s *MasterSubscription) BeforeCreate(tx *gorm.DB) error {
	if s.ID == uuid.Nil {
		s.ID = uuid.New()
	}
	return nil
}
