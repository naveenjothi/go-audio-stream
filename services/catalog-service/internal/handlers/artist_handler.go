package handlers

import (
	"go-audio-stream/pkg/database"
	"go-audio-stream/pkg/models"
	"net/http"

	"github.com/labstack/echo/v4"
)

func CreateArtistHandler(c echo.Context, db database.Service) error {

	artist := new(models.Artist)

	if err := c.Bind(artist); err != nil {
		return c.JSON(http.StatusBadRequest, echo.Map{"error": err.Error()})
	}

	_, err := db.Create(&artist)

	if err != nil {
		return c.JSON(http.StatusInternalServerError, echo.Map{"error": err.Error()})
	}

	return c.JSON(http.StatusCreated, artist)
}

func UpdateArtistHandler(c echo.Context, db database.Service) error {
	id := c.Param("id")
	artist := new(models.Artist)

	// Check if artist exists
	_, err := db.Find(artist, "id = ?", id)
	if err != nil {
		return c.JSON(http.StatusNotFound, echo.Map{"error": "Artist not found"})
	}

	// Bind new data
	if err := c.Bind(artist); err != nil {
		return c.JSON(http.StatusBadRequest, echo.Map{"error": err.Error()})
	}

	// Update
	_, err = db.Update(&models.Artist{}, artist, "id = ?", id)
	if err != nil {
		return c.JSON(http.StatusInternalServerError, echo.Map{"error": err.Error()})
	}

	return c.JSON(http.StatusOK, artist)
}

func DeleteArtistHandler(c echo.Context, db database.Service) error {
	id := c.Param("id")

	_, err := db.Delete(&models.Artist{}, "id = ?", id)
	if err != nil {
		return c.JSON(http.StatusInternalServerError, echo.Map{"error": err.Error()})
	}

	return c.JSON(http.StatusOK, echo.Map{"message": "Artist deleted successfully"})
}

func FindOneArtistById(c echo.Context, db database.Service) error {
	id := c.Param("id")
	var artist models.Artist

	_, err := db.Find(&artist, "id = ?", id)
	if err != nil {
		return c.JSON(http.StatusNotFound, echo.Map{"error": "Artist not found"})
	}

	return c.JSON(http.StatusOK, artist)
}

func FindAllArtists(c echo.Context, db database.Service) error {
	var artists []models.Artist
	_, err := db.Find(&artists)
	if err != nil {
		return c.JSON(http.StatusInternalServerError, echo.Map{"error": err.Error()})
	}

	return c.JSON(http.StatusOK, artists)
}
