```go
package handler

import (
	"net/http"

	"github.com/labstack/echo/v4"
)

// RBACHandler handles RBAC related endpoints
type RBACHandler struct{}

// NewRBACHandler creates a new RBACHandler
func NewRBACHandler() *RBACHandler {
	return &RBACHandler{}
}

// GetPermissionResources godoc
// @Summary Get available resources for permissions
// @Tags RBAC
// @Produce json
// @Success 200 {object} map[string]interface{}
// @Router /api/v1/rbac/resources [get]
func (h *RBACHandler) GetPermissionResources(c echo.Context) error {
	resources := []string{
		"users", "posts", "comments", "profile", "settings", "dashboard",
		"files", "reports", "notifications", "audit",
	}
	return c.JSON(http.StatusOK, map[string]interface{}{
		"resources": resources,
	})
}
```