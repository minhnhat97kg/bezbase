package handlers

import (
	"net/http"

	"bezbase/internal/dto"
	"bezbase/internal/i18n"
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

// @Summary Register a new user
// @Tags Auth
// @Accept json
// @Produce json
// @Param request body dto.RegisterRequest true "User registration request"
// @Success 201 {object} dto.AuthResponse
// @Failure 400 {object} map[string]interface{}
// @Failure 409 {object} map[string]interface{}
// @Failure 500 {object} map[string]interface{}
// @Router /auth/register [post]
func (h *AuthHandler) Register(c echo.Context) error {
	t := i18n.NewTranslator(c.Request().Context())
	
	var req dto.RegisterRequest
	if err := c.Bind(&req); err != nil {
		return echo.NewHTTPError(http.StatusBadRequest, t.InvalidRequestBody())
	}

	response, err := h.authService.Register(req)
	if err != nil {
		switch err.Error() {
		case "username already registered":
			return echo.NewHTTPError(http.StatusConflict, t.Error("username_already_registered"))
		case "username already taken":
			return echo.NewHTTPError(http.StatusConflict, t.UsernameAlreadyTaken())
		case "email already registered":
			return echo.NewHTTPError(http.StatusConflict, t.Error("email_already_registered"))
		default:
			return echo.NewHTTPError(http.StatusInternalServerError, err.Error())
		}
	}

	return c.JSON(http.StatusCreated, response)
}

// @Summary Login with username and password
// @Tags Auth
// @Accept json
// @Produce json
// @Param request body dto.LoginRequest true "User login request"
// @Success 200 {object} dto.AuthResponse
// @Failure 400 {object} map[string]interface{}
// @Failure 401 {object} map[string]interface{}
// @Failure 500 {object} map[string]interface{}
// @Router /auth/login [post]
func (h *AuthHandler) Login(c echo.Context) error {
	t := i18n.NewTranslator(c.Request().Context())
	
	var req dto.LoginRequest
	if err := c.Bind(&req); err != nil {
		return echo.NewHTTPError(http.StatusBadRequest, t.InvalidRequestBody())
	}

	response, err := h.authService.LoginWithUsername(req)
	if err != nil {
		switch err.Error() {
		case "invalid credentials":
			return echo.NewHTTPError(http.StatusUnauthorized, t.InvalidCredentials())
		default:
			return echo.NewHTTPError(http.StatusInternalServerError, err.Error())
		}
	}

	return c.JSON(http.StatusOK, response)
}
