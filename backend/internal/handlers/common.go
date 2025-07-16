package handlers

import (
	"net/http"

	"github.com/labstack/echo/v4"
)

type CommonHandler struct{}

func NewCommonHandler() *CommonHandler {
	return &CommonHandler{}
}

func (h *CommonHandler) HealthCheck(c echo.Context) error {
	return c.JSON(http.StatusOK, map[string]string{
		"status":  "healthy",
		"message": "Server is running",
	})
}
