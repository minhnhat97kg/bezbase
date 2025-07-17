package dto

import "bezbase/internal/models"

type ResourceResponse struct {
	ID          string `json:"id"`
	Name        string `json:"name"`
	Description string `json:"description"`
	CreatedAt   string `json:"created_at"`
}

type ActionResponse struct {
	ID          string `json:"id"`
	Name        string `json:"name"`
	Description string `json:"description"`
	CreatedAt   string `json:"created_at"`
}

type PaginatedResourceResponse struct {
	Data        []ResourceResponse `json:"data"`
	Page        int                `json:"page"`
	PageSize    int                `json:"page_size"`
	TotalItems  int                `json:"total_items"`
	TotalPages  int                `json:"total_pages"`
}

type PaginatedActionResponse struct {
	Data        []ActionResponse `json:"data"`
	Page        int              `json:"page"`
	PageSize    int              `json:"page_size"`
	TotalItems  int              `json:"total_items"`
	TotalPages  int              `json:"total_pages"`
}

// Hardcoded resource data
func GetAllResources() []ResourceResponse {
	return []ResourceResponse{
		{
			ID:          string(models.ResourceTypeUser),
			Name:        "Users",
			Description: "User account management and profile operations",
			CreatedAt:   "2024-01-01T00:00:00Z",
		},
		{
			ID:          string(models.ResourceTypePost),
			Name:        "Posts",
			Description: "Blog posts and content management",
			CreatedAt:   "2024-01-01T00:00:00Z",
		},
		{
			ID:          string(models.ResourceTypeProfile),
			Name:        "Profile",
			Description: "User profile information and settings",
			CreatedAt:   "2024-01-01T00:00:00Z",
		},
		{
			ID:          string(models.ResourceTypeAdmin),
			Name:        "Admin",
			Description: "Administrative functions and system settings",
			CreatedAt:   "2024-01-01T00:00:00Z",
		},
		{
			ID:          string(models.ResourceTypePermission),
			Name:        "Permissions",
			Description: "Role and permission management",
			CreatedAt:   "2024-01-01T00:00:00Z",
		},
		{
			ID:          string(models.ResourceTypeAll),
			Name:        "All Resources",
			Description: "Global access to all system resources",
			CreatedAt:   "2024-01-01T00:00:00Z",
		},
	}
}

// Hardcoded action data
func GetAllActions() []ActionResponse {
	return []ActionResponse{
		{
			ID:          string(models.ActionTypeCreate),
			Name:        "Create",
			Description: "Create new resources and entities",
			CreatedAt:   "2024-01-01T00:00:00Z",
		},
		{
			ID:          string(models.ActionTypeRead),
			Name:        "Read",
			Description: "View and retrieve existing resources",
			CreatedAt:   "2024-01-01T00:00:00Z",
		},
		{
			ID:          string(models.ActionTypeUpdate),
			Name:        "Update",
			Description: "Modify existing resources and entities",
			CreatedAt:   "2024-01-01T00:00:00Z",
		},
		{
			ID:          string(models.ActionTypeDelete),
			Name:        "Delete",
			Description: "Remove resources and entities",
			CreatedAt:   "2024-01-01T00:00:00Z",
		},
		{
			ID:          string(models.ActionTypeAll),
			Name:        "All Actions",
			Description: "Permission to perform all actions",
			CreatedAt:   "2024-01-01T00:00:00Z",
		},
	}
}