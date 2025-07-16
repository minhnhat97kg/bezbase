# RBAC Implementation Usage Guide

This guide shows how to use the RBAC (Role-Based Access Control) system implemented with Casbin and role table management.

## Overview

The RBAC system provides:
- Role-based access control using Casbin
- Database storage for permissions and roles with metadata
- Role table with rich metadata (display name, description, system flags)
- REST API endpoints for managing permissions and roles
- Middleware for authorization with role validation

## Role Table Structure

The system uses a dedicated `roles` table with the following fields:
- `id`: Primary key
- `name`: Unique role name (used in Casbin)
- `display_name`: Human-readable role name
- `description`: Role description
- `is_system`: System roles (cannot be deleted)
- `is_active`: Active status (inactive roles cannot be assigned)
- `created_at`, `updated_at`, `deleted_at`: Timestamps

## Default Roles and Permissions

The system initializes with these default roles:
- `admin`: Administrator with full system access
- `moderator`: Content moderation and user management
- `user`: Regular user with basic permissions

Default permissions:
- `admin`: `*:*` (all resources, all actions)
- `moderator`: `users:read,update`, `posts:*`
- `user`: `profile:read,update`, `posts:create,read`

## API Endpoints

### Authentication Required
All RBAC endpoints require JWT authentication via `Authorization: Bearer <token>` header.

### Role Management (Admin Only)

**Create Role**
```
POST /api/v1/rbac/roles
{
    "name": "editor",
    "display_name": "Content Editor",
    "description": "Can create and edit content",
    "is_active": true
}
```

**Get All Roles**
```
GET /api/v1/rbac/roles
```

**Get Role by ID**
```
GET /api/v1/rbac/roles/{role_id}
```

**Update Role**
```
PUT /api/v1/rbac/roles/{role_id}
{
    "display_name": "Senior Editor",
    "description": "Senior content editor with additional permissions",
    "is_active": true
}
```

**Delete Role**
```
DELETE /api/v1/rbac/roles/{role}
```

**Get Users with Role**
```
GET /api/v1/rbac/roles/{role}/users
```

**Get Role Permissions**
```
GET /api/v1/rbac/roles/{role}/permissions
```

### User Role Management (Admin Only)

**Assign Role to User**
```
POST /api/v1/rbac/users/assign-role
{
    "user_id": 1,
    "role": "editor"
}
```

**Remove Role from User**
```
POST /api/v1/rbac/users/remove-role
{
    "user_id": 1,
    "role": "editor"
}
```

**Get User Roles**
```
GET /api/v1/rbac/users/{user_id}/roles
```

### Permission Management (Admin Only)

**Add Permission to Role**
```
POST /api/v1/rbac/permissions
{
    "role": "editor",
    "resource": "posts",
    "action": "update"
}
```

**Remove Permission from Role**
```
DELETE /api/v1/rbac/permissions
{
    "role": "editor",
    "resource": "posts",
    "action": "update"
}
```

**Check User Permission**
```
GET /api/v1/rbac/users/{user_id}/check-permission?resource=posts&action=read
```

## Middleware Usage

### Require Specific Permission
```go
// Require specific permission for an endpoint
apiV1.GET("/posts", handler.GetPosts, middleware.RequirePermission(rbacService, "posts", "read"))
```

### Require Specific Role
```go
// Require specific role
apiV1.GET("/admin", handler.AdminPanel, middleware.RequireRole(rbacService, "admin"))
```

### Require Any of Multiple Roles
```go
// Require any of the specified roles
apiV1.GET("/moderate", handler.Moderate, middleware.RequireAnyRole(rbacService, "admin", "moderator"))
```

### Require All Specified Roles
```go
// Require all specified roles
apiV1.GET("/special", handler.Special, middleware.RequireAllRoles(rbacService, "admin", "editor"))
```

## Code Usage Examples

### Initialize RBAC Service
```go
rbacService, err := rbac.NewRBACService(db)
if err != nil {
    log.Fatal("Failed to initialize RBAC service:", err)
}
```

### Check Permission Programmatically
```go
allowed, err := rbacService.CheckPermission(userID, "posts", "read")
if err != nil {
    return err
}
if !allowed {
    return echo.NewHTTPError(http.StatusForbidden, "Insufficient permissions")
}
```

### Assign Role to User
```go
err := rbacService.AssignRoleToUser(userID, "editor")
if err != nil {
    return err
}
```

### Add Permission to Role
```go
err := rbacService.AddPermission("editor", "posts", "update")
if err != nil {
    return err
}
```

## Database Schema

The system uses Casbin's default table structure:
- `casbin_rule`: Stores all policies and role mappings
  - `ptype`: Policy type (p for permission, g for grouping/role)
  - `v0, v1, v2`: Values for subject, object, action

## Common Use Cases

### 1. Create a New Role with Permissions
```bash
# Create role
curl -X POST http://localhost:8080/api/v1/rbac/roles \
  -H "Authorization: Bearer <admin_token>" \
  -H "Content-Type: application/json" \
  -d '{"role": "editor"}'

# Add permissions
curl -X POST http://localhost:8080/api/v1/rbac/permissions \
  -H "Authorization: Bearer <admin_token>" \
  -H "Content-Type: application/json" \
  -d '{"role": "editor", "resource": "posts", "action": "update"}'
```

### 2. Assign Role to User
```bash
curl -X POST http://localhost:8080/api/v1/rbac/users/assign-role \
  -H "Authorization: Bearer <admin_token>" \
  -H "Content-Type: application/json" \
  -d '{"user_id": 1, "role": "editor"}'
```

### 3. Check User Permissions
```bash
curl -X GET "http://localhost:8080/api/v1/rbac/users/1/check-permission?resource=posts&action=read" \
  -H "Authorization: Bearer <admin_token>"
```

## Security Considerations

1. **Admin Access**: Only users with `admin` role can access RBAC management endpoints
2. **JWT Required**: All endpoints require valid JWT authentication
3. **Permission Checks**: Use middleware to protect sensitive endpoints
4. **Resource Naming**: Use consistent resource names across your application
5. **Action Naming**: Use standard actions like `read`, `write`, `update`, `delete`

## Error Handling

The system provides proper error responses:
- `400 Bad Request`: Invalid request parameters
- `401 Unauthorized`: Missing or invalid JWT token
- `403 Forbidden`: Insufficient permissions
- `500 Internal Server Error`: Database or system errors