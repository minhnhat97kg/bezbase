package models

// Permission represents a resource-action permission
type Permission struct {
	Resource   ResourceType `json:"resource"`
	Action     ActionType   `json:"action"`
	Permission string       `json:"permission"`
}

var (
	PermissionCreateUsers         = Permission{Resource: ResourceTypeUser, Action: ActionTypeCreate, Permission: "Create Users"}
	PermissionViewUsers           = Permission{Resource: ResourceTypeUser, Action: ActionTypeRead, Permission: "View Users"}
	PermissionEditUsers           = Permission{Resource: ResourceTypeUser, Action: ActionTypeUpdate, Permission: "Edit Users"}
	PermissionDeleteUsers         = Permission{Resource: ResourceTypeUser, Action: ActionTypeDelete, Permission: "Delete Users"}
	PermissionCreateRoles         = Permission{Resource: ResourceTypeRole, Action: ActionTypeCreate, Permission: "Create Roles"}
	PermissionViewRoles           = Permission{Resource: ResourceTypeRole, Action: ActionTypeRead, Permission: "View Roles"}
	PermissionEditRoles           = Permission{Resource: ResourceTypeRole, Action: ActionTypeUpdate, Permission: "Edit Roles"}
	PermissionDeleteRoles         = Permission{Resource: ResourceTypeRole, Action: ActionTypeDelete, Permission: "Delete Roles"}
	PermissionCreatePermissions   = Permission{Resource: ResourceTypePermission, Action: ActionTypeCreate, Permission: "Create Permissions"}
	PermissionViewPermissions     = Permission{Resource: ResourceTypePermission, Action: ActionTypeRead, Permission: "View Permissions"}
	PermissionEditPermissions     = Permission{Resource: ResourceTypePermission, Action: ActionTypeUpdate, Permission: "Edit Permissions"}
	PermissionDeletePermissions   = Permission{Resource: ResourceTypePermission, Action: ActionTypeDelete, Permission: "Delete Permissions"}
	PermissionCreateOrganizations = Permission{Resource: ResourceTypeOrganization, Action: ActionTypeCreate, Permission: "Create Organizations"}
	PermissionViewOrganizations   = Permission{Resource: ResourceTypeOrganization, Action: ActionTypeRead, Permission: "View Organizations"}
	PermissionEditOrganizations   = Permission{Resource: ResourceTypeOrganization, Action: ActionTypeUpdate, Permission: "Edit Organizations"}
	PermissionDeleteOrganizations = Permission{Resource: ResourceTypeOrganization, Action: ActionTypeDelete, Permission: "Delete Organizations"}
	PermissionViewDashboard       = Permission{Resource: ResourceTypeDashboard, Action: ActionTypeRead, Permission: "View Dashboard"}
	PermissionViewProfile         = Permission{Resource: ResourceTypeProfile, Action: ActionTypeRead, Permission: "View Profile"}
	PermissionEditProfile         = Permission{Resource: ResourceTypeProfile, Action: ActionTypeUpdate, Permission: "Edit Profile"}
)

// GetHardcodedPermissions returns a hardcoded list of permissions
func GetHardcodedPermissions() []Permission {
	return []Permission{
		PermissionCreateUsers,
		PermissionViewUsers,
		PermissionEditUsers,
		PermissionDeleteUsers,
		PermissionCreateRoles,
		PermissionViewRoles,
		PermissionEditRoles,
		PermissionDeleteRoles,
		PermissionCreatePermissions,
		PermissionViewPermissions,
		PermissionEditPermissions,
		PermissionDeletePermissions,
		PermissionCreateOrganizations,
		PermissionViewOrganizations,
		PermissionEditOrganizations,
		PermissionDeleteOrganizations,
		PermissionViewDashboard,
		PermissionViewProfile,
		PermissionEditProfile,
	}
}

// GetPermissionByResourceAction finds a specific permission by resource and action
func GetPermissionByResourceAction(resource ResourceType, action ActionType) (Permission, bool) {
	permissions := GetHardcodedPermissions()
	for _, perm := range permissions {
		if perm.Resource == resource && perm.Action == action {
			return perm, true
		}
	}
	return Permission{}, false
}

// Common permission constants for easy access
