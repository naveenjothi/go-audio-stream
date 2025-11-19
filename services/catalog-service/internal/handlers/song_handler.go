package handlers

import (
	"go-audio-stream/pkg/database"
	"go-audio-stream/pkg/models"
	"net/http"

	"github.com/labstack/echo/v4"
)

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

func DeleteSongHandler(c echo.Context, db database.Service) error {
	id := c.Param("id")

	_, err := db.Delete(&models.Song{}, "id = ?", id)
	if err != nil {
		return c.JSON(http.StatusInternalServerError, echo.Map{"error": err.Error()})
	}

	return c.JSON(http.StatusOK, echo.Map{"message": "Song deleted successfully"})
}

func FindOneSongById(c echo.Context, db database.Service) error {
	id := c.Param("id")
	var song models.Song

	_, err := db.Find(&song, "id = ?", id)
	if err != nil {
		return c.JSON(http.StatusNotFound, echo.Map{"error": "Song not found"})
	}

	return c.JSON(http.StatusOK, song)
}

func FindAllSongs(c echo.Context, db database.Service) error {
	var songs []models.Song
	_, err := db.Find(&songs)
	if err != nil {
		return c.JSON(http.StatusInternalServerError, echo.Map{"error": err.Error()})
	}

	return c.JSON(http.StatusOK, songs)
}
