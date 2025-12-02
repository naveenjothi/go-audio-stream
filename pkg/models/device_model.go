package models

import "time"

type Device struct {
	BaseModel
	UserID       string     `gorm:"index" json:"user_id"`
	DeviceType   string     `json:"device_type"` // mobile / browser
	DeviceName   string     `json:"device_name"`
	LastOnlineAt *time.Time `json:"last_online_at"`

	User User `gorm:"foreignKey:UserID"`
}
