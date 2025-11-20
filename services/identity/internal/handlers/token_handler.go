package handlers

import (
	"context"
	"go-audio-stream/pkg/database"
	"go-audio-stream/pkg/models"
	"log"
	"net/http"
	"strings"

	firebase "firebase.google.com/go/v4"
	"github.com/labstack/echo/v4"
)

func VerifyTokenHandler(c echo.Context, db database.Service, firebase_app *firebase.App) error {
	token := strings.Split(c.Request().Header.Get("Authorization"), " ")[1]

	client, err := firebase_app.Auth(context.Background())
	if err != nil {
		log.Fatalf("error getting Auth client: %v\n", err)
	}

	verified_token, err := client.VerifyIDToken(context.Background(), token)
	if err != nil {
		log.Fatalf("error verifying ID token: %v\n", err)
	}

	log.Printf("Verified ID token: %v\n", verified_token.Claims)

	var user models.User

	_, err = db.Find(&user, "email = ?", verified_token.Claims["email"])

	if err != nil {
		return c.JSON(http.StatusInternalServerError, echo.Map{"error": err.Error()})
	}

	if user.IsSuspended {
		return c.JSON(http.StatusForbidden, echo.Map{"error": "user is suspended"})
	}

	return c.JSON(http.StatusOK, user)
}
