package server

import (
	"go-audio-stream/pkg/database"
	common_handlers "go-audio-stream/pkg/handlers"
	"go-audio-stream/pkg/middlewares"
	"go-audio-stream/services/catalog-service/internal/handlers"
	"net/http"

	"github.com/labstack/echo/v4"
	"github.com/labstack/echo/v4/middleware"
)

func (s *Server) RegisterRoutes() http.Handler {
	e := echo.New()
	e.Use(middleware.Logger())
	e.Use(middleware.Recover())

	e.Use(middleware.CORSWithConfig(middleware.CORSConfig{
		AllowOrigins:     []string{"https://*", "http://*"},
		AllowMethods:     []string{"GET", "POST", "PUT", "DELETE", "OPTIONS", "PATCH"},
		AllowHeaders:     []string{"Accept", "Authorization", "Content-Type", "X-CSRF-Token"},
		AllowCredentials: true,
		MaxAge:           300,
	}))

	e.Use(middlewares.CustomResponseMiddleware)

	e.GET("/health", s.withClient(common_handlers.HealthHandler))
	e.GET("/hello", s.withClient(common_handlers.HelloWorldHandler))
	protectedGroup := e.Group("/api/v1")
	protectedGroup.Use(middlewares.NewAuthMiddleware(s.identityClient))
	userEndpointGroup := protectedGroup.Group("/users")

	userEndpointGroup.POST("/", s.withClient(handlers.CreateUserHandler))
	userEndpointGroup.GET("/:id", s.withClient(handlers.FindOneUserById))
	userEndpointGroup.PUT("/:id", s.withClient(handlers.UpdateUserHandler))
	userEndpointGroup.DELETE("/:id", s.withClient(handlers.DeleteUserHandler))

	artistGroup := protectedGroup.Group("/artists")
	artistGroup.POST("/", s.withClient(handlers.CreateArtistHandler))
	artistGroup.GET("/", s.withClient(handlers.FindAllArtists))
	artistGroup.GET("/:id", s.withClient(handlers.FindOneArtistById))
	artistGroup.PUT("/:id", s.withClient(handlers.UpdateArtistHandler))
	artistGroup.DELETE("/:id", s.withClient(handlers.DeleteArtistHandler))

	songGroup := protectedGroup.Group("/songs")
	songGroup.POST("/", s.withClient(handlers.CreateSongHandler))
	songGroup.GET("/", s.withClient(handlers.FindAllSongs))
	songGroup.GET("/:id", s.withClient(handlers.FindOneSongById))
	songGroup.PUT("/:id", s.withClient(handlers.UpdateSongHandler))
	songGroup.DELETE("/:id", s.withClient(handlers.DeleteSongHandler))

	playlistGroup := protectedGroup.Group("/playlists")
	playlistGroup.POST("/", s.withClient(handlers.CreatePlaylistHandler))
	playlistGroup.GET("/", s.withClient(handlers.FindAllPlaylists))
	playlistGroup.GET("/:id", s.withClient(handlers.FindOnePlaylistById))
	playlistGroup.PUT("/:id", s.withClient(handlers.UpdatePlaylistHandler))
	playlistGroup.DELETE("/:id", s.withClient(handlers.DeletePlaylistHandler))
	playlistGroup.POST("/:id/songs", s.withClient(handlers.AddSongToPlaylistHandler))
	playlistGroup.DELETE("/:id/songs/:song_id", s.withClient(handlers.RemoveSongFromPlaylistHandler))

	return e
}

func (s *Server) withClient(handler func(echo.Context, database.Service) error) echo.HandlerFunc {
	return func(c echo.Context) error {
		return handler(c, s.db)
	}
}
