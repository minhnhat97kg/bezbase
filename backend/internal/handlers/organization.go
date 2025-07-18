package handlers

import (
	"net/http"
	"strconv"

	"bezbase/internal/dto"
	"bezbase/internal/middleware"
	"bezbase/internal/pkg/contextx"
	"bezbase/internal/services"

	"github.com/labstack/echo/v4"
)

type OrganizationHandler struct {
	orgService *services.OrganizationService
}

func NewOrganizationHandler(orgService *services.OrganizationService) *OrganizationHandler {
	return &OrganizationHandler{
		orgService: orgService,
	}
}

// CreateOrganization creates a new organization
// @Summary Create organization
// @Description Create a new organization with the current user as owner
// @Tags Organizations
// @Accept json
// @Produce json
// @Param organization body services.CreateOrganizationRequest true "Organization data"
// @Success 201 {object} models.Organization
// @Failure 400 {object} dto.ErrorResponse
// @Failure 401 {object} dto.ErrorResponse
// @Failure 500 {object} dto.ErrorResponse
// @Router /api/v1/organizations [post]
func (h *OrganizationHandler) CreateOrganization(c echo.Context) error {
	userID := middleware.GetUserIDFromContext(c)
	if userID == nil {
		return c.JSON(http.StatusUnauthorized, dto.ErrorResponse{
			Message: "Authentication required",
		})
	}

	var req services.CreateOrganizationRequest
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

	org, err := h.orgService.CreateOrganization(contextx.NewWithRequestContext(c), *userID, req)
	if err != nil {
		return c.JSON(http.StatusInternalServerError, dto.ErrorResponse{
			Message: "Failed to create organization",
			Details: err.Error(),
		})
	}

	return c.JSON(http.StatusCreated, org)
}

// GetOrganization retrieves an organization by ID
// @Summary Get organization
// @Description Get organization details by ID
// @Tags Organizations
// @Produce json
// @Param id path int true "Organization ID"
// @Success 200 {object} models.Organization
// @Failure 400 {object} dto.ErrorResponse
// @Failure 404 {object} dto.ErrorResponse
// @Failure 500 {object} dto.ErrorResponse
// @Router /api/v1/organizations/{id} [get]
func (h *OrganizationHandler) GetOrganization(c echo.Context) error {
	orgIDStr := c.Param("id")
	orgID, err := strconv.ParseUint(orgIDStr, 10, 32)
	if err != nil {
		return c.JSON(http.StatusBadRequest, dto.ErrorResponse{
			Message: "Invalid organization ID",
		})
	}

	org, err := h.orgService.GetOrganization(contextx.NewWithRequestContext(c), uint(orgID))
	if err != nil {
		return c.JSON(http.StatusNotFound, dto.ErrorResponse{
			Message: "Organization not found",
		})
	}

	return c.JSON(http.StatusOK, org)
}

// GetUserOrganizations retrieves all organizations for the current user
// @Summary Get user organizations
// @Description Get all organizations that the current user is a member of
// @Tags Organizations
// @Produce json
// @Success 200 {array} models.OrganizationUser
// @Failure 401 {object} dto.ErrorResponse
// @Failure 500 {object} dto.ErrorResponse
// @Router /api/v1/user/organizations [get]
func (h *OrganizationHandler) GetUserOrganizations(c echo.Context) error {
	userID := middleware.GetUserIDFromContext(c)
	if userID == nil {
		return c.JSON(http.StatusUnauthorized, dto.ErrorResponse{
			Message: "Authentication required",
		})
	}

	orgs, err := h.orgService.GetUserOrganizations(contextx.NewWithRequestContext(c), *userID)
	if err != nil {
		return c.JSON(http.StatusInternalServerError, dto.ErrorResponse{
			Message: "Failed to get user organizations",
			Details: err.Error(),
		})
	}

	return c.JSON(http.StatusOK, orgs)
}

// UpdateOrganization updates organization details
// @Summary Update organization
// @Description Update organization information
// @Tags Organizations
// @Accept json
// @Produce json
// @Param id path int true "Organization ID"
// @Param organization body services.UpdateOrganizationRequest true "Organization data"
// @Success 200 {object} models.Organization
// @Failure 400 {object} dto.ErrorResponse
// @Failure 401 {object} dto.ErrorResponse
// @Failure 403 {object} dto.ErrorResponse
// @Failure 404 {object} dto.ErrorResponse
// @Failure 500 {object} dto.ErrorResponse
// @Router /api/v1/organizations/{id} [put]
func (h *OrganizationHandler) UpdateOrganization(c echo.Context) error {
	userID := middleware.GetUserIDFromContext(c)
	if userID == nil {
		return c.JSON(http.StatusUnauthorized, dto.ErrorResponse{
			Message: "Authentication required",
		})
	}

	orgIDStr := c.Param("id")
	orgID, err := strconv.ParseUint(orgIDStr, 10, 32)
	if err != nil {
		return c.JSON(http.StatusBadRequest, dto.ErrorResponse{
			Message: "Invalid organization ID",
		})
	}

	var req services.UpdateOrganizationRequest
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

	org, err := h.orgService.UpdateOrganization(contextx.NewWithRequestContext(c), uint(orgID), *userID, req)
	if err != nil {
		if err.Error() == "insufficient permissions to update organization" {
			return c.JSON(http.StatusForbidden, dto.ErrorResponse{
				Message: err.Error(),
			})
		}
		return c.JSON(http.StatusInternalServerError, dto.ErrorResponse{
			Message: "Failed to update organization",
			Details: err.Error(),
		})
	}

	return c.JSON(http.StatusOK, org)
}

