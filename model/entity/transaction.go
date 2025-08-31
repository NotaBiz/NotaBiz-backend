package entity

import (
	"time"

	"github.com/google/uuid"
	"gorm.io/gorm"
	"gorm.io/plugin/soft_delete"
)

type TransactionType string

const (
	TransactionIn  TransactionType = "in"
	TransactionOut TransactionType = "out"
)

type Transaction struct {
	ID          uuid.UUID             `gorm:"type:uuid;primaryKey"`
	Type        TransactionType       `gorm:"type:varchar(10);not null"`
	Amount      float64               `gorm:"type:float;not null"`
	Date        time.Time             `gorm:"type:timestamp;not null;default:now()"`
	Description *string               `gorm:"type:text"`
	CompanyID   uuid.UUID             `gorm:"type:uuid;not null"`
	Company     MasterCompany         `gorm:"foreignKey:CompanyID"`
	CreatedAt   time.Time             `gorm:"default:now();type:timestamp"`
	CreatedBy   *uuid.UUID            `gorm:"type:uuid"`
	UpdatedAt   time.Time             `gorm:"default:now();type:timestamp"`
	UpdatedBy   *uuid.UUID            `gorm:"type:uuid"`
	DeletedAt   soft_delete.DeletedAt `gorm:"softDelete:nano;default:0"`
	Attachments []Attachment          `gorm:"foreignKey:TransactionID"`
}

func (t *Transaction) BeforeCreate(tx *gorm.DB) error {
	if t.ID == uuid.Nil {
		t.ID = uuid.New()
	}
	return nil
}
