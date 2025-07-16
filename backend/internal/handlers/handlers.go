package handlers

import (
	"net/http"

	"bezbase/internal/auth"
	"bezbase/internal/models"
	"bezbase/internal/services"

	"github.com/casbin/casbin/v2"

	"github.com/labstack/echo/v4"
)

type Handler struct {
	authService *services.AuthService
	userService *services.UserService
	enforcer    *casbin.Enforcer
}

func NewHandler(authService *services.AuthService, userService *services.UserService) *Handler {
	return &Handler{
		authService: authService,
		userService: userService,
		enforcer:    nil,
	}
}

func NewHandlerWithEnforcer(authService *services.AuthService, userService *services.UserService, enforcer *casbin.Enforcer) *Handler {
	return &Handler{
		authService: authService,
		userService: userService,
		enforcer:    enforcer,
	}
}

func (h *Handler) Register(c echo.Context) error {
	var req models.RegisterRequest
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

func (h *Handler) Login(c echo.Context) error {
	var req models.LoginRequest
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

func (h *Handler) GetProfile(c echo.Context) error {
	claims := c.Get("user").(*auth.Claims)

	user, err := h.userService.GetProfile(claims.UserID)
	if err != nil {
		return echo.NewHTTPError(http.StatusNotFound, err.Error())
	}

	return c.JSON(http.StatusOK, user)
}

func (h *Handler) UpdateProfile(c echo.Context) error {
	claims := c.Get("user").(*auth.Claims)

	var req models.UpdateProfileRequest
	if err := c.Bind(&req); err != nil {
		return echo.NewHTTPError(http.StatusBadRequest, "Invalid request body")
	}

	user, err := h.userService.UpdateProfile(claims.UserID, req)
	if err != nil {
		return echo.NewHTTPError(http.StatusInternalServerError, err.Error())
	}

	return c.JSON(http.StatusOK, user)
}

func (h *Handler) HealthCheck(c echo.Context) error {
	return c.JSON(http.StatusOK, map[string]string{
		"status":  "healthy",
		"message": "Server is running",
	})
}
