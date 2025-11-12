package models

import "time"

// 1. Define the Junction Model
type PlaylistSong struct {
	// Foreign Key to Playlist (use the type of your BaseModel.ID)
	PlaylistID string `gorm:"primaryKey"`
	// Foreign Key to Song (use the type of your BaseModel.ID)
	SongID    string `gorm:"primaryKey"`
	Position  int    `gorm:"not null"`
	CreatedAt time.Time
}
