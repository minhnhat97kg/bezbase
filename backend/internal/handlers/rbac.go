package handlers

import (
	"net/http"
	"strconv"

	"bezbase/internal/dto"
	"bezbase/internal/services"

	"github.com/labstack/echo/v4"
)

type RBACHandler struct {
	rbacService *services.RBACService
}

func NewRBACHandler(rbacService *services.RBACService) *RBACHandler {
	return &RBACHandler{
		rbacService: rbacService,
	}
}

// @Summary Create a new role
// @Tags RBAC
// @Accept json
// @Produce json
// @Param request body dto.CreateRoleRequest true "Role creation request"
// @Success 200 {object} dto.RoleResponse
// @Failure 400 {object} map[string]interface{}
// @Failure 500 {object} map[string]interface{}
// @Router /api/v1/rbac/roles [post]
func (h *RBACHandler) CreateRole(c echo.Context) error {
	var req dto.CreateRoleRequest
	if err := c.Bind(&req); err != nil {
		return echo.NewHTTPError(http.StatusBadRequest, "Invalid request body")
	}

	role, err := h.rbacService.CreateRole(req)
	if err != nil {
		return echo.NewHTTPError(http.StatusInternalServerError, err.Error())
	}

	return c.JSON(http.StatusOK, dto.ToRoleResponse(role))
}

// @Summary Get all roles
// @Tags RBAC
// @Produce json
// @Success 200 {object} []dto.RoleResponse
// @Failure 500 {object} map[string]interface{}
// @Router /api/v1/rbac/roles [get]
func (h *RBACHandler) GetRoles(c echo.Context) error {
	roles, err := h.rbacService.GetAllRolesWithMetadata()
	if err != nil {
		return echo.NewHTTPError(http.StatusInternalServerError, err.Error())
	}

	return c.JSON(http.StatusOK, dto.ToRoleResponses(roles))
}

// @Summary Get role by ID
// @Tags RBAC
// @Param role_id path string true "Role ID"
// @Success 200 {object} dto.RoleResponse
// @Failure 400 {object} map[string]interface{}
// @Failure 404 {object} map[string]interface{}
// @Failure 500 {object} map[string]interface{}
// @Router /api/v1/rbac/roles/{role_id} [get]
func (h *RBACHandler) GetRole(c echo.Context) error {
	roleIDStr := c.Param("role_id")
	roleID, err := strconv.ParseUint(roleIDStr, 10, 32)
	if err != nil {
		return echo.NewHTTPError(http.StatusBadRequest, "Invalid role ID")
	}

	role, err := h.rbacService.GetRoleByID(uint(roleID))
	if err != nil {
		return echo.NewHTTPError(http.StatusNotFound, err.Error())
	}

	return c.JSON(http.StatusOK, dto.ToRoleResponse(role))
}

// @Summary Update a role
// @Tags RBAC
// @Accept json
// @Produce json
// @Param role_id path string true "Role ID"
// @Param request body dto.UpdateRoleRequest true "Role update request"
// @Success 200 {object} dto.RoleResponse
// @Failure 400 {object} map[string]interface{}
// @Failure 404 {object} map[string]interface{}
// @Failure 500 {object} map[string]interface{}
// @Router /api/v1/rbac/roles/{role_id} [put]
func (h *RBACHandler) UpdateRole(c echo.Context) error {
	roleIDStr := c.Param("role_id")
	roleID, err := strconv.ParseUint(roleIDStr, 10, 32)
	if err != nil {
		return echo.NewHTTPError(http.StatusBadRequest, "Invalid role ID")
	}

	var req dto.UpdateRoleRequest
	if err := c.Bind(&req); err != nil {
		return echo.NewHTTPError(http.StatusBadRequest, "Invalid request body")
	}

	role, err := h.rbacService.UpdateRole(uint(roleID), req)
	if err != nil {
		return echo.NewHTTPError(http.StatusInternalServerError, err.Error())
	}

	return c.JSON(http.StatusOK, dto.ToRoleResponse(role))
}

// @Summary Delete a role
// @Tags RBAC
// @Param role path string true "Role name"
// @Success 200 {object} map[string]interface{}
// @Failure 400 {object} map[string]interface{}
// @Failure 500 {object} map[string]interface{}
// @Router /api/v1/rbac/roles/{role} [delete]
func (h *RBACHandler) DeleteRole(c echo.Context) error {
	role := c.Param("role")
	if role == "" {
		return echo.NewHTTPError(http.StatusBadRequest, "Role parameter is required")
	}

	if err := h.rbacService.DeleteRole(role); err != nil {
		return echo.NewHTTPError(http.StatusInternalServerError, err.Error())
	}

	return c.JSON(http.StatusOK, map[string]any{
		"message": "Role deleted successfully",
		"role":    role,
	})
}

