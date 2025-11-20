package middlewares

import (
	"go-audio-stream/pkg/clients"
	"go-audio-stream/pkg/models"
	"net/http"

	"github.com/labstack/echo/v4"
)

const UserContextKey string = "VerifiedUser"

func NewAuthMiddleware(client *clients.IdentityClient) echo.MiddlewareFunc {
	return func(next echo.HandlerFunc) echo.HandlerFunc {
		return func(c echo.Context) error {
			authHeader := c.Request().Header.Get("Authorization")
			if authHeader == "" {
				return echo.ErrUnauthorized
			}

			// Verify token using gRPC client
			// Extract token from "Bearer <token>"
			token := authHeader
			if len(authHeader) > 7 && authHeader[:7] == "Bearer " {
				token = authHeader[7:]
			}

			resp, err := client.VerifyToken(token)
			if err != nil {
				return echo.NewHTTPError(http.StatusUnauthorized, "Invalid token")
			}

			// Map response to User model
			user := models.User{
				BaseModel: models.BaseModel{
					ID:          resp.Id,
					IsSuspended: resp.IsSuspended,
				},
				Email:     resp.Email,
				FirstName: resp.Name, // Simplified mapping, as Name was combined
			}

			// Set the user in the context
			c.Set(UserContextKey, user)

			return next(c)
		}
	}
}
