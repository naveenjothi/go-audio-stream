package handlers

import (
	"context"
	"go-audio-stream/internal/models"
	"net/http"
	"strings"

	"github.com/labstack/echo/v4"
	"github.com/labstack/gommon/random"
	"gorm.io/gorm"
)

func calculateUserName(email string) string {
	emailStr := strings.Split(email, "@")[0]
	randomNumber := random.String(4, random.Numeric)
	return emailStr + randomNumber
}

func CreateUserHandler(c echo.Context, db *gorm.DB) error {
	first_name := c.FormValue("first_name")
	last_name := c.FormValue("last_name")
	email := c.FormValue("email")
	mobile := c.FormValue("mobile")
	user_name := calculateUserName(email)
	ctx := context.Background()
	user := models.UserModel{FirstName: first_name, LastName: last_name, Email: email, Mobile: mobile, Username: user_name}
	err := gorm.G[models.UserModel](db).Create(ctx, &user)
	if err != nil {
		return err
	}
	return c.JSON(http.StatusCreated, user)
}

// func FindOneUserById(c echo.Context, db database.Service) error {
// 	id := c.Param("id")

// 	var user models.User

// 	result, err := db.Find(`SELECT id, email, first_name, last_name, mobile, user_name FROM users WHERE id=$1;`, id)
// 	if err != nil {
// 		return err
// 	}

// 	err = result.Scan(&user.Id, &user.Email, &user.FirstName, &user.LastName, &user.Mobile, &user.Username)

// 	if err != nil {
// 		return err
// 	}

// 	return c.JSON(http.StatusOK, user)
// }