// @Summary Assign role to user
// @Tags RBAC
// @Accept json
// @Produce json
// @Param request body AssignRoleRequest true "Role assignment request"
// @Success 200 {object} map[string]interface{}
// @Failure 400 {object} map[string]interface{}
// @Failure 500 {object} map[string]interface{}
// @Router /api/v1/rbac/users/assign-role [post]
func (h *RBACHandler) AssignRole(c echo.Context) error {
	var req dto.AssignRoleRequest
	if err := c.Bind(&req); err != nil {
		return echo.NewHTTPError(http.StatusBadRequest, "Invalid request body")
	}

	if err := h.rbacService.AssignRoleToUser(req.UserID, req.Role); err != nil {
		return echo.NewHTTPError(http.StatusInternalServerError, err.Error())
	}

	return c.JSON(http.StatusOK, map[string]any{
		"message": "Role assigned successfully",
		"user_id": req.UserID,
		"role":    req.Role,
	})
}

// @Summary Remove role from user
// @Tags RBAC
// @Accept json
// @Produce json
// @Param request body AssignRoleRequest true "Role removal request"
// @Success 200 {object} map[string]interface{}
// @Failure 400 {object} map[string]interface{}
// @Failure 500 {object} map[string]interface{}
// @Router /api/v1/rbac/users/remove-role [post]
func (h *RBACHandler) RemoveRole(c echo.Context) error {
	var req dto.AssignRoleRequest
	if err := c.Bind(&req); err != nil {
		return echo.NewHTTPError(http.StatusBadRequest, "Invalid request body")
	}

	if err := h.rbacService.RemoveRoleFromUser(req.UserID, req.Role); err != nil {
		return echo.NewHTTPError(http.StatusInternalServerError, err.Error())
	}

	return c.JSON(http.StatusOK, map[string]any{
		"message": "Role removed successfully",
		"user_id": req.UserID,
		"role":    req.Role,
	})
}

// @Summary Get user roles
// @Tags RBAC
// @Param user_id path string true "User ID"
// @Success 200 {object} map[string]interface{}
// @Failure 400 {object} map[string]interface{}
// @Failure 500 {object} map[string]interface{}
// @Router /api/v1/rbac/users/{user_id}/roles [get]
func (h *RBACHandler) GetUserRoles(c echo.Context) error {
	userIDStr := c.Param("user_id")
	userID, err := strconv.ParseUint(userIDStr, 10, 32)
	if err != nil {
		return echo.NewHTTPError(http.StatusBadRequest, "Invalid user ID")
	}

	roles, err := h.rbacService.GetUserRoles(uint(userID))
	if err != nil {
		return echo.NewHTTPError(http.StatusInternalServerError, err.Error())
	}

	return c.JSON(http.StatusOK, map[string]any{
		"user_id": userID,
		"roles":   roles,
	})
}

// @Summary Get users with specific role
// @Tags RBAC
// @Param role path string true "Role name"
// @Success 200 {object} map[string]interface{}
// @Failure 400 {object} map[string]interface{}
// @Failure 500 {object} map[string]interface{}
// @Router /api/v1/rbac/roles/{role}/users [get]
func (h *RBACHandler) GetUsersWithRole(c echo.Context) error {
	role := c.Param("role")
	if role == "" {
		return echo.NewHTTPError(http.StatusBadRequest, "Role parameter is required")
	}

	userIDs, err := h.rbacService.GetUsersWithRole(role)
	if err != nil {
		return echo.NewHTTPError(http.StatusInternalServerError, err.Error())
	}

	return c.JSON(http.StatusOK, map[string]any{
		"role":     role,
		"user_ids": userIDs,
	})
}

// @Summary List all permissions with pagination
// @Tags RBAC
// @Produce json
// @Param page query int false "Page number (default: 1)"
// @Param page_size query int false "Page size (default: 10, max: 100)"
// @Param role query string false "Filter by role"
// @Param resource query string false "Filter by resource"
// @Param action query string false "Filter by action"
// @Param sort query string false "Sort field (role, resource, action)"
// @Param order query string false "Sort order (asc, desc)"
// @Success 200 {object} PermissionsListResponse
// @Failure 400 {object} map[string]interface{}
// @Failure 500 {object} map[string]interface{}
// @Router /api/v1/rbac/permissions [get]
func (h *RBACHandler) GetPermissions(c echo.Context) error {
	// Parse pagination parameters
	page, err := strconv.Atoi(c.QueryParam("page"))
	if err != nil || page < 1 {
		page = 1
	}

	pageSize, err := strconv.Atoi(c.QueryParam("page_size"))
	if err != nil || pageSize < 1 {
		pageSize = 10
	}
	if pageSize > 100 {
		pageSize = 100
	}

	// Parse filter parameters
	roleFilter := c.QueryParam("role")
	resourceFilter := c.QueryParam("resource")
	actionFilter := c.QueryParam("action")

	// Parse sort parameters
	sortField := c.QueryParam("sort")
	sortOrder := c.QueryParam("order")
	if sortOrder != "desc" {
		sortOrder = "asc"
	}

	permissions, total, err := h.rbacService.GetAllPermissions(page, pageSize, roleFilter, resourceFilter, actionFilter, sortField, sortOrder)
	if err != nil {
		return echo.NewHTTPError(http.StatusInternalServerError, err.Error())
	}

	totalPages := (total + pageSize - 1) / pageSize

	response := dto.PermissionsListResponse{
		Permissions: permissions,
		Total:       total,
		Page:        page,
		PageSize:    pageSize,
		TotalPages:  totalPages,
	}

	return c.JSON(http.StatusOK, response)
}

