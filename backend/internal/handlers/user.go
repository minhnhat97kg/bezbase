package handlers

import (
	"net/http"

	"bezbase/internal/dto"
	"bezbase/internal/pkg/auth"
	"bezbase/internal/services"

	"github.com/labstack/echo/v4"
)

type UserHandler struct {
	userService *services.UserService
}

func NewUserHandler(userService *services.UserService) *UserHandler {
	return &UserHandler{
		userService: userService,
	}
}

func (h *UserHandler) GetProfile(c echo.Context) error {
	claims := c.Get("user").(*auth.Claims)

	user, err := h.userService.GetProfile(claims.UserID)
	if err != nil {
		return echo.NewHTTPError(http.StatusNotFound, err.Error())
	}

	return c.JSON(http.StatusOK, user)
}

func (h *UserHandler) UpdateProfile(c echo.Context) error {
	claims := c.Get("user").(*auth.Claims)

	var req dto.UpdateProfileRequest
	if err := c.Bind(&req); err != nil {
		return echo.NewHTTPError(http.StatusBadRequest, "Invalid request body")
	}

	user, err := h.userService.UpdateProfile(claims.UserID, req)
	if err != nil {
		return echo.NewHTTPError(http.StatusInternalServerError, err.Error())
	}

	return c.JSON(http.StatusOK, user)
}

func (h *UserHandler) GetUsers(c echo.Context) error {
	// Get search query parameter
	search := c.QueryParam("search")

	var users []dto.UserResponse
	var err error

	if search != "" {
		users, err = h.userService.SearchUsers(search)
	} else {
		users, err = h.userService.GetAllUsers()
	}

	if err != nil {
		return echo.NewHTTPError(http.StatusInternalServerError, err.Error())
	}

	return c.JSON(http.StatusOK, users)
}
