package handlers

import (
	"context"
	"go-audio-stream/pkg/database"
	"log"

	firebase "firebase.google.com/go/v4"
	"github.com/labstack/echo/v4"
)

func LoginHandler(c echo.Context, db database.Service, firebase_app *firebase.App) error {
	client, err := firebase_app.Auth(context.Background())
	if err != nil {
		log.Fatalf("error getting Auth client: %v\n", err)
	}

	token, err := client.VerifyIDToken(context.Background(), "idToken")
	if err != nil {
		log.Fatalf("error verifying ID token: %v\n", err)
	}

	log.Printf("Verified ID token: %v\n", token)
	return nil
}
