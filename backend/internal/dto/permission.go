package dto

// Role management endpoints
type AssignRoleRequest struct {
	UserID uint   `json:"user_id" validate:"required"`
	Role   string `json:"role" validate:"required"`
}

type PermissionRequest struct {
	Role     string `json:"role" validate:"required"`
	Resource string `json:"resource" validate:"required"`
	Action   string `json:"action" validate:"required"`
}

type PermissionResponse struct {
	ID         int    `json:"id"`
	Role       string `json:"role"`
	Resource   string `json:"resource"`
	Action     string `json:"action"`
	Permission string `json:"permission"`
}

type PermissionsListResponse = PaginatedResponse[PermissionResponse]
