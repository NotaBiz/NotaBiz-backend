package entity

import (
	"fmt"
	"time"

	"github.com/google/uuid"
	"gorm.io/gorm"
	"gorm.io/plugin/soft_delete"
)

type MasterUser struct {
	ID          uuid.UUID             `gorm:"type:uuid;primary_key;"`
	Name        *string               `gorm:"type:varchar(255);"`
	Email       *string               `gorm:"type:varchar(255);uniqueIndex:idx_deleted_at_email"`
	PhoneNumber *string               `gorm:"type:varchar(20);uniqueIndex:idx_deleted_at_phone_number"`
	Password    string                `gorm:"not null;type:varchar(255)"`
	CompanyID   uuid.UUID             `gorm:"type:uuid;not null"`
	Company     MasterCompany         `gorm:"foreignKey:CompanyID"`
	RoleID      uuid.UUID             `gorm:"type:uuid;not null"`
	Role        MasterRole            `gorm:"foreignKey:RoleID"`
	CreatedAt   time.Time             `gorm:"default:now();type:timestamp"`
	UpdatedAt   time.Time             `gorm:"default:now();type:timestamp"`
	DeletedAt   soft_delete.DeletedAt `gorm:"softDelete:nano;default:0;uniqueIndex:idx_deleted_at_email;uniqueIndex:idx_deleted_at_phone_number"`
}

func (u *MasterUser) BeforeCreate(tx *gorm.DB) error {
	if u.ID == uuid.Nil {
		u.ID = uuid.New()
	}
	return nil
}

func (u *MasterUser) Validate() error {
    if (u.Email == nil || *u.Email == "") && (u.PhoneNumber == nil || *u.PhoneNumber == "") {
        return fmt.Errorf("either email or phone number must be provided")
    }
    return nil
}