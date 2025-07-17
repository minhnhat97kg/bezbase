package models

type ResourceType string

func (o ResourceType) String() string {
	return string(o)
}

// Define resource types for RBAC
const (
	ResourceTypeUser       ResourceType = "users"
	ResourceTypePost       ResourceType = "posts"
	ResourceTypeProfile    ResourceType = "profile"
	ResourceTypeAdmin      ResourceType = "admin"
	ResourceTypePermission ResourceType = "permissions"
	ResourceTypeRole       ResourceType = "roles"
	ResourceTypeDashboard  ResourceType = "dashboard"
	ResourceTypeSettings   ResourceType = "settings"
	ResourceTypeReports    ResourceType = "reports"
	ResourceTypeAudit      ResourceType = "audit"
	ResourceTypeSystem     ResourceType = "system"
	ResourceTypeBackup     ResourceType = "backup"
	ResourceTypeAll        ResourceType = "*"
)

type ActionType string

func (a ActionType) String() string {
	return string(a)
}

// Define action types for RBAC
const (
	ActionTypeCreate  ActionType = "create"
	ActionTypeRead    ActionType = "read"
	ActionTypeUpdate  ActionType = "update"
	ActionTypeDelete  ActionType = "delete"
	ActionTypeExport  ActionType = "export"
	ActionTypeRestore ActionType = "restore"
	ActionTypeAll     ActionType = "*"
)

// Apply pagination and get results
type Rule struct {
	ID int    `json:"id" gorm:"primaryKey"`
	V0 string `json:"v0"`
	V1 string `json:"v1"`
	V2 string `json:"v2"`
	V3 string `json:"v3"`
	V4 string `json:"v4"`
	V5 string `json:"v5"`
}