// InviteUser invites a user to join the organization
// @Summary Invite user to organization
// @Description Send an invitation to a user to join the organization
// @Tags Organizations
// @Accept json
// @Produce json
// @Param id path int true "Organization ID"
// @Param invitation body services.InviteUserRequest true "Invitation data"
// @Success 201 {object} models.OrganizationInvitation
// @Failure 400 {object} dto.ErrorResponse
// @Failure 401 {object} dto.ErrorResponse
// @Failure 403 {object} dto.ErrorResponse
// @Failure 500 {object} dto.ErrorResponse
// @Router /api/v1/organizations/{id}/invite [post]
func (h *OrganizationHandler) InviteUser(c echo.Context) error {
	userID := middleware.GetUserIDFromContext(c)
	if userID == nil {
		return c.JSON(http.StatusUnauthorized, dto.ErrorResponse{
			Message: "Authentication required",
		})
	}

	orgIDStr := c.Param("id")
	orgID, err := strconv.ParseUint(orgIDStr, 10, 32)
	if err != nil {
		return c.JSON(http.StatusBadRequest, dto.ErrorResponse{
			Message: "Invalid organization ID",
		})
	}

	var req services.InviteUserRequest
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

	invitation, err := h.orgService.InviteUser(contextx.NewWithRequestContext(c), uint(orgID), *userID, req)
	if err != nil {
		if err.Error() == "insufficient permissions to invite users" {
			return c.JSON(http.StatusForbidden, dto.ErrorResponse{
				Message: err.Error(),
			})
		}
		return c.JSON(http.StatusBadRequest, dto.ErrorResponse{
			Message: err.Error(),
		})
	}

	return c.JSON(http.StatusCreated, invitation)
}

// AcceptInvitation accepts an organization invitation
// @Summary Accept organization invitation
// @Description Accept an invitation to join an organization
// @Tags Organizations
// @Accept json
// @Produce json
// @Param token path string true "Invitation token"
// @Success 200 {object} dto.SuccessResponse
// @Failure 400 {object} dto.ErrorResponse
// @Failure 401 {object} dto.ErrorResponse
// @Failure 500 {object} dto.ErrorResponse
// @Router /api/v1/organizations/invitations/{token}/accept [post]
func (h *OrganizationHandler) AcceptInvitation(c echo.Context) error {
	userID := middleware.GetUserIDFromContext(c)
	if userID == nil {
		return c.JSON(http.StatusUnauthorized, dto.ErrorResponse{
			Message: "Authentication required",
		})
	}

	token := c.Param("token")
	if token == "" {
		return c.JSON(http.StatusBadRequest, dto.ErrorResponse{
			Message: "Invalid invitation token",
		})
	}

	err := h.orgService.AcceptInvitation(contextx.NewWithRequestContext(c), token, *userID)
	if err != nil {
		return c.JSON(http.StatusBadRequest, dto.ErrorResponse{
			Message: err.Error(),
		})
	}

	return c.JSON(http.StatusOK, dto.SuccessResponse{
		Message: "Invitation accepted successfully",
	})
}

// SwitchOrganization switches user's current organization
// @Summary Switch organization context
// @Description Switch the user's current organization context
// @Tags Organizations
// @Accept json
// @Produce json
// @Param id path int true "Organization ID"
// @Success 200 {object} dto.SuccessResponse
// @Failure 400 {object} dto.ErrorResponse
// @Failure 401 {object} dto.ErrorResponse
// @Failure 403 {object} dto.ErrorResponse
// @Failure 500 {object} dto.ErrorResponse
// @Router /api/v1/organizations/{id}/switch [post]
func (h *OrganizationHandler) SwitchOrganization(c echo.Context) error {
	userID := middleware.GetUserIDFromContext(c)
	if userID == nil {
		return c.JSON(http.StatusUnauthorized, dto.ErrorResponse{
			Message: "Authentication required",
		})
	}

	orgIDStr := c.Param("id")
	orgID, err := strconv.ParseUint(orgIDStr, 10, 32)
	if err != nil {
		return c.JSON(http.StatusBadRequest, dto.ErrorResponse{
			Message: "Invalid organization ID",
		})
	}

	err = h.orgService.SwitchOrganization(*userID, uint(orgID))
	if err != nil {
		if err.Error() == "user is not a member of this organization" {
			return c.JSON(http.StatusForbidden, dto.ErrorResponse{
				Message: err.Error(),
			})
		}
		return c.JSON(http.StatusInternalServerError, dto.ErrorResponse{
			Message: "Failed to switch organization",
			Details: err.Error(),
		})
	}

	return c.JSON(http.StatusOK, dto.SuccessResponse{
		Message: "Organization switched successfully",
	})
}

