package models

import "time"

type DevicePairing struct {
	BaseModel
	UserID          string `gorm:"index" json:"user_id"`
	MobileDeviceID  string `gorm:"index" json:"mobile_device_id"`
	BrowserDeviceID string `gorm:"index" json:"browser_device_id"`

	PairCode  string    `json:"pair_code"`
	Status    string    `json:"status"` // pending, paired, expired
	ExpiresAt time.Time `json:"expires_at"`

	User          User   `gorm:"foreignKey:UserID"`
	MobileDevice  Device `gorm:"foreignKey:MobileDeviceID"`
	BrowserDevice Device `gorm:"foreignKey:BrowserDeviceID"`
}
