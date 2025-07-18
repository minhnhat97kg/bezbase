package handlers

import (
	"net/http"
	"strconv"

	"bezbase/internal/dto"
	"bezbase/internal/middleware"
	"bezbase/internal/models"
	"bezbase/internal/pkg/contextx"
	"bezbase/internal/repository"
	"bezbase/internal/services"

	"github.com/labstack/echo/v4"
	"gorm.io/gorm"
)

type AdvancedRBACHandler struct {
	rbacService         *services.RBACService
	roleTemplateRepo    repository.RoleTemplateRepository
	contextualPermRepo  repository.ContextualPermissionRepository
	roleInheritanceRepo repository.RoleInheritanceRepository
	db                  *gorm.DB
}

func NewAdvancedRBACHandler(
	rbacService *services.RBACService,
	roleTemplateRepo repository.RoleTemplateRepository,
	contextualPermRepo repository.ContextualPermissionRepository,
	roleInheritanceRepo repository.RoleInheritanceRepository,
	db *gorm.DB,
) *AdvancedRBACHandler {
	return &AdvancedRBACHandler{
		rbacService:         rbacService,
		roleTemplateRepo:    roleTemplateRepo,
		contextualPermRepo:  contextualPermRepo,
		roleInheritanceRepo: roleInheritanceRepo,
		db:                  db,
	}
}

// CreateRoleFromTemplate creates a role from a template
// @Summary Create role from template
// @Description Create a new role based on a predefined template
// @Tags Advanced RBAC
// @Accept json
// @Produce json
// @Param request body CreateRoleFromTemplateRequest true "Role creation request"
// @Success 201 {object} models.Role
// @Failure 400 {object} dto.ErrorResponse
// @Failure 401 {object} dto.ErrorResponse
// @Failure 500 {object} dto.ErrorResponse
// @Router /api/v1/rbac/roles/from-template [post]
func (h *AdvancedRBACHandler) CreateRoleFromTemplate(c echo.Context) error {
	userID := middleware.GetUserIDFromContext(c)
	if userID == nil {
		return c.JSON(http.StatusUnauthorized, dto.ErrorResponse{
			Message: "Authentication required",
		})
	}

	orgID := middleware.GetOrganizationIDFromContext(c)

	var req CreateRoleFromTemplateRequest
	if err := c.Bind(&req); err != nil {
		return c.JSON(http.StatusBadRequest, dto.ErrorResponse{
			Message: "Invalid request format",
		})
	}

	if err := c.Validate(&req); err != nil {
		return c.JSON(http.StatusBadRequest, dto.ErrorResponse{
			Message: "Validation failed",
			Details: err.Error(),
		})
	}

	role, err := h.rbacService.CreateRoleFromTemplate(req.TemplateID, orgID, req.CustomName)
	if err != nil {
		return c.JSON(http.StatusInternalServerError, dto.ErrorResponse{
			Message: "Failed to create role from template",
			Details: err.Error(),
		})
	}

	return c.JSON(http.StatusCreated, role)
}

// SetRoleParent sets or updates role hierarchy
// @Summary Set role parent
// @Description Set parent role for hierarchical inheritance
// @Tags Advanced RBAC
// @Accept json
// @Produce json
// @Param id path int true "Role ID"
// @Param request body SetRoleParentRequest true "Parent role request"
// @Success 200 {object} dto.SuccessResponse
// @Failure 400 {object} dto.ErrorResponse
// @Failure 401 {object} dto.ErrorResponse
// @Failure 500 {object} dto.ErrorResponse
// @Router /api/v1/rbac/roles/{id}/parent [put]
func (h *AdvancedRBACHandler) SetRoleParent(c echo.Context) error {
	userID := middleware.GetUserIDFromContext(c)
	if userID == nil {
		return c.JSON(http.StatusUnauthorized, dto.ErrorResponse{
			Message: "Authentication required",
		})
	}

	roleIDStr := c.Param("id")
	roleID, err := strconv.ParseUint(roleIDStr, 10, 32)
	if err != nil {
		return c.JSON(http.StatusBadRequest, dto.ErrorResponse{
			Message: "Invalid role ID",
		})
	}

	var req SetRoleParentRequest
	if err := c.Bind(&req); err != nil {
		return c.JSON(http.StatusBadRequest, dto.ErrorResponse{
			Message: "Invalid request format",
		})
	}

	err = h.rbacService.SetRoleParent(uint(roleID), req.ParentRoleID)
	if err != nil {
		return c.JSON(http.StatusInternalServerError, dto.ErrorResponse{
			Message: "Failed to set role parent",
			Details: err.Error(),
		})
	}

	return c.JSON(http.StatusOK, dto.SuccessResponse{
		Message: "Role hierarchy updated successfully",
	})
}

