package runner

import (
	"fmt"
	"go-audio-stream/pkg/models"
	"log"
	"sort"
	"time"

	"gorm.io/gorm"
)

// Migration defines the interface for database migrations
type Migration interface {
	Version() string
	Name() string
	Up(db *gorm.DB) error
	Down(db *gorm.DB) error
}

// Runner handles migration execution
type Runner struct {
	db         *gorm.DB
	migrations []Migration
}

// NewRunner creates a new migration runner
func NewRunner(db *gorm.DB, migrations []Migration) *Runner {
	// Sort migrations by version
	sort.Slice(migrations, func(i, j int) bool {
		return migrations[i].Version() < migrations[j].Version()
	})

	return &Runner{
		db:         db,
		migrations: migrations,
	}
}

// Run executes all pending migrations
func (r *Runner) Run() error {
	// Ensure schema_migrations table exists
	if err := r.db.AutoMigrate(&models.SchemaMigration{}); err != nil {
		return fmt.Errorf("failed to create schema_migrations table: %w", err)
	}

	applied, err := r.getAppliedMigrations()
	if err != nil {
		return err
	}

	pending := 0
	for _, m := range r.migrations {
		if _, ok := applied[m.Version()]; ok {
			continue
		}

		log.Printf("Running migration %s: %s", m.Version(), m.Name())

		if err := m.Up(r.db); err != nil {
			return fmt.Errorf("failed to run migration %s: %w", m.Version(), err)
		}

		// Record migration
		record := models.SchemaMigration{
			Version:   m.Version(),
			Name:      m.Name(),
			AppliedAt: time.Now(),
		}
		if err := r.db.Create(&record).Error; err != nil {
			return fmt.Errorf("failed to record migration %s: %w", m.Version(), err)
		}

		log.Printf("Completed migration %s", m.Version())
		pending++
	}

	if pending == 0 {
		log.Println("No pending migrations")
	} else {
		log.Printf("Applied %d migrations", pending)
	}

	return nil
}

// Rollback reverts the last applied migration
func (r *Runner) Rollback() error {
	var lastMigration models.SchemaMigration
	result := r.db.Order("version DESC").First(&lastMigration)
	if result.Error != nil {
		if result.Error == gorm.ErrRecordNotFound {
			log.Println("No migrations to rollback")
			return nil
		}
		return fmt.Errorf("failed to get last migration: %w", result.Error)
	}

	// Find the migration
	var migration Migration
	for _, m := range r.migrations {
		if m.Version() == lastMigration.Version {
			migration = m
			break
		}
	}

	if migration == nil {
		return fmt.Errorf("migration %s not found in registry", lastMigration.Version)
	}

	log.Printf("Rolling back migration %s: %s", migration.Version(), migration.Name())

	if err := migration.Down(r.db); err != nil {
		return fmt.Errorf("failed to rollback migration %s: %w", migration.Version(), err)
	}

	// Remove migration record
	if err := r.db.Delete(&lastMigration).Error; err != nil {
		return fmt.Errorf("failed to remove migration record %s: %w", migration.Version(), err)
	}

	log.Printf("Rolled back migration %s", migration.Version())
	return nil
}

// Status prints the status of all migrations
func (r *Runner) Status() error {
	// Ensure schema_migrations table exists
	if err := r.db.AutoMigrate(&models.SchemaMigration{}); err != nil {
		return fmt.Errorf("failed to create schema_migrations table: %w", err)
	}

	applied, err := r.getAppliedMigrations()
	if err != nil {
		return err
	}

	fmt.Println("\n=== Migration Status ===")
	fmt.Printf("%-20s %-40s %s\n", "VERSION", "NAME", "STATUS")
	fmt.Println("-------------------------------------------------------------------")

	for _, m := range r.migrations {
		status := "Pending"
		if record, ok := applied[m.Version()]; ok {
			status = fmt.Sprintf("Applied at %s", record.AppliedAt.Format("2006-01-02 15:04:05"))
		}
		fmt.Printf("%-20s %-40s %s\n", m.Version(), m.Name(), status)
	}
	fmt.Println()

	return nil
}

func (r *Runner) getAppliedMigrations() (map[string]models.SchemaMigration, error) {
	var migrations []models.SchemaMigration
	if err := r.db.Find(&migrations).Error; err != nil {
		return nil, fmt.Errorf("failed to get applied migrations: %w", err)
	}

	applied := make(map[string]models.SchemaMigration)
	for _, m := range migrations {
		applied[m.Version] = m
	}

	return applied, nil
}
