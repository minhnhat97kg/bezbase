package handlers

import (
	"fmt"
	"net/http"

	"bezbase/internal/dto"
	"bezbase/internal/pkg/auth"
	"bezbase/internal/services"

	"github.com/labstack/echo/v4"
)

type UserHandler struct {
	userService *services.UserService
	rbacService *services.RBACService
}

func NewUserHandler(userService *services.UserService, rbacService *services.RBACService) *UserHandler {
	return &UserHandler{
		userService: userService,
		rbacService: rbacService,
	}
}

// GetCurrentUserPermissions returns all permissions for the current user
func (h *UserHandler) GetCurrentUserPermissions(c echo.Context) error {
	claims, ok := c.Get("user").(*auth.Claims)
	if !ok || claims == nil {
		return echo.NewHTTPError(http.StatusUnauthorized, "Invalid user context")
	}
	if h.rbacService == nil {
		return echo.NewHTTPError(http.StatusInternalServerError, "RBAC service not available")
	}
	permissions, err := h.rbacService.GetPermissionsForUser(claims.UserID)
	if err != nil {
		return echo.NewHTTPError(http.StatusInternalServerError, err.Error())
	}
	return c.JSON(http.StatusOK, map[string]any{
		"user_id":     claims.UserID,
		"permissions": permissions,
	})
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

func (h *UserHandler) GetUser(c echo.Context) error {
	userID := c.Param("id")

	// Convert userID to uint
	var id uint
	if _, err := fmt.Sscanf(userID, "%d", &id); err != nil {
		return echo.NewHTTPError(http.StatusBadRequest, "Invalid user ID")
	}

	user, err := h.userService.GetUserByIDDetailed(id)
	if err != nil {
		return echo.NewHTTPError(http.StatusNotFound, err.Error())
	}

	return c.JSON(http.StatusOK, user)
}

func (h *UserHandler) CreateUser(c echo.Context) error {
	var req dto.CreateUserRequest
	if err := c.Bind(&req); err != nil {
		return echo.NewHTTPError(http.StatusBadRequest, "Invalid request body")
	}

	// Validate required fields
	if req.FirstName == "" || req.LastName == "" || req.Email == "" || req.Password == "" {
		return echo.NewHTTPError(http.StatusBadRequest, "Missing required fields")
	}

	user, err := h.userService.CreateUser(req)
	if err != nil {
		return echo.NewHTTPError(http.StatusInternalServerError, err.Error())
	}

	return c.JSON(http.StatusCreated, user)
}

func (h *UserHandler) UpdateUser(c echo.Context) error {
	userID := c.Param("id")

	// Convert userID to uint
	var id uint
	if _, err := fmt.Sscanf(userID, "%d", &id); err != nil {
		return echo.NewHTTPError(http.StatusBadRequest, "Invalid user ID")
	}

	var req dto.UpdateUserRequest
	if err := c.Bind(&req); err != nil {
		return echo.NewHTTPError(http.StatusBadRequest, "Invalid request body")
	}

	user, err := h.userService.UpdateUser(id, req)
	if err != nil {
		return echo.NewHTTPError(http.StatusInternalServerError, err.Error())
	}

	return c.JSON(http.StatusOK, user)
}

func (h *UserHandler) DeleteUser(c echo.Context) error {
	userID := c.Param("id")

	// Convert userID to uint
	var id uint
	if _, err := fmt.Sscanf(userID, "%d", &id); err != nil {
		return echo.NewHTTPError(http.StatusBadRequest, "Invalid user ID")
	}

	// Get current user from JWT token
	claims := c.Get("user").(*auth.Claims)
	currentUserID := claims.UserID

	// Prevent self-deletion
	if id == currentUserID {
		return echo.NewHTTPError(http.StatusForbidden, "Cannot delete your own account")
	}

	err := h.userService.DeleteUser(id)
	if err != nil {
		return echo.NewHTTPError(http.StatusInternalServerError, err.Error())
	}

	return c.JSON(http.StatusOK, map[string]interface{}{
		"message": "User deleted successfully",
	})
}
