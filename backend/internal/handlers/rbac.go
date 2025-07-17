package handlers

import (
	"net/http"
	"strconv"
	"strings"

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
// @Security BearerAuth
// @Accept json
// @Produce json
// @Param request body dto.CreateRoleRequest true "Role creation request"
// @Success 200 {object} dto.RoleResponse
// @Failure 400 {object} map[string]interface{}
// @Failure 401 {object} map[string]interface{}
// @Failure 403 {object} map[string]interface{}
// @Failure 500 {object} map[string]interface{}
// @Router /v1/rbac/roles [post]
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

// @Summary List all roles with pagination
// @Tags RBAC
// @Security BearerAuth
// @Produce json
// @Param page query int false "Page number (default: 1)"
// @Param page_size query int false "Page size (default: 10, max: 100)"
// @Param search query string false "Search by name or display name"
// @Param status query string false "Filter by status (active, inactive)"
// @Param is_system query bool false "Filter by system roles"
// @Param sort query string false "Sort field (name, display_name, created_at)"
// @Param order query string false "Sort order (asc, desc)"
// @Success 200 {object} dto.RolesListResponse
// @Failure 400 {object} map[string]interface{}
// @Failure 401 {object} map[string]interface{}
// @Failure 403 {object} map[string]interface{}
// @Failure 500 {object} map[string]interface{}
// @Router /v1/rbac/roles [get]
func (h *RBACHandler) GetRoles(c echo.Context) error {
	// Parse pagination parameters
	var pagination dto.PaginationParams
	if err := c.Bind(&pagination); err != nil {
		return echo.NewHTTPError(http.StatusBadRequest, "Invalid pagination parameters")
	}
	pagination.SetDefaults()

	// Parse filter parameters
	searchFilter := c.QueryParam("search")
	statusFilter := c.QueryParam("status")
	isSystemParam := c.QueryParam("is_system")
	
	var isSystemFilter *bool
	if isSystemParam != "" {
		isSystem, err := strconv.ParseBool(isSystemParam)
		if err == nil {
			isSystemFilter = &isSystem
		}
	}

	// Parse sort parameters
	sortField := c.QueryParam("sort")
	sortOrder := c.QueryParam("order")
	if sortOrder != "desc" {
		sortOrder = "asc"
	}

	roles, total, err := h.rbacService.GetAllRolesWithPagination(pagination.Page, pagination.PageSize, searchFilter, statusFilter, isSystemFilter, sortField, sortOrder)
	if err != nil {
		return echo.NewHTTPError(http.StatusInternalServerError, err.Error())
	}

	response := dto.NewPaginatedResponse(dto.ToRoleResponses(roles), pagination.Page, pagination.PageSize, int64(total))

	return c.JSON(http.StatusOK, response)
}

// @Summary Get role by ID
// @Tags RBAC
// @Security BearerAuth
// @Produce json
// @Param role_id path string true "Role ID"
// @Success 200 {object} dto.RoleResponse
// @Failure 400 {object} map[string]interface{}
// @Failure 401 {object} map[string]interface{}
// @Failure 403 {object} map[string]interface{}
// @Failure 404 {object} map[string]interface{}
// @Failure 500 {object} map[string]interface{}
// @Router /v1/rbac/roles/{role_id} [get]
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
// @Security BearerAuth
// @Accept json
// @Produce json
// @Param role_id path string true "Role ID"
// @Param request body dto.UpdateRoleRequest true "Role update request"
// @Success 200 {object} dto.RoleResponse
// @Failure 400 {object} map[string]interface{}
// @Failure 401 {object} map[string]interface{}
// @Failure 403 {object} map[string]interface{}
// @Failure 404 {object} map[string]interface{}
// @Failure 500 {object} map[string]interface{}
// @Router /v1/rbac/roles/{role_id} [put]
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
// @Security BearerAuth
// @Produce json
// @Param role path string true "Role name"
// @Success 200 {object} map[string]interface{}
// @Failure 400 {object} map[string]interface{}
// @Failure 401 {object} map[string]interface{}
// @Failure 403 {object} map[string]interface{}
// @Failure 500 {object} map[string]interface{}
// @Router /v1/rbac/roles/{role} [delete]
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
// @Security BearerAuth
// @Accept json
// @Produce json
// @Param request body dto.AssignRoleRequest true "Role assignment request"
// @Success 200 {object} map[string]interface{}
// @Failure 400 {object} map[string]interface{}
// @Failure 401 {object} map[string]interface{}
// @Failure 403 {object} map[string]interface{}
// @Failure 500 {object} map[string]interface{}
// @Router /v1/rbac/users/assign-role [post]
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
// @Security BearerAuth
// @Accept json
// @Produce json
// @Param request body dto.AssignRoleRequest true "Role removal request"
// @Success 200 {object} map[string]interface{}
// @Failure 400 {object} map[string]interface{}
// @Failure 401 {object} map[string]interface{}
// @Failure 403 {object} map[string]interface{}
// @Failure 500 {object} map[string]interface{}
// @Router /v1/rbac/users/remove-role [post]
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
// @Security BearerAuth
// @Produce json
// @Param user_id path string true "User ID"
// @Success 200 {object} map[string]interface{}
// @Failure 400 {object} map[string]interface{}
// @Failure 401 {object} map[string]interface{}
// @Failure 403 {object} map[string]interface{}
// @Failure 500 {object} map[string]interface{}
// @Router /v1/rbac/users/{user_id}/roles [get]
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
// @Security BearerAuth
// @Produce json
// @Param role path string true "Role name"
// @Success 200 {object} map[string]interface{}
// @Failure 400 {object} map[string]interface{}
// @Failure 401 {object} map[string]interface{}
// @Failure 403 {object} map[string]interface{}
// @Failure 500 {object} map[string]interface{}
// @Router /v1/rbac/roles/{role}/users [get]
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
// @Security BearerAuth
// @Produce json
// @Param page query int false "Page number (default: 1)"
// @Param page_size query int false "Page size (default: 10, max: 100)"
// @Param role query string false "Filter by role"
// @Param resource query string false "Filter by resource"
// @Param action query string false "Filter by action"
// @Param sort query string false "Sort field (role, resource, action)"
// @Param order query string false "Sort order (asc, desc)"
// @Success 200 {object} dto.PermissionsListResponse
// @Failure 400 {object} map[string]interface{}
// @Failure 401 {object} map[string]interface{}
// @Failure 403 {object} map[string]interface{}
// @Failure 500 {object} map[string]interface{}
// @Router /v1/rbac/permissions [get]
func (h *RBACHandler) GetPermissions(c echo.Context) error {
	// Parse pagination parameters
	var pagination dto.PaginationParams
	if err := c.Bind(&pagination); err != nil {
		return echo.NewHTTPError(http.StatusBadRequest, "Invalid pagination parameters")
	}
	pagination.SetDefaults()

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

	permissions, total, err := h.rbacService.GetAllPermissions(pagination.Page, pagination.PageSize, roleFilter, resourceFilter, actionFilter, sortField, sortOrder)
	if err != nil {
		return echo.NewHTTPError(http.StatusInternalServerError, err.Error())
	}

	response := dto.NewPaginatedResponse(permissions, pagination.Page, pagination.PageSize, int64(total))

	return c.JSON(http.StatusOK, response)
}

// @Summary Add permission to role
// @Tags RBAC
// @Security BearerAuth
// @Accept json
// @Produce json
// @Param request body dto.PermissionRequest true "Permission request"
// @Success 200 {object} map[string]interface{}
// @Failure 400 {object} map[string]interface{}
// @Failure 401 {object} map[string]interface{}
// @Failure 403 {object} map[string]interface{}
// @Failure 500 {object} map[string]interface{}
// @Router /v1/rbac/permissions [post]
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
// @Security BearerAuth
// @Accept json
// @Produce json
// @Param request body dto.PermissionRequest true "Permission request"
// @Success 200 {object} map[string]interface{}
// @Failure 400 {object} map[string]interface{}
// @Failure 401 {object} map[string]interface{}
// @Failure 403 {object} map[string]interface{}
// @Failure 500 {object} map[string]interface{}
// @Router /v1/rbac/permissions [delete]
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
// @Security BearerAuth
// @Produce json
// @Param role path string true "Role name"
// @Success 200 {object} map[string]interface{}
// @Failure 400 {object} map[string]interface{}
// @Failure 401 {object} map[string]interface{}
// @Failure 403 {object} map[string]interface{}
// @Failure 500 {object} map[string]interface{}
// @Router /v1/rbac/roles/{role}/permissions [get]
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
// @Security BearerAuth
// @Produce json
// @Param user_id path string true "User ID"
// @Param resource query string true "Resource name"
// @Param action query string true "Action name"
// @Success 200 {object} map[string]interface{}
// @Failure 400 {object} map[string]interface{}
// @Failure 401 {object} map[string]interface{}
// @Failure 403 {object} map[string]interface{}
// @Failure 500 {object} map[string]interface{}
// @Router /v1/rbac/users/{user_id}/check-permission [get]
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

// @Summary Get paginated list of resources
// @Tags RBAC
// @Security BearerAuth
// @Produce json
// @Param page query int false "Page number" default(1)
// @Param page_size query int false "Page size" default(10)
// @Param search query string false "Search term"
// @Success 200 {object} dto.PaginatedResourceResponse
// @Failure 400 {object} map[string]interface{}
// @Failure 401 {object} map[string]interface{}
// @Failure 403 {object} map[string]interface{}
// @Failure 500 {object} map[string]interface{}
// @Router /v1/rbac/resources [get]
func (h *RBACHandler) GetResources(c echo.Context) error {
	pagination := dto.ParsePagination(c)
	search := c.QueryParam("search")

	// Get all resources
	allResources := dto.GetAllResources()

	// Apply filters
	var filteredResources []dto.ResourceResponse
	for _, resource := range allResources {
		matchesSearch := search == "" || 
			strings.Contains(strings.ToLower(resource.Name), strings.ToLower(search)) ||
			strings.Contains(strings.ToLower(resource.Description), strings.ToLower(search))
		
		if matchesSearch {
			filteredResources = append(filteredResources, resource)
		}
	}

	// Calculate pagination
	totalItems := len(filteredResources)
	totalPages := (totalItems + pagination.PageSize - 1) / pagination.PageSize
	
	// Apply pagination
	startIndex := (pagination.Page - 1) * pagination.PageSize
	endIndex := startIndex + pagination.PageSize
	
	if startIndex > totalItems {
		startIndex = totalItems
	}
	if endIndex > totalItems {
		endIndex = totalItems
	}
	
	var paginatedResources []dto.ResourceResponse
	if startIndex < totalItems {
		paginatedResources = filteredResources[startIndex:endIndex]
	}

	response := dto.PaginatedResourceResponse{
		Data:       paginatedResources,
		Page:       pagination.Page,
		PageSize:   pagination.PageSize,
		TotalItems: totalItems,
		TotalPages: totalPages,
	}

	return c.JSON(http.StatusOK, response)
}

// @Summary Get paginated list of actions
// @Tags RBAC
// @Security BearerAuth
// @Produce json
// @Param page query int false "Page number" default(1)
// @Param page_size query int false "Page size" default(10)
// @Param search query string false "Search term"
// @Success 200 {object} dto.PaginatedActionResponse
// @Failure 400 {object} map[string]interface{}
// @Failure 401 {object} map[string]interface{}
// @Failure 403 {object} map[string]interface{}
// @Failure 500 {object} map[string]interface{}
// @Router /v1/rbac/actions [get]
func (h *RBACHandler) GetActions(c echo.Context) error {
	pagination := dto.ParsePagination(c)
	search := c.QueryParam("search")

	// Get all actions
	allActions := dto.GetAllActions()

	// Apply filters
	var filteredActions []dto.ActionResponse
	for _, action := range allActions {
		matchesSearch := search == "" || 
			strings.Contains(strings.ToLower(action.Name), strings.ToLower(search)) ||
			strings.Contains(strings.ToLower(action.Description), strings.ToLower(search))
		
		if matchesSearch {
			filteredActions = append(filteredActions, action)
		}
	}

	// Calculate pagination
	totalItems := len(filteredActions)
	totalPages := (totalItems + pagination.PageSize - 1) / pagination.PageSize
	
	// Apply pagination
	startIndex := (pagination.Page - 1) * pagination.PageSize
	endIndex := startIndex + pagination.PageSize
	
	if startIndex > totalItems {
		startIndex = totalItems
	}
	if endIndex > totalItems {
		endIndex = totalItems
	}
	
	var paginatedActions []dto.ActionResponse
	if startIndex < totalItems {
		paginatedActions = filteredActions[startIndex:endIndex]
	}

	response := dto.PaginatedActionResponse{
		Data:       paginatedActions,
		Page:       pagination.Page,
		PageSize:   pagination.PageSize,
		TotalItems: totalItems,
		TotalPages: totalPages,
	}

	return c.JSON(http.StatusOK, response)
}
