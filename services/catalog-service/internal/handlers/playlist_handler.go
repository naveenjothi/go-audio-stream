package handlers

import (
	"go-audio-stream/pkg/database"
	"go-audio-stream/pkg/models"
	"net/http"

	"github.com/labstack/echo/v4"
)

// CreatePlaylistHandler creates a new playlist.
// @Summary      Create a new playlist
// @Description  Create a new playlist with the provided details
// @Tags         playlists
// @Accept       json
// @Produce      json
// @Param        playlist  body      models.Playlist  true  "Playlist Data"
// @Success      201       {object}  models.Playlist
// @Failure      400       {object}  map[string]string
// @Failure      500       {object}  map[string]string
// @Router       /api/v1/playlists/ [post]
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

// UpdatePlaylistHandler updates an existing playlist.
// @Summary      Update a playlist
// @Description  Update a playlist's details
// @Tags         playlists
// @Accept       json
// @Produce      json
// @Param        id        path      string           true  "Playlist ID"
// @Param        playlist  body      models.Playlist  true  "Playlist Data"
// @Success      200       {object}  models.Playlist
// @Failure      400       {object}  map[string]string
// @Failure      404       {object}  map[string]string
// @Failure      500       {object}  map[string]string
// @Router       /api/v1/playlists/{id} [put]
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

// DeletePlaylistHandler deletes a playlist.
// @Summary      Delete a playlist
// @Description  Delete a playlist by ID
// @Tags         playlists
// @Accept       json
// @Produce      json
// @Param        id   path      string  true  "Playlist ID"
// @Success      200  {object}  map[string]string
// @Failure      500  {object}  map[string]string
// @Router       /api/v1/playlists/{id} [delete]
func DeletePlaylistHandler(c echo.Context, db database.Service) error {
	id := c.Param("id")

	_, err := db.Delete(&models.Playlist{}, "id = ?", id)
	if err != nil {
		return c.JSON(http.StatusInternalServerError, echo.Map{"error": err.Error()})
	}

	return c.JSON(http.StatusOK, echo.Map{"message": "Playlist deleted successfully"})
}

// FindOnePlaylistById retrieves a playlist by ID.
// @Summary      Get a playlist
// @Description  Get a playlist by ID
// @Tags         playlists
// @Accept       json
// @Produce      json
// @Param        id   path      string  true  "Playlist ID"
// @Success      200  {object}  models.Playlist
// @Failure      404  {object}  map[string]string
// @Router       /api/v1/playlists/{id} [get]
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

// FindAllPlaylists retrieves all playlists.
// @Summary      Get all playlists
// @Description  Get a list of all playlists
// @Tags         playlists
// @Accept       json
// @Produce      json
// @Success      200  {array}   models.Playlist
// @Failure      500  {object}  map[string]string
// @Router       /api/v1/playlists/ [get]
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

// AddSongToPlaylistHandler adds a song to a playlist.
// @Summary      Add song to playlist
// @Description  Add a song to a playlist
// @Tags         playlists
// @Accept       json
// @Produce      json
// @Param        id   path      string          true  "Playlist ID"
// @Param        req  body      AddSongRequest  true  "Song Data"
// @Success      201  {object}  models.PlaylistSong
// @Failure      400  {object}  map[string]string
// @Failure      500  {object}  map[string]string
// @Router       /api/v1/playlists/{id}/songs [post]
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

// RemoveSongFromPlaylistHandler removes a song from a playlist.
// @Summary      Remove song from playlist
// @Description  Remove a song from a playlist
// @Tags         playlists
// @Accept       json
// @Produce      json
// @Param        id       path      string  true  "Playlist ID"
// @Param        song_id  path      string  true  "Song ID"
// @Success      200      {object}  map[string]string
// @Failure      500      {object}  map[string]string
// @Router       /api/v1/playlists/{id}/songs/{song_id} [delete]
func RemoveSongFromPlaylistHandler(c echo.Context, db database.Service) error {
	playlistID := c.Param("id")
	songID := c.Param("song_id")

	_, err := db.Delete(&models.PlaylistSong{}, "playlist_id = ? AND song_id = ?", playlistID, songID)
	if err != nil {
		return c.JSON(http.StatusInternalServerError, echo.Map{"error": err.Error()})
	}

	return c.JSON(http.StatusOK, echo.Map{"message": "Song removed from playlist"})
}
