package migrations

import "go-audio-stream/services/migration/internal/runner"

// GetMigrations returns all registered migrations in order
func GetMigrations() []runner.Migration {
	return []runner.Migration{
		&AddARRahmanShowkali{},
	}
}
