package server

import (
	"context"
	"fmt"
	"log"
	"net/http"
	"os"
	"strconv"
	"time"

	firebase "firebase.google.com/go/v4"
	_ "github.com/joho/godotenv/autoload"

	"go-audio-stream/pkg/database"
)

type Server struct {
	port int

	db database.Service

	firebase_app *firebase.App
}

func NewServer() *http.Server {
	port, _ := strconv.Atoi(os.Getenv("IDENTITY_PORT"))

	app, err := firebase.NewApp(context.Background(), nil)
	if err != nil {
		log.Fatalf("error initializing app: %v\n", err)
	}

	NewServer := &Server{
		port:         port,
		db:           database.New(),
		firebase_app: app,
	}

	// Declare Server config
	server := &http.Server{
		Addr:         fmt.Sprintf(":%d", NewServer.port),
		Handler:      NewServer.RegisterRoutes(),
		IdleTimeout:  time.Minute,
		ReadTimeout:  10 * time.Second,
		WriteTimeout: 30 * time.Second,
	}

	return server
}
