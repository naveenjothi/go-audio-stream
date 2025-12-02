package models

import "time"

type UserListenHistory struct {
	BaseModel
	UserID           string    `json:"user_id"`
	SongID           string    `json:"song_id"`
	PlayedAt         time.Time `json:"played_at"`
	DurationPlayedMs int       `json:"duration_played_ms"`
	IsCompleted      bool      `json:"is_completed"`

	User User `gorm:"foreignKey:UserID"`
	Song Song `gorm:"foreignKey:SongID"`
}
