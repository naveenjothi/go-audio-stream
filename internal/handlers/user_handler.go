package handlers

import (
	"go-audio-stream/internal/database"
	"go-audio-stream/internal/models"
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
	first_name := c.FormValue("first_name")
	last_name := c.FormValue("last_name")
	email := c.FormValue("email")
	mobile := c.FormValue("mobile")
	user_name := calculateUserName(email)

	result, err := db.Create(`INSERT INTO users (mobile, email, first_name, last_name, user_name)
	VALUES ($1, $2, $3, $4, $5)
	RETURNING id;`, mobile, email, first_name, last_name, user_name)

	if err != nil {
		return err
	}
	return c.JSON(http.StatusCreated, result)
}

func FindOneUserById(c echo.Context, db database.Service) error {
	id := c.Param("id")

	var user models.User

	result, err := db.Find(`SELECT id, email, first_name, last_name, mobile, user_name FROM users WHERE id=$1;`, id)
	if err != nil {
		return err
	}

	err = result.Scan(&user.Id, &user.Email, &user.FirstName, &user.LastName, &user.Mobile, &user.Username)

	if err != nil {
		return err
	}

	return c.JSON(http.StatusOK, user)
}
