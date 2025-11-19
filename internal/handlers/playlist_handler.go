package handlers

import (
	"go-audio-stream/internal/database"
	"go-audio-stream/internal/models"
	"net/http"

	"github.com/labstack/echo/v4"
)

func CreatePlaylistHandler(c echo.Context, db database.Service) error {
	playlist := new(models.Playlist)

	if err := c.Bind(playlist); err != nil {
		return c.JSON(http.StatusBadRequest, echo.Map{"error": err.Error()})
	}

	_, err := db.Create(playlist)
	if err != nil {
		return c.JSON(http.StatusInternalServerError, echo.Map{"error": err.Error()})
	}

	return c.JSON(http.StatusCreated, playlist)
}

func UpdatePlaylistHandler(c echo.Context, db database.Service) error {
	id := c.Param("id")
	playlist := new(models.Playlist)

	_, err := db.Find(playlist, "id = ?", id)
	if err != nil {
		return c.JSON(http.StatusNotFound, echo.Map{"error": "Playlist not found"})
	}

	if err := c.Bind(playlist); err != nil {
		return c.JSON(http.StatusBadRequest, echo.Map{"error": err.Error()})
	}

	_, err = db.Update(&models.Playlist{}, playlist, "id = ?", id)
	if err != nil {
		return c.JSON(http.StatusInternalServerError, echo.Map{"error": err.Error()})
	}

	return c.JSON(http.StatusOK, playlist)
}

func DeletePlaylistHandler(c echo.Context, db database.Service) error {
	id := c.Param("id")

	_, err := db.Delete(&models.Playlist{}, "id = ?", id)
	if err != nil {
		return c.JSON(http.StatusInternalServerError, echo.Map{"error": err.Error()})
	}

	return c.JSON(http.StatusOK, echo.Map{"message": "Playlist deleted successfully"})
}

func FindOnePlaylistById(c echo.Context, db database.Service) error {
	id := c.Param("id")
	var playlist models.Playlist

	// Preload songs if needed, but for now just basic info
	_, err := db.Find(&playlist, "id = ?", id)
	if err != nil {
		return c.JSON(http.StatusNotFound, echo.Map{"error": "Playlist not found"})
	}

	return c.JSON(http.StatusOK, playlist)
}

func FindAllPlaylists(c echo.Context, db database.Service) error {
	var playlists []models.Playlist
	_, err := db.Find(&playlists)
	if err != nil {
		return c.JSON(http.StatusInternalServerError, echo.Map{"error": err.Error()})
	}

	return c.JSON(http.StatusOK, playlists)
}

type AddSongRequest struct {
	SongID   string `json:"song_id"`
	Position int    `json:"position"`
}

func AddSongToPlaylistHandler(c echo.Context, db database.Service) error {
	playlistID := c.Param("id")
	req := new(AddSongRequest)

	if err := c.Bind(req); err != nil {
		return c.JSON(http.StatusBadRequest, echo.Map{"error": err.Error()})
	}

	playlistSong := &models.PlaylistSong{
		PlaylistID: playlistID,
		SongID:     req.SongID,
		Position:   req.Position,
	}

	_, err := db.Create(playlistSong)
	if err != nil {
		return c.JSON(http.StatusInternalServerError, echo.Map{"error": err.Error()})
	}

	return c.JSON(http.StatusCreated, playlistSong)
}

func RemoveSongFromPlaylistHandler(c echo.Context, db database.Service) error {
	playlistID := c.Param("id")
	songID := c.Param("song_id")

	_, err := db.Delete(&models.PlaylistSong{}, "playlist_id = ? AND song_id = ?", playlistID, songID)
	if err != nil {
		return c.JSON(http.StatusInternalServerError, echo.Map{"error": err.Error()})
	}

	return c.JSON(http.StatusOK, echo.Map{"message": "Song removed from playlist"})
}
