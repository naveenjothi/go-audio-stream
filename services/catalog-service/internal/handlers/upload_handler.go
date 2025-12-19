package handlers

import (
	"fmt"
	"io"
	"net/http"
	"path/filepath"
	"strings"
	"time"

	"go-audio-stream/pkg/storage"

	"github.com/google/uuid"
	"github.com/labstack/echo/v4"
)

// UploadHandler holds the storage client for upload operations
type UploadHandler struct {
	storage *storage.Client
}

// NewUploadHandler creates a new upload handler
func NewUploadHandler(storageClient *storage.Client) *UploadHandler {
	return &UploadHandler{
		storage: storageClient,
	}
}

// UploadAudioRequest contains metadata for audio upload
type UploadAudioRequest struct {
	SongID string `form:"song_id" json:"song_id"`
}

// UploadResponse is returned after successful upload
type UploadResponse struct {
	Key         string `json:"key"`
	URL         string `json:"url"`
	ContentType string `json:"content_type"`
	Size        int64  `json:"size"`
}

// UploadAudio handles audio file uploads
// POST /api/upload/audio
func (h *UploadHandler) UploadAudio(c echo.Context) error {
	// Get song ID from form
	songID := c.FormValue("song_id")
	if songID == "" {
		songID = uuid.New().String()
	}

	// Get the file from the request
	file, err := c.FormFile("file")
	if err != nil {
		return c.JSON(http.StatusBadRequest, echo.Map{"error": "No file provided"})
	}

	// Validate audio content type
	contentType := file.Header.Get("Content-Type")
	if !isValidAudioType(contentType) {
		return c.JSON(http.StatusBadRequest, echo.Map{
			"error": fmt.Sprintf("Invalid content type: %s. Allowed: audio/mpeg, audio/wav, audio/flac, audio/aac, audio/ogg", contentType),
		})
	}

	// Open the file
	src, err := file.Open()
	if err != nil {
		return c.JSON(http.StatusInternalServerError, echo.Map{"error": "Failed to read file"})
	}
	defer src.Close()

	// Generate storage key: songs/{song_id}/audio.{ext}
	ext := filepath.Ext(file.Filename)
	if ext == "" {
		ext = getExtensionFromContentType(contentType)
	}
	key := fmt.Sprintf("songs/%s/audio%s", songID, ext)

	// Upload to B2 (audio is private, use presigned URLs for access)
	uploadedKey, err := h.storage.Upload(c.Request().Context(), key, src, contentType)
	if err != nil {
		return c.JSON(http.StatusInternalServerError, echo.Map{"error": "Upload failed: " + err.Error()})
	}

	return c.JSON(http.StatusCreated, UploadResponse{
		Key:         uploadedKey,
		URL:         "", // Audio files use presigned URLs, not direct access
		ContentType: contentType,
		Size:        file.Size,
	})
}

// UploadImage handles image file uploads (album art, artist images)
// POST /api/upload/image
func (h *UploadHandler) UploadImage(c echo.Context) error {
	// Get entity info from form
	entityType := c.FormValue("entity_type") // "song", "artist", "playlist"
	entityID := c.FormValue("entity_id")

	if entityType == "" || entityID == "" {
		return c.JSON(http.StatusBadRequest, echo.Map{"error": "entity_type and entity_id are required"})
	}

	// Validate entity type
	if entityType != "song" && entityType != "artist" && entityType != "playlist" {
		return c.JSON(http.StatusBadRequest, echo.Map{"error": "Invalid entity_type. Allowed: song, artist, playlist"})
	}

	// Get the file from the request
	file, err := c.FormFile("file")
	if err != nil {
		return c.JSON(http.StatusBadRequest, echo.Map{"error": "No file provided"})
	}

	// Validate image content type
	contentType := file.Header.Get("Content-Type")
	if !isValidImageType(contentType) {
		return c.JSON(http.StatusBadRequest, echo.Map{
			"error": fmt.Sprintf("Invalid content type: %s. Allowed: image/jpeg, image/png, image/webp, image/gif", contentType),
		})
	}

	// Open the file
	src, err := file.Open()
	if err != nil {
		return c.JSON(http.StatusInternalServerError, echo.Map{"error": "Failed to read file"})
	}
	defer src.Close()

	// Generate storage key based on entity: {entity_type}s/{entity_id}/cover.{ext}
	ext := filepath.Ext(file.Filename)
	if ext == "" {
		ext = getImageExtensionFromContentType(contentType)
	}
	key := fmt.Sprintf("%ss/%s/cover%s", entityType, entityID, ext)

	// Upload to B2 with public-read ACL (images are public)
	uploadedKey, err := h.storage.UploadWithACL(c.Request().Context(), key, src, contentType, "public-read")
	if err != nil {
		return c.JSON(http.StatusInternalServerError, echo.Map{"error": "Upload failed: " + err.Error()})
	}

	// Return public URL for images
	publicURL := h.storage.GetPublicURL(uploadedKey)

	return c.JSON(http.StatusCreated, UploadResponse{
		Key:         uploadedKey,
		URL:         publicURL,
		ContentType: contentType,
		Size:        file.Size,
	})
}

