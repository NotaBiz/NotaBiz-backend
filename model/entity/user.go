package entity

import (
	"time"

	"github.com/google/uuid"
	"gorm.io/gorm"
	"gorm.io/plugin/soft_delete"
)

type User struct {
	ID        uuid.UUID             `gorm:"type:uuid;primary_key;"`
	Email     string                `gorm:"unique;not null;type:varchar(255);uniqueIndex:idx_email"`
	Username  string                `gorm:"unique;not null;type:varchar(255);"`
	Password  string                `gorm:"not null;type:varchar(255)"`
	CreatedAt time.Time             `gorm:"default:now();type:timestamp"`
	UpdatedAt time.Time             `gorm:"default:now();type:timestamp"`
	DeletedAt soft_delete.DeletedAt `gorm:"softDelete:nano;default:0"`
}

func (u *User) BeforeCreate(tx *gorm.DB) error {
	if u.ID == uuid.Nil {
		u.ID = uuid.New()
	}
	return nil
}