package entity

import (
	"time"

	"github.com/google/uuid"
	"gorm.io/gorm"
	"gorm.io/plugin/soft_delete"
)

type Attachment struct {
	ID            uuid.UUID             `gorm:"type:uuid;primaryKey"`
	TransactionID uuid.UUID             `gorm:"type:uuid;not null"`
	Transaction   Transaction           `gorm:"foreignKey:TransactionID"`
	FileName      *string               `gorm:"type:varchar(255);not null"`
	FilePath      string                `gorm:"type:varchar(255);not null"`
	FileType      *string               `gorm:"type:varchar(100);not null"`
	CreatedAt     time.Time             `gorm:"default:now();type:timestamp"`
	CreatedBy     *uuid.UUID            `gorm:"type:uuid"`
	UpdatedAt     time.Time             `gorm:"default:now();type:timestamp"`
	UpdatedBy     *uuid.UUID            `gorm:"type:uuid"`
	DeletedAt     soft_delete.DeletedAt `gorm:"softDelete:nano;default:0"`
}

func (a *Attachment) BeforeCreate(tx *gorm.DB) error {
	if a.ID == uuid.Nil {
		a.ID = uuid.New()
	}
	return nil
}
