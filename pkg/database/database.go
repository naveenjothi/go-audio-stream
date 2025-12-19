package database

import (
	"context"
	"fmt"
	"go-audio-stream/pkg/models"
	"log"
	"os"
	"strconv"
	"time"

	_ "github.com/jackc/pgx/v5/stdlib"
	_ "github.com/joho/godotenv/autoload"
	"gorm.io/driver/postgres"
	"gorm.io/gorm"
)

// Service represents a service that interacts with a database.
type Service interface {
	// Health returns a map of health status information.
	// The keys and values in the map are service-specific.
	Health() map[string]string

	// Close terminates the database connection.
	// It returns an error if the connection cannot be closed.
	Close() error
	Create(value interface{}) (*gorm.DB, error)
	Update(model interface{}, updates interface{}, where interface{}, whereArgs ...interface{}) (*gorm.DB, error)
	Find(dest interface{}, conditions ...interface{}) (*gorm.DB, error)
	Delete(model interface{}, where interface{}, whereArgs ...interface{}) (*gorm.DB, error)
	GetDB() *gorm.DB
}

type service struct {
	gorm_db *gorm.DB
}

var (
	database   = os.Getenv("BLUEPRINT_DB_DATABASE")
	password   = os.Getenv("BLUEPRINT_DB_PASSWORD")
	username   = os.Getenv("BLUEPRINT_DB_USERNAME")
	port       = os.Getenv("BLUEPRINT_DB_PORT")
	host       = os.Getenv("BLUEPRINT_DB_HOST")
	schema     = os.Getenv("BLUEPRINT_DB_SCHEMA")
	dbInstance *service
)

func New() Service {
	// Reuse Connection
	if dbInstance != nil {
		return dbInstance
	}

	dsn := fmt.Sprintf("host=%s user=%s password=%s dbname=%s port=%s sslmode=disable TimeZone=Asia/Shanghai search_path=%s", host, username, password, database, port, schema)
	gorm_db, err := gorm.Open(postgres.Open(dsn), &gorm.Config{})
	if err != nil {
		log.Fatal(err)
	}

	gorm_db.Exec("CREATE EXTENSION IF NOT EXISTS vector")

	gorm_db.AutoMigrate(
		&models.User{},
		&models.Artist{},
		&models.Song{},
		&models.Playlist{},
		&models.PlaylistSong{},
		&models.Device{},
		&models.DevicePairing{},
		&models.UserListenHistory{},
		&models.PlaybackEvent{},
		&models.PlaybackSession{},
		&models.SongFeatures{},
		&models.SongInstrument{},
		&models.SongTag{},
		&models.UserLocalSong{},
		&models.SchemaMigration{},
	)

	dbInstance = &service{
		gorm_db: gorm_db,
	}

	return dbInstance
}

// Health checks the health of the database connection by pinging the database.
// It returns a map with keys indicating various health statistics.
func (s *service) Health() map[string]string {
	ctx, cancel := context.WithTimeout(context.Background(), 1*time.Second)
	defer cancel()

	stats := make(map[string]string)

	// Ping the database
	sqlDB, err := s.gorm_db.DB()
	if err != nil {
		stats["status"] = "down"
		stats["error"] = fmt.Sprintf("db down: %v", err)
		log.Fatalf("db down: %v", err) // Log the error and terminate the program
		return stats
	}
	err = sqlDB.PingContext(ctx)
	if err != nil {
		stats["status"] = "down"
		stats["error"] = fmt.Sprintf("db down: %v", err)
		log.Fatalf("db down: %v", err) // Log the error and terminate the program
		return stats
	}

	// Database is up, add more statistics
	stats["status"] = "up"
	stats["message"] = "It's healthy"

	// Get database stats (like open connections, in use, idle, etc.)
	dbStats := sqlDB.Stats()
	stats["open_connections"] = strconv.Itoa(dbStats.OpenConnections)
	stats["in_use"] = strconv.Itoa(dbStats.InUse)
	stats["idle"] = strconv.Itoa(dbStats.Idle)
	stats["wait_count"] = strconv.FormatInt(dbStats.WaitCount, 10)
	stats["wait_duration"] = dbStats.WaitDuration.String()
	stats["max_idle_closed"] = strconv.FormatInt(dbStats.MaxIdleClosed, 10)
	stats["max_lifetime_closed"] = strconv.FormatInt(dbStats.MaxLifetimeClosed, 10)

	// Evaluate stats to provide a health message
	if dbStats.OpenConnections > 40 { // Assuming 50 is the max for this example
		stats["message"] = "The database is experiencing heavy load."
	}

	if dbStats.WaitCount > 1000 {
		stats["message"] = "The database has a high number of wait events, indicating potential bottlenecks."
	}

	if dbStats.MaxIdleClosed > int64(dbStats.OpenConnections)/2 {
		stats["message"] = "Many idle connections are being closed, consider revising the connection pool settings."
	}

	if dbStats.MaxLifetimeClosed > int64(dbStats.OpenConnections)/2 {
		stats["message"] = "Many connections are being closed due to max lifetime, consider increasing max lifetime or revising the connection usage pattern."
	}

	return stats
}

func (s *service) Close() error {
	log.Printf("Disconnected from database: %s", database)
	sqlDB, err := s.gorm_db.DB()
	if err != nil {
		return err
	}
	return sqlDB.Close()
}

func (s *service) Create(value interface{}) (*gorm.DB, error) {
	result := s.gorm_db.Create(value)
	if result.Error != nil {
		return nil, fmt.Errorf("failed to create record: %w", result.Error)
	}
	return result, nil
}

// Parameters:
//   - model: The GORM model instance or a pointer to it (used for table name).
//   - updates: The data to update (e.g., map[string]interface{} or a struct).
//   - where: The WHERE clause condition (e.g., "id = ?", or a struct/map).
//   - whereArgs: Optional arguments for the WHERE clause.
//
// The function signature is modified to accept the dynamic 'where' and 'whereArgs'.
func (s *service) Update(model interface{}, updates interface{}, where interface{}, whereArgs ...interface{}) (*gorm.DB, error) {
	// 1. Start with the model to scope the query.
	db := s.gorm_db.Model(model)

	// 2. Apply the dynamic WHERE filter.
	// This uses the provided 'where' expression and 'whereArgs' to construct the filter.
	db = db.Where(where, whereArgs...)

	result := db.Updates(updates)

	if result.Error != nil {
		return nil, fmt.Errorf("failed to update record: %w", result.Error)
	}

	return result, nil
}

func (s *service) Find(dest interface{}, conditions ...interface{}) (*gorm.DB, error) {
	result := s.gorm_db.Find(dest, conditions...)
	if result.Error != nil {
		return nil, fmt.Errorf("failed to find records: %w", result.Error)
	}
	return result, nil
}

func (s *service) Delete(model interface{}, where interface{}, whereArgs ...interface{}) (*gorm.DB, error) {
	// 1. Start with the model to scope the query.
	db := s.gorm_db.Model(model)

	// 2. Apply the dynamic WHERE filter.
	// This uses the provided 'where' expression and 'whereArgs' to construct the filter.
	db = db.Where(where, whereArgs...)

	result := db.Delete(model)

	if result.Error != nil {
		return nil, fmt.Errorf("failed to delete record: %w", result.Error)
	}

	return result, nil
}

func (s *service) GetDB() *gorm.DB {
	return s.gorm_db
}
