package models

import (
	"time"
)

type BaseModel struct {
	Id        string    `json:"id"`
	IsActive  bool      `json:"isActive"`
	IsDeleted bool      `json:"isDeleted"`
	CreatedAt time.Time `json:"createdAt"`
	UpdatedAt time.Time `json:"updatedAt"`
}

func (bm *BaseModel) InitiliseDefaultValue() {
	now := time.Now().UTC()
	bm.IsDeleted = false
	bm.IsActive = true
	bm.CreatedAt = now
	bm.UpdatedAt = now
}

func (bm *BaseModel) UpdateDefaultValue() {
	now := time.Now().UTC()
	bm.UpdatedAt = now
}
