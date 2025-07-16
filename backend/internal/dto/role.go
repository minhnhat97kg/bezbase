package dto

import (
	"time"

	"bezbase/internal/models"
)

// Role requests

type CreateRoleRequest struct {
	Name        string `json:"name" validate:"required,min=2,max=100"`
	DisplayName string `json:"display_name" validate:"required,min=2,max=255"`
	Description string `json:"description" validate:"max=500"`
	IsActive    *bool  `json:"is_active,omitempty"`
}

type UpdateRoleRequest struct {
	DisplayName string `json:"display_name" validate:"min=2,max=255"`
	Description string `json:"description" validate:"max=500"`
	IsActive    *bool  `json:"is_active,omitempty"`
}

// Role responses

type RoleResponse struct {
	ID          uint      `json:"id"`
	Name        string    `json:"name"`
	DisplayName string    `json:"display_name"`
	Description string    `json:"description"`
	IsSystem    bool      `json:"is_system"`
	IsActive    bool      `json:"is_active"`
	CreatedAt   time.Time `json:"created_at"`
	UpdatedAt   time.Time `json:"updated_at"`
}

type RoleWithPermissions struct {
	RoleResponse
	Permissions []PermissionResponse `json:"permissions"`
}

type PermissionResponse struct {
	ID       int    `json:"id"`
	Role     string `json:"role"`
	Resource string `json:"resource"`
	Action   string `json:"action"`
}

type RolesListResponse = PaginatedResponse[RoleResponse]

func ToRoleResponse(role *models.Role) RoleResponse {
	return RoleResponse{
		ID:          role.ID,
		Name:        role.Name,
		DisplayName: role.DisplayName,
		Description: role.Description,
		IsSystem:    role.IsSystem,
		IsActive:    role.IsActive,
		CreatedAt:   role.CreatedAt,
		UpdatedAt:   role.UpdatedAt,
	}
}

func ToRoleResponses(roles []models.Role) []RoleResponse {
	responses := make([]RoleResponse, len(roles))
	for i, role := range roles {
		responses[i] = ToRoleResponse(&role)
	}
	return responses
}
