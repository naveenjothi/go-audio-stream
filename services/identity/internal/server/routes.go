package server

import (
	"go-audio-stream/pkg/database"
	common_handlers "go-audio-stream/pkg/handlers"
	"go-audio-stream/pkg/middlewares"
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

	return e
}

func (s *Server) withClient(handler func(echo.Context, database.Service) error) echo.HandlerFunc {
	return func(c echo.Context) error {
		return handler(c, s.db)
	}
}
