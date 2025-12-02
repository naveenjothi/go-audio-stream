package handlers

import (
	"go-audio-stream/pkg/database"
	"go-audio-stream/pkg/models"
	"net/http"
	"strings"

	"github.com/labstack/echo/v4"
	"github.com/labstack/gommon/random"
)

func calculateUserName(email string) string {
	emailStr := strings.Split(email, "@")[0]
	randomNumber := random.String(4, random.Numeric)
	return emailStr + randomNumber
}

func CreateUserHandler(c echo.Context, db database.Service) error {
	user := new(models.User)

	if err := c.Bind(user); err != nil {
		return c.JSON(http.StatusBadRequest, echo.Map{"error": err.Error()})
	}

	// Validate required fields
	if user.Email == "" {
		return c.JSON(http.StatusBadRequest, echo.Map{"error": "email is required"})
	}
	if user.FirstName == "" {
		return c.JSON(http.StatusBadRequest, echo.Map{"error": "first_name is required"})
	}

	user.Username = calculateUserName(user.Email)
	_, err := db.Create(user)
	if err != nil {
		return c.JSON(http.StatusInternalServerError, echo.Map{"error": err.Error()})
	}

	return c.JSON(http.StatusCreated, user)
}

func FindOneUserById(c echo.Context, db database.Service) error {
	id := c.Param("id")

	var user models.User

	_, err := db.Find(&user, "firebase_id = ?", id)

	if err != nil {
		return c.JSON(http.StatusInternalServerError, echo.Map{"error": err.Error()})
	}

	return c.JSON(http.StatusOK, user)
}

func UpdateUserHandler(c echo.Context, db database.Service) error {
	id := c.Param("id")
	user := new(models.User)

	_, err := db.Find(&user, "id = ?", id)

	if err != nil {
		return c.JSON(http.StatusInternalServerError, echo.Map{"error": err.Error()})
	}

	if err := c.Bind(user); err != nil {
		return c.JSON(http.StatusBadRequest, echo.Map{"error": err.Error()})
	}

	_, err = db.Update(&models.User{}, user, "id = ?", id)

	if err != nil {
		return c.JSON(http.StatusInternalServerError, echo.Map{"error": err.Error()})
	}

	return c.JSON(http.StatusOK, user)
}

func DeleteUserHandler(c echo.Context, db database.Service) error {
	id := c.Param("id")

	_, err := db.Delete(&models.User{}, "id = ?", id)

	if err != nil {
		return c.JSON(http.StatusInternalServerError, echo.Map{"error": err.Error()})
	}

	return c.JSON(http.StatusOK, echo.Map{"message": "User deleted successfully"})
}
