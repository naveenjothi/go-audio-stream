package migrations

import (
	"context"
	"fmt"
	"go-audio-stream/pkg/models"
	"go-audio-stream/pkg/storage"
	"log"
	"os"
	"path/filepath"

	"gorm.io/gorm"
)

// AddARRahmanShowkali adds AR Rahman artist and Showkali song
type AddARRahmanShowkali struct{}

func (m *AddARRahmanShowkali) Version() string {
	return "20251219230000"
}

func (m *AddARRahmanShowkali) Name() string {
	return "add_ar_rahman_showkali"
}

func (m *AddARRahmanShowkali) Up(db *gorm.DB) error {
	ctx := context.Background()

	// Initialize storage client for Backblaze B2
	storageClient, err := storage.NewClient(storage.LoadConfig())
	if err != nil {
		return fmt.Errorf("failed to initialize storage client: %w", err)
	}

	// Upload Showkali.mp3 from assets folder
	assetsPath := getAssetsPath()
	songFilePath := filepath.Join(assetsPath, "Showkali.mp3")

	log.Printf("Uploading audio file from: %s", songFilePath)

	songFile, err := os.Open(songFilePath)
	if err != nil {
		return fmt.Errorf("failed to open song file: %w", err)
	}
	defer songFile.Close()

	// Upload to artist/album/songs/ folder structure in B2
	// Path: /artist-slug/album-slug/songs/song-name.mp3
	// Note: B2 doesn't support per-object ACLs like S3, public access is controlled at bucket level
	songKey := "a-r-rahman/acham-yenbadhu-madamaiyada/songs/showkali.mp3"
	_, err = storageClient.Upload(ctx, songKey, songFile, "audio/mpeg")
	if err != nil {
		return fmt.Errorf("failed to upload song to B2: %w", err)
	}
	log.Printf("Uploaded song to B2: %s", songKey)

	// Get the public URL for the uploaded song
	songURL := storageClient.GetPublicURL(songKey)

	// Create AR Rahman artist
	// A.R. Rahman is an Indian composer, record producer, singer, and songwriter.
	// He has won two Academy Awards, two Grammy Awards, a BAFTA Award, and a Golden Globe Award.
	artist := models.Artist{
		Name:      "A.R. Rahman",
		Image:     fmt.Sprintf("https://%s.s3.%s.backblazeb2.com/artists/ar-rahman.jpg", os.Getenv("B2_BUCKET_NAME"), os.Getenv("B2_REGION")),
		Followers: 15000000,
		Verified:  true,
		Bio:       "Allah Rakha Rahman, professionally known as A.R. Rahman, is an Indian composer, record producer, singer, and songwriter. He is one of the world's best-selling music artists and has won two Academy Awards, two Grammy Awards, a BAFTA Award, a Golden Globe Award, six National Film Awards, and numerous Filmfare Awards.",
	}

	if err := db.Create(&artist).Error; err != nil {
		return fmt.Errorf("failed to create artist A.R. Rahman: %w", err)
	}
	log.Println("Created artist: A.R. Rahman")

	// Create Showkali song from movie "Acham Yenbadhu Madamaiyada" (2016)
	// Movie Details:
	// Starring: Silambarasan, Manjima Mohan, Baba Sehgal
	// Music: A.R. Rahman
	// Director: Gautham Vasudev Menon
	// Lyricists: Thamarai, Madhan Karky, Vignesh Shivan, ADK, Sri Raskol, Pavendhar Bharathidasan
	// Year: 2016
	// Language: Tamil
	trackNumber := int16(1)
	song := models.Song{
		Name:        "Showkali",
		Image:       fmt.Sprintf("https://%s.s3.%s.backblazeb2.com/albums/acham-yenbadhu-madamaiyada.jpg", os.Getenv("B2_BUCKET_NAME"), os.Getenv("B2_REGION")),
		URL:         songURL,
		VideoURL:    "",
		Duration:    270, // 4:30 in seconds
		Explicit:    false,
		TrackNumber: &trackNumber,
		Language:    "Tamil",
		Artists:     []models.Artist{artist},
	}

	if err := db.Create(&song).Error; err != nil {
		return fmt.Errorf("failed to create song Showkali: %w", err)
	}
	log.Println("Created song: Showkali from Acham Yenbadhu Madamaiyada (2016)")

	return nil
}

func (m *AddARRahmanShowkali) Down(db *gorm.DB) error {
	ctx := context.Background()

	// Initialize storage client for Backblaze B2
	storageClient, err := storage.NewClient(storage.LoadConfig())
	if err != nil {
		log.Printf("Warning: failed to initialize storage client for cleanup: %v", err)
	} else {
		// Delete uploaded song from B2 (using the same path structure)
		songKey := "a-r-rahman/acham-yenbadhu-madamaiyada/songs/showkali.mp3"
		if err := storageClient.Delete(ctx, songKey); err != nil {
			log.Printf("Warning: failed to delete song from B2: %v", err)
		} else {
			log.Printf("Deleted song from B2: %s", songKey)
		}
	}

	// Delete song first (due to foreign key)
	if err := db.Where("name = ?", "Showkali").Delete(&models.Song{}).Error; err != nil {
		return fmt.Errorf("failed to delete song Showkali: %w", err)
	}

	// Delete artist
	if err := db.Where("name = ?", "A.R. Rahman").Delete(&models.Artist{}).Error; err != nil {
		return fmt.Errorf("failed to delete artist A.R. Rahman: %w", err)
	}

	return nil
}

// getAssetsPath returns the path to the assets directory
// It looks for the assets folder relative to the project root
func getAssetsPath() string {
	// Try to find assets folder - check common locations
	possiblePaths := []string{
		"assets",             // Running from project root
		"../../assets",       // Running from services/migration
		"../../../../assets", // Running from services/migration/cmd
		"/Volumes/Work/projects/go-audio-stream/assets", // Absolute path as fallback
	}

	for _, path := range possiblePaths {
		if _, err := os.Stat(path); err == nil {
			absPath, _ := filepath.Abs(path)
			return absPath
		}
	}

	// Default to absolute path
	return "/Volumes/Work/projects/go-audio-stream/assets"
}
