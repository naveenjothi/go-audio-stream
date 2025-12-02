package models

type UserLocalSong struct {
	BaseModel
	UserID   string `json:"user_id"`
	DeviceID string `json:"device_id"`

	Title      string `json:"title"`
	Artist     string `json:"artist"`
	Album      string `json:"album"`
	DurationMS int    `json:"duration_ms"`
	FilePath   string `json:"file_path"`
	Language   string `json:"language"`

	User   User   `gorm:"foreignKey:UserID"`
	Device Device `gorm:"foreignKey:DeviceID"`
}
