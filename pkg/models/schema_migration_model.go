package models

import "time"

// SchemaMigration tracks applied database migrations
type SchemaMigration struct {
	Version   string    `gorm:"primaryKey" json:"version"`
	Name      string    `json:"name"`
	AppliedAt time.Time `json:"applied_at"`
}
