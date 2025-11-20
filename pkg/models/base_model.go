package models

import (
	"time"

	"github.com/google/uuid"
	"gorm.io/gorm"
)

type BaseModel struct {
	ID          string          `json:"id" gorm:"type:uuid;primary_key;"`
	CreatedAt   time.Time       `json:"created_at"`
	UpdatedAt   time.Time       `json:"updated_at"`
	DeletedAt   *gorm.DeletedAt `json:"deleted_at,omitempty" gorm:"index"`
	IsSuspended bool            `json:"is_suspended"`
}

func (u *BaseModel) BeforeCreate(tx *gorm.DB) (err error) {
	u.ID = uuid.New().String()
	return
}