// GetRolesByOrganization retrieves roles for an organization
// @Summary Get organization roles
// @Description Get all roles available in an organization context
// @Tags Advanced RBAC
// @Produce json
// @Param org_id query int false "Organization ID"
// @Success 200 {array} models.Role
// @Failure 401 {object} dto.ErrorResponse
// @Failure 500 {object} dto.ErrorResponse
// @Router /api/v1/rbac/roles [get]
func (h *AdvancedRBACHandler) GetRolesByOrganization(c echo.Context) error {
	userID := middleware.GetUserIDFromContext(c)
	if userID == nil {
		return c.JSON(http.StatusUnauthorized, dto.ErrorResponse{
			Message: "Authentication required",
		})
	}

	var orgID *uint
	if orgIDStr := c.QueryParam("org_id"); orgIDStr != "" {
		if id, err := strconv.ParseUint(orgIDStr, 10, 32); err == nil {
			orgIDValue := uint(id)
			orgID = &orgIDValue
		}
	}

	// If no org_id provided, try to get from context
	if orgID == nil {
		orgID = middleware.GetOrganizationIDFromContext(c)
	}

	roles, err := h.rbacService.GetRolesByOrganization(orgID)
	if err != nil {
		return c.JSON(http.StatusInternalServerError, dto.ErrorResponse{
			Message: "Failed to get roles",
			Details: err.Error(),
		})
	}

	return c.JSON(http.StatusOK, roles)
}

// GetRoleTemplates retrieves available role templates
// @Summary Get role templates
// @Description Get all available role templates
// @Tags Advanced RBAC
// @Produce json
// @Param category query string false "Template category"
// @Success 200 {array} models.RoleTemplate
// @Failure 500 {object} dto.ErrorResponse
// @Router /api/v1/rbac/role-templates [get]
func (h *AdvancedRBACHandler) GetRoleTemplates(c echo.Context) error {
	ctxx := contextx.NewWithRequestContext(c)
	category := c.QueryParam("category")

	var templates []models.RoleTemplate
	var err error

	if category != "" {
		templates, err = h.roleTemplateRepo.GetByCategory(ctxx, category)
	} else {
		templates, err = h.roleTemplateRepo.GetActive(ctxx)
	}

	if err != nil {
		return c.JSON(http.StatusInternalServerError, dto.ErrorResponse{
			Message: "Failed to get role templates",
			Details: err.Error(),
		})
	}

	return c.JSON(http.StatusOK, templates)
}

// CreateContextualPermission creates a contextual permission
// @Summary Create contextual permission
// @Description Create a context-aware permission for a role
// @Tags Advanced RBAC
// @Accept json
// @Produce json
// @Param permission body CreateContextualPermissionRequest true "Permission data"
// @Success 201 {object} models.ContextualPermission
// @Failure 400 {object} dto.ErrorResponse
// @Failure 401 {object} dto.ErrorResponse
// @Failure 500 {object} dto.ErrorResponse
// @Router /api/v1/rbac/contextual-permissions [post]
func (h *AdvancedRBACHandler) CreateContextualPermission(c echo.Context) error {
	ctxx := contextx.NewWithRequestContext(c)
	userID := middleware.GetUserIDFromContext(c)
	if userID == nil {
		return c.JSON(http.StatusUnauthorized, dto.ErrorResponse{
			Message: "Authentication required",
		})
	}

	var req CreateContextualPermissionRequest
	if err := c.Bind(&req); err != nil {
		return c.JSON(http.StatusBadRequest, dto.ErrorResponse{
			Message: "Invalid request format",
		})
	}

	if err := c.Validate(&req); err != nil {
		return c.JSON(http.StatusBadRequest, dto.ErrorResponse{
			Message: "Validation failed",
			Details: err.Error(),
		})
	}

	permission := &models.ContextualPermission{
		RoleID:       req.RoleID,
		Resource:     req.Resource,
		Action:       req.Action,
		ContextType:  req.ContextType,
		ContextValue: req.ContextValue,
		IsGranted:    req.IsGranted,
	}

	if err := h.contextualPermRepo.Create(ctxx, permission); err != nil {
		return c.JSON(http.StatusInternalServerError, dto.ErrorResponse{
			Message: "Failed to create contextual permission",
			Details: err.Error(),
		})
	}

	return c.JSON(http.StatusCreated, permission)
}

