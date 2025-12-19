package handlers

import (
	"go-audio-stream/pkg/database"
	"go-audio-stream/pkg/models"
	"net/http"

	"github.com/labstack/echo/v4"
)

// CreateSongHandler creates a new song.
// @Summary      Create a new song
// @Description  Create a new song with the provided details
// @Tags         songs
// @Accept       json
// @Produce      json
// @Param        song  body      models.Song  true  "Song Data"
// @Success      201   {object}  models.Song
// @Failure      400   {object}  map[string]string
// @Failure      500   {object}  map[string]string
// @Router       /api/v1/songs/ [post]
func CreateSongHandler(c echo.Context, db database.Service) error {
	song := new(models.Song)

	if err := c.Bind(song); err != nil {
		return c.JSON(http.StatusBadRequest, echo.Map{"error": err.Error()})
	}

	_, err := db.Create(song)
	if err != nil {
		return c.JSON(http.StatusInternalServerError, echo.Map{"error": err.Error()})
	}

	return c.JSON(http.StatusCreated, song)
}

// UpdateSongHandler updates an existing song.
// @Summary      Update a song
// @Description  Update a song's details
// @Tags         songs
// @Accept       json
// @Produce      json
// @Param        id    path      string       true  "Song ID"
// @Param        song  body      models.Song  true  "Song Data"
// @Success      200   {object}  models.Song
// @Failure      400   {object}  map[string]string
// @Failure      404   {object}  map[string]string
// @Failure      500   {object}  map[string]string
// @Router       /api/v1/songs/{id} [put]
func UpdateSongHandler(c echo.Context, db database.Service) error {
	id := c.Param("id")
	song := new(models.Song)

	_, err := db.Find(song, "id = ?", id)
	if err != nil {
		return c.JSON(http.StatusNotFound, echo.Map{"error": "Song not found"})
	}

	if err := c.Bind(song); err != nil {
		return c.JSON(http.StatusBadRequest, echo.Map{"error": err.Error()})
	}

	_, err = db.Update(&models.Song{}, song, "id = ?", id)
	if err != nil {
		return c.JSON(http.StatusInternalServerError, echo.Map{"error": err.Error()})
	}

	return c.JSON(http.StatusOK, song)
}

// DeleteSongHandler deletes a song.
// @Summary      Delete a song
// @Description  Delete a song by ID
// @Tags         songs
// @Accept       json
// @Produce      json
// @Param        id   path      string  true  "Song ID"
// @Success      200  {object}  map[string]string
// @Failure      500  {object}  map[string]string
// @Router       /api/v1/songs/{id} [delete]
func DeleteSongHandler(c echo.Context, db database.Service) error {
	id := c.Param("id")

	_, err := db.Delete(&models.Song{}, "id = ?", id)
	if err != nil {
		return c.JSON(http.StatusInternalServerError, echo.Map{"error": err.Error()})
	}

	return c.JSON(http.StatusOK, echo.Map{"message": "Song deleted successfully"})
}

// FindOneSongById retrieves a song by ID.
// @Summary      Get a song
// @Description  Get a song by ID
// @Tags         songs
// @Accept       json
// @Produce      json
// @Param        id   path      string  true  "Song ID"
// @Success      200  {object}  models.Song
// @Failure      404  {object}  map[string]string
// @Router       /api/v1/songs/{id} [get]
func FindOneSongById(c echo.Context, db database.Service) error {
	id := c.Param("id")
	var song models.Song

	_, err := db.Find(&song, "id = ?", id)
	if err != nil {
		return c.JSON(http.StatusNotFound, echo.Map{"error": "Song not found"})
	}

	return c.JSON(http.StatusOK, song)
}

// FindAllSongs retrieves all songs.
// @Summary      Get all songs
// @Description  Get a list of all songs
// @Tags         songs
// @Accept       json
// @Produce      json
// @Success      200  {array}   models.Song
// @Failure      500  {object}  map[string]string
// @Router       /api/v1/songs/ [get]
func FindAllSongs(c echo.Context, db database.Service) error {
	var songs []models.Song
	_, err := db.Find(&songs)
	if err != nil {
		return c.JSON(http.StatusInternalServerError, echo.Map{"error": err.Error()})
	}

	return c.JSON(http.StatusOK, songs)
}