// @Summary Add permission to role
// @Tags RBAC
// @Accept json
// @Produce json
// @Param request body PermissionRequest true "Permission request"
// @Success 200 {object} map[string]interface{}
// @Failure 400 {object} map[string]interface{}
// @Failure 500 {object} map[string]interface{}
// @Router /api/v1/rbac/permissions [post]
func (h *RBACHandler) AddPermission(c echo.Context) error {
	var req dto.PermissionRequest
	if err := c.Bind(&req); err != nil {
		return echo.NewHTTPError(http.StatusBadRequest, "Invalid request body")
	}

	if err := h.rbacService.AddPermission(req.Role, req.Resource, req.Action); err != nil {
		return echo.NewHTTPError(http.StatusInternalServerError, err.Error())
	}

	return c.JSON(http.StatusOK, map[string]any{
		"message":  "Permission added successfully",
		"role":     req.Role,
		"resource": req.Resource,
		"action":   req.Action,
	})
}

// @Summary Remove permission from role
// @Tags RBAC
// @Accept json
// @Produce json
// @Param request body PermissionRequest true "Permission request"
// @Success 200 {object} map[string]interface{}
// @Failure 400 {object} map[string]interface{}
// @Failure 500 {object} map[string]interface{}
// @Router /api/v1/rbac/permissions [delete]
func (h *RBACHandler) RemovePermission(c echo.Context) error {
	var req dto.PermissionRequest
	if err := c.Bind(&req); err != nil {
		return echo.NewHTTPError(http.StatusBadRequest, "Invalid request body")
	}

	if err := h.rbacService.RemovePermission(req.Role, req.Resource, req.Action); err != nil {
		return echo.NewHTTPError(http.StatusInternalServerError, err.Error())
	}

	return c.JSON(http.StatusOK, map[string]any{
		"message":  "Permission removed successfully",
		"role":     req.Role,
		"resource": req.Resource,
		"action":   req.Action,
	})
}

// @Summary Get permissions for role
// @Tags RBAC
// @Param role path string true "Role name"
// @Success 200 {object} map[string]interface{}
// @Failure 400 {object} map[string]interface{}
// @Failure 500 {object} map[string]interface{}
// @Router /api/v1/rbac/roles/{role}/permissions [get]
func (h *RBACHandler) GetRolePermissions(c echo.Context) error {
	role := c.Param("role")
	if role == "" {
		return echo.NewHTTPError(http.StatusBadRequest, "Role parameter is required")
	}

	permissions, err := h.rbacService.GetPermissionsForRole(role)
	if err != nil {
		return echo.NewHTTPError(http.StatusInternalServerError, err.Error())
	}

	return c.JSON(http.StatusOK, map[string]any{
		"role":        role,
		"permissions": permissions,
	})
}

// @Summary Check user permission
// @Tags RBAC
// @Param user_id path string true "User ID"
// @Param resource query string true "Resource name"
// @Param action query string true "Action name"
// @Success 200 {object} map[string]interface{}
// @Failure 400 {object} map[string]interface{}
// @Failure 500 {object} map[string]interface{}
// @Router /api/v1/rbac/users/{user_id}/check-permission [get]
func (h *RBACHandler) CheckPermission(c echo.Context) error {
	userIDStr := c.Param("user_id")
	userID, err := strconv.ParseUint(userIDStr, 10, 32)
	if err != nil {
		return echo.NewHTTPError(http.StatusBadRequest, "Invalid user ID")
	}

	resource := c.QueryParam("resource")
	action := c.QueryParam("action")

	if resource == "" || action == "" {
		return echo.NewHTTPError(http.StatusBadRequest, "Resource and action parameters are required")
	}

	allowed, err := h.rbacService.CheckPermission(uint(userID), resource, action)
	if err != nil {
		return echo.NewHTTPError(http.StatusInternalServerError, err.Error())
	}

	return c.JSON(http.StatusOK, map[string]any{
		"user_id":  userID,
		"resource": resource,
		"action":   action,
		"allowed":  allowed,
	})
}
