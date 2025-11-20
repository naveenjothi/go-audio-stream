package server

import (
	"fmt"
	"log"
	"net/http"
	"os"
	"strconv"
	"time"

	_ "github.com/joho/godotenv/autoload"

	"go-audio-stream/pkg/clients"
	"go-audio-stream/pkg/database"
)

type Server struct {
	port int

	db             database.Service
	identityClient *clients.IdentityClient
}

func NewServer() *http.Server {
	port, _ := strconv.Atoi(os.Getenv("API_GATEWAY_PORT"))
	identityServiceURL := os.Getenv("IDENTITY_SERVICE_URL")

	identityClient, err := clients.NewIdentityClient(identityServiceURL)
	if err != nil {
		log.Fatalf("Failed to create identity client: %v", err)
	}

	NewServer := &Server{
		port:           port,
		db:             database.New(),
		identityClient: identityClient,
	}

	// Declare Server config
	server := &http.Server{
		Addr:         fmt.Sprintf(":%d", NewServer.port),
		Handler:      NewServer.RegisterRoutes(),
		IdleTimeout:  time.Minute,
		ReadTimeout:  10 * time.Second,
		WriteTimeout: 30 * time.Second,
	}

	fmt.Printf("Server running on port %d\n", port)

	return server
}