// GetPresignedURL generates a time-limited download URL for audio files
// GET /api/files/:key/presigned
func (h *UploadHandler) GetPresignedURL(c echo.Context) error {
	key := c.Param("key")
	if key == "" {
		return c.JSON(http.StatusBadRequest, echo.Map{"error": "Key is required"})
	}

	// Decode URL-encoded key (e.g., songs%2F123%2Faudio.mp3 -> songs/123/audio.mp3)
	// The key comes from path param, which is already decoded by Echo

	// Generate presigned URL valid for 1 hour
	presignedURL, err := h.storage.GetPresignedURL(c.Request().Context(), key, 1*time.Hour)
	if err != nil {
		return c.JSON(http.StatusInternalServerError, echo.Map{"error": "Failed to generate URL: " + err.Error()})
	}

	return c.JSON(http.StatusOK, echo.Map{
		"url":        presignedURL,
		"expires_in": 3600, // seconds
	})
}

// DeleteFile removes a file from storage
// DELETE /api/files/:key
func (h *UploadHandler) DeleteFile(c echo.Context) error {
	key := c.Param("key")
	if key == "" {
		return c.JSON(http.StatusBadRequest, echo.Map{"error": "Key is required"})
	}

	err := h.storage.Delete(c.Request().Context(), key)
	if err != nil {
		return c.JSON(http.StatusInternalServerError, echo.Map{"error": "Failed to delete: " + err.Error()})
	}

	return c.JSON(http.StatusOK, echo.Map{"message": "File deleted successfully"})
}

// StreamAudio streams audio file with range support for seeking
// GET /api/stream/:key
func (h *UploadHandler) StreamAudio(c echo.Context) error {
	key := c.Param("key")
	if key == "" {
		return c.JSON(http.StatusBadRequest, echo.Map{"error": "Key is required"})
	}

	// Generate a short-lived presigned URL and redirect
	presignedURL, err := h.storage.GetPresignedURL(c.Request().Context(), key, 5*time.Minute)
	if err != nil {
		return c.JSON(http.StatusInternalServerError, echo.Map{"error": "Failed to generate stream URL: " + err.Error()})
	}

	return c.Redirect(http.StatusTemporaryRedirect, presignedURL)
}

// Helper functions

func isValidAudioType(contentType string) bool {
	validTypes := []string{
		"audio/mpeg",
		"audio/mp3",
		"audio/wav",
		"audio/wave",
		"audio/x-wav",
		"audio/flac",
		"audio/x-flac",
		"audio/aac",
		"audio/ogg",
		"audio/vorbis",
		"audio/webm",
	}
	for _, t := range validTypes {
		if strings.EqualFold(contentType, t) {
			return true
		}
	}
	return false
}

func isValidImageType(contentType string) bool {
	validTypes := []string{
		"image/jpeg",
		"image/jpg",
		"image/png",
		"image/webp",
		"image/gif",
	}
	for _, t := range validTypes {
		if strings.EqualFold(contentType, t) {
			return true
		}
	}
	return false
}

func getExtensionFromContentType(contentType string) string {
	switch strings.ToLower(contentType) {
	case "audio/mpeg", "audio/mp3":
		return ".mp3"
	case "audio/wav", "audio/wave", "audio/x-wav":
		return ".wav"
	case "audio/flac", "audio/x-flac":
		return ".flac"
	case "audio/aac":
		return ".aac"
	case "audio/ogg", "audio/vorbis":
		return ".ogg"
	case "audio/webm":
		return ".webm"
	default:
		return ".mp3"
	}
}

func getImageExtensionFromContentType(contentType string) string {
	switch strings.ToLower(contentType) {
	case "image/jpeg", "image/jpg":
		return ".jpg"
	case "image/png":
		return ".png"
	case "image/webp":
		return ".webp"
	case "image/gif":
		return ".gif"
	default:
		return ".jpg"
	}
}

// Utility to read up to a limit (for optional use)
func limitReader(r io.Reader, limit int64) io.Reader {
	return io.LimitReader(r, limit)
}
