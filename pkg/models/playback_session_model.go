package models

type PlaybackSession struct {
	BaseModel
	UserID         string `json:"user_id"`
	SongID         string `json:"song_id"`
	ActiveDeviceID string `json:"active_device_id"`

	Status     string  `json:"status"`
	PositionMS int     `json:"position_ms"`
	Volume     float32 `json:"volume"`

	User   User   `gorm:"foreignKey:UserID"`
	Song   Song   `gorm:"foreignKey:SongID"`
	Device Device `gorm:"foreignKey:ActiveDeviceID"`
}
