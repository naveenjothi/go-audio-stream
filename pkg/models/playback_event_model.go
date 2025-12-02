package models

import "gorm.io/datatypes"

type PlaybackEvent struct {
	BaseModel
	UserID   string `json:"user_id"`
	DeviceID string `json:"device_id"`
	SongID   string `json:"song_id"`

	EventType  string         `json:"event_type"` // PLAY, PAUSE, SEEK, etc.
	PositionMS int            `json:"position_ms"`
	Metadata   datatypes.JSON `json:"metadata"`

	User   User   `gorm:"foreignKey:UserID"`
	Device Device `gorm:"foreignKey:DeviceID"`
	Song   Song   `gorm:"foreignKey:SongID"`
}
