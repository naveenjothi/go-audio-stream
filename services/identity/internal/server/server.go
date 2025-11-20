package server

import (
	"context"
	"log"

	firebase "firebase.google.com/go/v4"
	_ "github.com/joho/godotenv/autoload"

	"go-audio-stream/pkg/database"
)

type Server struct {
	DB database.Service

	FirebaseApp *firebase.App
}

func NewServer() *Server {
	app, err := firebase.NewApp(context.Background(), nil)
	if err != nil {
		log.Fatalf("error initializing app: %v\n", err)
	}

	return &Server{
		DB:          database.New(),
		FirebaseApp: app,
	}
}
