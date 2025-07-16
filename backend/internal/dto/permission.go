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

type PermissionsListResponse = PaginatedResponse[PermissionResponse]