// GetEffectivePermissions gets effective permissions for a user
// @Summary Get effective permissions
// @Description Get all effective permissions for a user in organization context
// @Tags Advanced RBAC
// @Produce json
// @Param userId path int true "User ID"
// @Param org_id query int false "Organization ID"
// @Success 200 {array} models.ContextualPermission
// @Failure 400 {object} dto.ErrorResponse
// @Failure 401 {object} dto.ErrorResponse
// @Failure 500 {object} dto.ErrorResponse
// @Router /api/v1/rbac/users/{userId}/effective-permissions [get]
func (h *AdvancedRBACHandler) GetEffectivePermissions(c echo.Context) error {
	requestingUserID := middleware.GetUserIDFromContext(c)
	if requestingUserID == nil {
		return c.JSON(http.StatusUnauthorized, dto.ErrorResponse{
			Message: "Authentication required",
		})
	}

	userIDStr := c.Param("userId")
	userID, err := strconv.ParseUint(userIDStr, 10, 32)
	if err != nil {
		return c.JSON(http.StatusBadRequest, dto.ErrorResponse{
			Message: "Invalid user ID",
		})
	}

	var orgID *uint
	if orgIDStr := c.QueryParam("org_id"); orgIDStr != "" {
		if id, err := strconv.ParseUint(orgIDStr, 10, 32); err == nil {
			orgIDValue := uint(id)
			orgID = &orgIDValue
		}
	}

	// If no org_id provided, try to get from context
	if orgID == nil {
		orgID = middleware.GetOrganizationIDFromContext(c)
	}

	// Get user's roles in organization context
	roles, err := h.rbacService.GetUserRolesInOrganization(uint(userID), orgID)
	if err != nil {
		return c.JSON(http.StatusInternalServerError, dto.ErrorResponse{
			Message: "Failed to get user roles",
			Details: err.Error(),
		})
	}

	// Collect all effective permissions
	var allPermissions []models.ContextualPermission
	for _, role := range roles {
		// Add special rule: if role is 'agent', grant 'agent:access' permission
		if role.Name == "agent" {
			allPermissions = append(allPermissions, models.ContextualPermission{
				RoleID:       role.ID,
				Resource:     "agent",
				Action:       "access",
				ContextType:  "organization",
				ContextValue: orgIDString(orgID),
				IsGranted:    true,
			})
		}
		permissions, err := h.contextualPermRepo.GetEffectivePermissions(contextx.NewWithRequestContext(c), role.ID, orgID)
		if err != nil {
			continue
		}
		allPermissions = append(allPermissions, permissions...)
	}

	return c.JSON(http.StatusOK, allPermissions)
}

// orgIDString returns the string value of orgID or empty string if nil
func orgIDString(orgID *uint) string {
	if orgID == nil {
		return ""
	}
	return strconv.FormatUint(uint64(*orgID), 10)
}

// GetRoleHierarchy gets the role hierarchy for a role
// @Summary Get role hierarchy
// @Description Get parent and child roles in hierarchy
// @Tags Advanced RBAC
// @Produce json
// @Param id path int true "Role ID"
// @Success 200 {object} RoleHierarchyResponse
// @Failure 400 {object} dto.ErrorResponse
// @Failure 500 {object} dto.ErrorResponse
// @Router /api/v1/rbac/roles/{id}/hierarchy [get]
func (h *AdvancedRBACHandler) GetRoleHierarchy(c echo.Context) error {
	roleIDStr := c.Param("id")
	roleID, err := strconv.ParseUint(roleIDStr, 10, 32)
	if err != nil {
		return c.JSON(http.StatusBadRequest, dto.ErrorResponse{
			Message: "Invalid role ID",
		})
	}

	// Get parent roles
	parentRoles, err := models.GetAllParentRoles(h.db, uint(roleID))
	if err != nil {
		return c.JSON(http.StatusInternalServerError, dto.ErrorResponse{
			Message: "Failed to get parent roles",
			Details: err.Error(),
		})
	}

	// Get child roles
	childRoles, err := models.GetAllChildRoles(h.db, uint(roleID))
	if err != nil {
		return c.JSON(http.StatusInternalServerError, dto.ErrorResponse{
			Message: "Failed to get child roles",
			Details: err.Error(),
		})
	}

	response := RoleHierarchyResponse{
		RoleID:      uint(roleID),
		ParentRoles: parentRoles,
		ChildRoles:  childRoles,
	}

	return c.JSON(http.StatusOK, response)
}

// Request/Response DTOs
type CreateRoleFromTemplateRequest struct {
	TemplateID uint   `json:"template_id" validate:"required"`
	CustomName string `json:"custom_name"`
}

type SetRoleParentRequest struct {
	ParentRoleID *uint `json:"parent_role_id"`
}

type CreateContextualPermissionRequest struct {
	RoleID       uint   `json:"role_id" validate:"required"`
	Resource     string `json:"resource" validate:"required"`
	Action       string `json:"action" validate:"required"`
	ContextType  string `json:"context_type"`
	ContextValue string `json:"context_value"`
	IsGranted    bool   `json:"is_granted"`
}

type RoleHierarchyResponse struct {
	RoleID      uint          `json:"role_id"`
	ParentRoles []models.Role `json:"parent_roles"`
	ChildRoles  []models.Role `json:"child_roles"`
}
