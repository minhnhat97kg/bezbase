package handlers

import (
	"net/http"

	"bezbase/internal/i18n"

	"github.com/labstack/echo/v4"
)

type CommonHandler struct{}

func NewCommonHandler() *CommonHandler {
	return &CommonHandler{}
}

// @Summary Health check endpoint
// @Tags System
// @Produce json
// @Success 200 {object} map[string]string
// @Router /health [get]
func (h *CommonHandler) HealthCheck(c echo.Context) error {
	t := i18n.NewTranslator(c.Request().Context())
	
	return c.JSON(http.StatusOK, map[string]string{
		"status":  t.Status("healthy"),
		"message": t.Status("server_running"),
	})
}
