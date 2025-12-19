package main

import (
	"flag"
	"log"
	"os"

	"go-audio-stream/pkg/database"
	"go-audio-stream/services/migration/internal/migrations"
	"go-audio-stream/services/migration/internal/runner"

	_ "github.com/joho/godotenv/autoload"
)

func main() {
	// Parse command line flags
	command := flag.String("cmd", "status", "Migration command: up, down, status")
	flag.Parse()

	// Initialize database connection
	db := database.New()
	defer db.Close()

	// Create migration runner
	r := runner.NewRunner(db.GetDB(), migrations.GetMigrations())

	// Execute command
	var err error
	switch *command {
	case "up":
		log.Println("Running migrations...")
		err = r.Run()
	case "down":
		log.Println("Rolling back last migration...")
		err = r.Rollback()
	case "status":
		err = r.Status()
	default:
		log.Printf("Unknown command: %s", *command)
		log.Println("Available commands: up, down, status")
		os.Exit(1)
	}

	if err != nil {
		log.Fatalf("Migration failed: %v", err)
	}
}
