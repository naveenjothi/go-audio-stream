package models

import "time"

// 1. Define the Junction Model
type PlaylistSong struct {
	PlaylistID string    `gorm:"primaryKey" json:"playlist_id"`
	SongID     string    `gorm:"primaryKey" json:"song_id"`
	Position   int       `gorm:"not null" json:"position"`
	CreatedAt  time.Time `json:"created_at"`
}