// RemoveUser removes a user from the organization
// @Summary Remove user from organization
// @Description Remove a user from the organization
// @Tags Organizations
// @Param id path int true "Organization ID"
// @Param userId path int true "User ID"
// @Success 200 {object} dto.SuccessResponse
// @Failure 400 {object} dto.ErrorResponse
// @Failure 401 {object} dto.ErrorResponse
// @Failure 403 {object} dto.ErrorResponse
// @Failure 500 {object} dto.ErrorResponse
// @Router /api/v1/organizations/{id}/users/{userId} [delete]
func (h *OrganizationHandler) RemoveUser(c echo.Context) error {
	userID := middleware.GetUserIDFromContext(c)
	if userID == nil {
		return c.JSON(http.StatusUnauthorized, dto.ErrorResponse{
			Message: "Authentication required",
		})
	}

	orgIDStr := c.Param("id")
	orgID, err := strconv.ParseUint(orgIDStr, 10, 32)
	if err != nil {
		return c.JSON(http.StatusBadRequest, dto.ErrorResponse{
			Message: "Invalid organization ID",
		})
	}

	targetUserIDStr := c.Param("userId")
	targetUserID, err := strconv.ParseUint(targetUserIDStr, 10, 32)
	if err != nil {
		return c.JSON(http.StatusBadRequest, dto.ErrorResponse{
			Message: "Invalid user ID",
		})
	}

	err = h.orgService.RemoveUser(uint(orgID), uint(targetUserID), *userID)
	if err != nil {
		if err.Error() == "insufficient permissions to remove users" {
			return c.JSON(http.StatusForbidden, dto.ErrorResponse{
				Message: err.Error(),
			})
		}
		return c.JSON(http.StatusBadRequest, dto.ErrorResponse{
			Message: err.Error(),
		})
	}

	return c.JSON(http.StatusOK, dto.SuccessResponse{
		Message: "User removed from organization successfully",
	})
}

// UpdateUserRole updates a user's role in the organization
// @Summary Update user role in organization
// @Description Update a user's role within the organization
// @Tags Organizations
// @Accept json
// @Produce json
// @Param id path int true "Organization ID"
// @Param userId path int true "User ID"
// @Param role body struct{Role string `json:"role"`} true "New role"
// @Success 200 {object} dto.SuccessResponse
// @Failure 400 {object} dto.ErrorResponse
// @Failure 401 {object} dto.ErrorResponse
// @Failure 403 {object} dto.ErrorResponse
// @Failure 500 {object} dto.ErrorResponse
// @Router /api/v1/organizations/{id}/users/{userId}/role [put]
func (h *OrganizationHandler) UpdateUserRole(c echo.Context) error {
	userID := middleware.GetUserIDFromContext(c)
	if userID == nil {
		return c.JSON(http.StatusUnauthorized, dto.ErrorResponse{
			Message: "Authentication required",
		})
	}

	orgIDStr := c.Param("id")
	orgID, err := strconv.ParseUint(orgIDStr, 10, 32)
	if err != nil {
		return c.JSON(http.StatusBadRequest, dto.ErrorResponse{
			Message: "Invalid organization ID",
		})
	}

	targetUserIDStr := c.Param("userId")
	targetUserID, err := strconv.ParseUint(targetUserIDStr, 10, 32)
	if err != nil {
		return c.JSON(http.StatusBadRequest, dto.ErrorResponse{
			Message: "Invalid user ID",
		})
	}

	var req struct {
		Role string `json:"role" validate:"required,oneof=owner admin member viewer"`
	}
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

	err = h.orgService.UpdateUserRole(uint(orgID), uint(targetUserID), req.Role, *userID)
	if err != nil {
		if err.Error() == "insufficient permissions to update user roles" {
			return c.JSON(http.StatusForbidden, dto.ErrorResponse{
				Message: err.Error(),
			})
		}
		return c.JSON(http.StatusInternalServerError, dto.ErrorResponse{
			Message: "Failed to update user role",
			Details: err.Error(),
		})
	}

	return c.JSON(http.StatusOK, dto.SuccessResponse{
		Message: "User role updated successfully",
	})
}
