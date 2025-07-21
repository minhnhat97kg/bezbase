package handlers

import (
	"fmt"
	"net/http"

	"bezbase/internal/dto"
	"bezbase/internal/i18n"
	"bezbase/internal/pkg/auth"
	"bezbase/internal/pkg/contextx"
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

// @Summary Get current user permissions
// @Tags User
// @Security BearerAuth
// @Produce json
// @Success 200 {object} map[string]interface{}
// @Failure 401 {object} map[string]interface{}
// @Failure 500 {object} map[string]interface{}
// @Router /v1/rbac/me/permissions [get]
// GetCurrentUserPermissions returns all permissions for the current user
func (h *UserHandler) GetCurrentUserPermissions(c echo.Context) error {
	t := i18n.NewTranslator(c.Request().Context())

	claims, ok := c.Get("user").(*auth.Claims)
	if !ok || claims == nil {
		return echo.NewHTTPError(http.StatusUnauthorized, t.Error("invalid_user_context"))
	}
	if h.rbacService == nil {
		return echo.NewHTTPError(http.StatusInternalServerError, t.Error("rbac_service_not_available"))
	}
	permissions, err := h.rbacService.GetPermissionsForUser(contextx.NewWithRequestContext(c), claims.UserID)
	if err != nil {
		return echo.NewHTTPError(http.StatusInternalServerError, err.Error())
	}
	return c.JSON(http.StatusOK, map[string]any{
		"user_id":     claims.UserID,
		"permissions": permissions,
	})
}

// @Summary Get current user profile
// @Tags User
// @Security BearerAuth
// @Produce json
// @Success 200 {object} dto.UserResponse
// @Failure 401 {object} map[string]interface{}
// @Failure 404 {object} map[string]interface{}
// @Router /v1/profile [get]
func (h *UserHandler) GetProfile(c echo.Context) error {
	claims := c.Get("user").(*auth.Claims)
	ctx := contextx.NewWithRequestContext(c)
	user, err := h.userService.GetProfile(ctx, claims.UserID)
	if err != nil {
		return echo.NewHTTPError(http.StatusNotFound, err.Error())
	}
	return c.JSON(http.StatusOK, user)
}

// @Summary Update current user profile
// @Tags User
// @Security BearerAuth
// @Accept json
// @Produce json
// @Param request body dto.UpdateProfileRequest true "Profile update request"
// @Success 200 {object} dto.UserResponse
// @Failure 400 {object} map[string]interface{}
// @Failure 401 {object} map[string]interface{}
// @Failure 500 {object} map[string]interface{}
// @Router /v1/profile [put]
func (h *UserHandler) UpdateProfile(c echo.Context) error {
	claims := c.Get("user").(*auth.Claims)
	ctx := contextx.NewWithRequestContext(c)
	var req dto.UpdateProfileRequest
	if err := c.Bind(&req); err != nil {
		return echo.NewHTTPError(http.StatusBadRequest, "Invalid request body")
	}
	user, err := h.userService.UpdateProfile(ctx, claims.UserID, req)
	if err != nil {
		return echo.NewHTTPError(http.StatusInternalServerError, err.Error())
	}
	return c.JSON(http.StatusOK, user)
}

// @Summary Change current user password
// @Tags User
// @Security BearerAuth
// @Accept json
// @Produce json
// @Param request body dto.ChangePasswordRequest true "Password change request"
// @Success 200 {object} map[string]interface{}
// @Failure 400 {object} map[string]interface{}
// @Failure 401 {object} map[string]interface{}
// @Failure 500 {object} map[string]interface{}
// @Router /v1/profile/password [put]
func (h *UserHandler) ChangePassword(c echo.Context) error {
	claims := c.Get("user").(*auth.Claims)

	var req dto.ChangePasswordRequest
	if err := c.Bind(&req); err != nil {
		return echo.NewHTTPError(http.StatusBadRequest, "Invalid request body")
	}

	// Validate password confirmation
	if req.NewPassword != req.ConfirmPassword {
		return echo.NewHTTPError(http.StatusBadRequest, "New password and confirm password do not match")
	}

	err := h.userService.ChangePassword(contextx.NewWithRequestContext(c), claims.UserID, req.CurrentPassword, req.NewPassword)
	if err != nil {
		if err.Error() == "invalid current password" {
			return echo.NewHTTPError(http.StatusBadRequest, "Current password is incorrect")
		}
		return echo.NewHTTPError(http.StatusInternalServerError, err.Error())
	}

	return c.JSON(http.StatusOK, map[string]interface{}{
		"message": "Password changed successfully",
	})
}

// @Summary Get all users (admin only)
// @Tags User
// @Security BearerAuth
// @Produce json
// @Param search query string false "Search users by name or email"
// @Success 200 {array} dto.UserResponse
// @Failure 401 {object} map[string]interface{}
// @Failure 403 {object} map[string]interface{}
// @Failure 500 {object} map[string]interface{}
// @Router /v1/users [get]
func (h *UserHandler) GetUsers(c echo.Context) error {
	// Get search query parameter
	search := c.QueryParam("search")

	var users []dto.UserResponse
	var err error
	ctx := contextx.NewWithRequestContext(c)
	if search != "" {
		users, err = h.userService.SearchUsers(ctx, search)
	} else {
		users, err = h.userService.GetAllUsers(ctx)
	}

	if err != nil {
		return echo.NewHTTPError(http.StatusInternalServerError, err.Error())
	}

	return c.JSON(http.StatusOK, users)
}

// @Summary Get user by ID (admin only)
// @Tags User
// @Security BearerAuth
// @Produce json
// @Param id path string true "User ID"
// @Success 200 {object} dto.UserResponse
// @Failure 400 {object} map[string]interface{}
// @Failure 401 {object} map[string]interface{}
// @Failure 403 {object} map[string]interface{}
// @Failure 404 {object} map[string]interface{}
// @Router /v1/users/{id} [get]
func (h *UserHandler) GetUser(c echo.Context) error {
	userID := c.Param("id")

	// Convert userID to uint
	var id uint
	if _, err := fmt.Sscanf(userID, "%d", &id); err != nil {
		return echo.NewHTTPError(http.StatusBadRequest, "Invalid user ID")
	}

	user, err := h.userService.GetUserByIDDetailed(contextx.NewWithRequestContext(c), id)
	if err != nil {
		return echo.NewHTTPError(http.StatusNotFound, err.Error())
	}

	return c.JSON(http.StatusOK, user)
}

// @Summary Create a new user (admin only)
// @Tags User
// @Security BearerAuth
// @Accept json
// @Produce json
// @Param request body dto.CreateUserRequest true "User creation request"
// @Success 201 {object} dto.UserResponse
// @Failure 400 {object} map[string]interface{}
// @Failure 401 {object} map[string]interface{}
// @Failure 403 {object} map[string]interface{}
// @Failure 500 {object} map[string]interface{}
// @Router /v1/users [post]
func (h *UserHandler) CreateUser(c echo.Context) error {
	var req dto.CreateUserRequest
	if err := c.Bind(&req); err != nil {
		return echo.NewHTTPError(http.StatusBadRequest, "Invalid request body")
	}

	// Validate required fields
	if req.FirstName == "" || req.LastName == "" || req.Email == "" || req.Password == "" {
		return echo.NewHTTPError(http.StatusBadRequest, "Missing required fields")
	}

	user, err := h.userService.CreateUser(contextx.NewWithRequestContext(c), req)
	if err != nil {
		return echo.NewHTTPError(http.StatusInternalServerError, err.Error())
	}

	return c.JSON(http.StatusCreated, user)
}

// @Summary Update user by ID (admin only)
// @Tags User
// @Security BearerAuth
// @Accept json
// @Produce json
// @Param id path string true "User ID"
// @Param request body dto.UpdateUserRequest true "User update request"
// @Success 200 {object} dto.UserResponse
// @Failure 400 {object} map[string]interface{}
// @Failure 401 {object} map[string]interface{}
// @Failure 403 {object} map[string]interface{}
// @Failure 500 {object} map[string]interface{}
// @Router /v1/users/{id} [put]
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

	user, err := h.userService.UpdateUser(contextx.NewWithRequestContext(c), id, req)
	if err != nil {
		return echo.NewHTTPError(http.StatusInternalServerError, err.Error())
	}

	return c.JSON(http.StatusOK, user)
}

// @Summary Delete user by ID (admin only)
// @Tags User
// @Security BearerAuth
// @Param id path string true "User ID"
// @Success 200 {object} map[string]interface{}
// @Failure 400 {object} map[string]interface{}
// @Failure 401 {object} map[string]interface{}
// @Failure 403 {object} map[string]interface{}
// @Failure 500 {object} map[string]interface{}
// @Router /v1/users/{id} [delete]
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

	err := h.userService.DeleteUser(contextx.NewWithRequestContext(c), id)
	if err != nil {
		return echo.NewHTTPError(http.StatusInternalServerError, err.Error())
	}

	return c.JSON(http.StatusOK, map[string]interface{}{
		"message": "User deleted successfully",
	})
}
