package models

import "github.com/pgvector/pgvector-go"

type SongFeatures struct {
	BaseModel
	SongID string `gorm:"uniqueIndex" json:"song_id"`

	Tempo            float32 `json:"tempo"`
	Energy           float32 `json:"energy"`
	Valence          float32 `json:"valence"`
	Danceability     float32 `json:"danceability"`
	Loudness         float32 `json:"loudness"`
	Acousticness     float32 `json:"acousticness"`
	Speechiness      float32 `json:"speechiness"`
	Instrumentalness float32 `json:"instrumentalness"`
	DurationMs       int32   `json:"duration_ms"`

	Embedding pgvector.Vector `gorm:"type:vector(128)" json:"embedding"`

	Song Song `gorm:"foreignKey:SongID"`
}
