package handlers

import (
	"net/http"

	"bezbase/internal/dto"
	"bezbase/internal/services"

	"github.com/labstack/echo/v4"
)

type AuthHandler struct {
	authService *services.AuthService
}

func NewAuthHandler(authService *services.AuthService) *AuthHandler {
	return &AuthHandler{
		authService: authService,
	}
}

func (h *AuthHandler) Register(c echo.Context) error {
	var req dto.RegisterRequest
	if err := c.Bind(&req); err != nil {
		return echo.NewHTTPError(http.StatusBadRequest, "Invalid request body")
	}

	response, err := h.authService.RegisterWithEmail(req)
	if err != nil {
		switch err.Error() {
		case "username already registered":
			return echo.NewHTTPError(http.StatusConflict, err.Error())
		default:
			return echo.NewHTTPError(http.StatusInternalServerError, err.Error())
		}
	}

	return c.JSON(http.StatusCreated, response)
}

func (h *AuthHandler) Login(c echo.Context) error {
	var req dto.LoginRequest
	if err := c.Bind(&req); err != nil {
		return echo.NewHTTPError(http.StatusBadRequest, "Invalid request body")
	}

	response, err := h.authService.LoginWithUsername(req)
	if err != nil {
		switch err.Error() {
		case "invalid credentials":
			return echo.NewHTTPError(http.StatusUnauthorized, err.Error())
		default:
			return echo.NewHTTPError(http.StatusInternalServerError, err.Error())
		}
	}

	return c.JSON(http.StatusOK, response)
}
