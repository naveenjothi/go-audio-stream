package handlers

import (
	"go-audio-stream/pkg/database"
	"net/http"

	"github.com/labstack/echo/v4"
)

func HealthHandler(c echo.Context, db database.Service) error {
	return c.JSON(http.StatusOK, db.Health())
}

func HelloWorldHandler(c echo.Context, db database.Service) error {
	resp := map[string]string{"message": "Hello World"}
	return c.JSON(http.StatusOK, resp)
}
