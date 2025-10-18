package entity

import (
	"time"

	"github.com/google/uuid"
	"gorm.io/gorm"
	"gorm.io/plugin/soft_delete"
)

type RoleName string

const (
	RoleOwner RoleName = "owner"
	RoleAdmin RoleName = "admin"
	RoleStaff RoleName = "staff"
)

type MasterRole struct {
	ID        uuid.UUID             `gorm:"type:uuid;primaryKey"`
	Role      RoleName              `gorm:"type:enum('owner','admin','staff');unique;not null"`
	CreatedAt time.Time             `gorm:"default:now();type:timestamp"`
	UpdatedAt time.Time             `gorm:"default:now();type:timestamp"`
	DeletedAt soft_delete.DeletedAt `gorm:"softDelete:nano;default:0"`
	Users     []MasterUser          `gorm:"foreignKey:RoleID"`
}

func (r *MasterRole) BeforeCreate(tx *gorm.DB) error {
	if r.ID == uuid.Nil {
		r.ID = uuid.New()
	}
	return nil
}
