package entity

import (
	"time"

	"github.com/google/uuid"
	"gorm.io/gorm"
	"gorm.io/plugin/soft_delete"
)

type Feedback struct {
	ID        uuid.UUID             `gorm:"type:uuid;primaryKey"`
	Feedback  string                `gorm:"type:varchar(255);not null"`
	UserID    uuid.UUID             `gorm:"type:uuid;not null"`
	User      MasterUser            `gorm:"foreignKey:UserID"`
	CreatedAt time.Time             `gorm:"default:now();type:timestamp"`
	CreatedBy *uuid.UUID            `gorm:"type:uuid"`
	UpdatedAt time.Time             `gorm:"default:now();type:timestamp"`
	UpdatedBy *uuid.UUID            `gorm:"type:uuid"`
	DeletedAt soft_delete.DeletedAt `gorm:"softDelete:nano;default:0"`
}

func (f *Feedback) BeforeCreate(tx *gorm.DB) error {
	if f.ID == uuid.Nil {
		f.ID = uuid.New()
	}
	return nil
}
