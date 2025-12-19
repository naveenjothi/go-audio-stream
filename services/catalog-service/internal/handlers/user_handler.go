package handlers

import (
	"go-audio-stream/pkg/database"
	"go-audio-stream/pkg/models"
	"net/http"
	"strings"

	"github.com/labstack/echo/v4"
	"github.com/labstack/gommon/random"
)

func calculateUserName(email string) string {
	emailStr := strings.Split(email, "@")[0]
	randomNumber := random.String(4, random.Numeric)
	return emailStr + randomNumber
}

// CreateUserHandler creates a new user.
// @Summary      Create a new user
// @Description  Create a new user with the provided details
// @Tags         users
// @Accept       json
// @Produce      json
// @Param        user  body      models.User  true  "User Data"
// @Success      201   {object}  models.User
// @Failure      400   {object}  map[string]string
// @Failure      500   {object}  map[string]string
// @Router       /api/v1/users [post]
func CreateUserHandler(c echo.Context, db database.Service) error {
	user := new(models.User)

	if err := c.Bind(user); err != nil {
		return c.JSON(http.StatusBadRequest, echo.Map{"error": err.Error()})
	}

	// Validate required fields
	if user.Email == "" {
		return c.JSON(http.StatusBadRequest, echo.Map{"error": "email is required"})
	}
	if user.FirstName == "" {
		return c.JSON(http.StatusBadRequest, echo.Map{"error": "first_name is required"})
	}

	user.Username = calculateUserName(user.Email)
	_, err := db.Create(user)
	if err != nil {
		return c.JSON(http.StatusInternalServerError, echo.Map{"error": err.Error()})
	}

	return c.JSON(http.StatusCreated, user)
}

// FindOneUserById retrieves a user by ID.
// @Summary      Get a user
// @Description  Get a user by their Firebase ID
// @Tags         users
// @Accept       json
// @Produce      json
// @Param        id   path      string  true  "User ID"
// @Success      200  {object}  models.User
// @Failure      500  {object}  map[string]string
// @Router       /api/v1/users/users/{id} [get]
func FindOneUserById(c echo.Context, db database.Service) error {
	id := c.Param("id")

	var user models.User

	_, err := db.Find(&user, "firebase_id = ?", id)

	if err != nil {
		return c.JSON(http.StatusInternalServerError, echo.Map{"error": err.Error()})
	}

	return c.JSON(http.StatusOK, user)
}

// UpdateUserHandler updates an existing user.
// @Summary      Update a user
// @Description  Update a user's details
// @Tags         users
// @Accept       json
// @Produce      json
// @Param        id    path      string       true  "User ID"
// @Param        user  body      models.User  true  "User Data"
// @Success      200   {object}  models.User
// @Failure      400   {object}  map[string]string
// @Failure      500   {object}  map[string]string
// @Router       /api/v1/users/users/{id} [put]
func UpdateUserHandler(c echo.Context, db database.Service) error {
	id := c.Param("id")
	user := new(models.User)

	_, err := db.Find(&user, "id = ?", id)

	if err != nil {
		return c.JSON(http.StatusInternalServerError, echo.Map{"error": err.Error()})
	}

	if err := c.Bind(user); err != nil {
		return c.JSON(http.StatusBadRequest, echo.Map{"error": err.Error()})
	}

	_, err = db.Update(&models.User{}, user, "id = ?", id)

	if err != nil {
		return c.JSON(http.StatusInternalServerError, echo.Map{"error": err.Error()})
	}

	return c.JSON(http.StatusOK, user)
}

// DeleteUserHandler deletes a user.
// @Summary      Delete a user
// @Description  Delete a user by ID
// @Tags         users
// @Accept       json
// @Produce      json
// @Param        id   path      string  true  "User ID"
// @Success      200  {object}  map[string]string
// @Failure      500  {object}  map[string]string
// @Router       /api/v1/users/users/{id} [delete]
func DeleteUserHandler(c echo.Context, db database.Service) error {
	id := c.Param("id")

	_, err := db.Delete(&models.User{}, "id = ?", id)

	if err != nil {
		return c.JSON(http.StatusInternalServerError, echo.Map{"error": err.Error()})
	}

	return c.JSON(http.StatusOK, echo.Map{"message": "User deleted successfully"})
}
