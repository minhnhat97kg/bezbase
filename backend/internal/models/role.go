package models

import (
	"time"

	"gorm.io/gorm"
)

type Role struct {
	ID          uint           `json:"id" gorm:"primaryKey"`
	Name        string         `json:"name" gorm:"uniqueIndex;not null;size:100"`
	DisplayName string         `json:"display_name" gorm:"not null;size:255"`
	Description string         `json:"description" gorm:"size:500"`
	IsSystem    bool           `json:"is_system" gorm:"default:false"`
	IsActive    bool           `json:"is_active" gorm:"default:true"`
	CreatedAt   time.Time      `json:"created_at"`
	UpdatedAt   time.Time      `json:"updated_at"`
	DeletedAt   gorm.DeletedAt `json:"-" gorm:"index"`
}



func (r *Role) BeforeDelete(tx *gorm.DB) error {
	if r.IsSystem {
		return gorm.ErrInvalidValue
	}
	return nil
}

func (r *Role) Validate() error {
	if r.Name == "" {
		return gorm.ErrInvalidValue
	}
	if r.DisplayName == "" {
		return gorm.ErrInvalidValue
	}
	return nil
}