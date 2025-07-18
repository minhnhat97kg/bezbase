# Multi-Tenancy and Advanced RBAC Implementation

This document describes the implementation of Multi-tenancy and Advanced RBAC features in BezBase, representing Priority 5 features from the roadmap (v2.0+).

## 🏗️ Overview

BezBase now supports enterprise-grade multi-tenancy with organization isolation and advanced hierarchical role-based access control (RBAC). These features enable organizations to manage multiple tenants with sophisticated permission systems.

## 🔧 Core Features Implemented

### 1. Multi-Tenancy (Organization Isolation)

**Database Schema:**
- `organizations` - Organization entities with settings and plan types
- `organization_users` - Many-to-many relationship between users and organizations
- `organization_invitations` - Invitation system for joining organizations
- Added `current_org_id` to users for context switching
- Added `org_id` to roles for organization-specific roles

**Key Capabilities:**
- ✅ Organization creation and management
- ✅ User membership across multiple organizations
- ✅ Organization-specific role assignment
- ✅ Context switching between organizations
- ✅ Invitation system with email notifications
- ✅ Organization isolation and data separation

### 2. Advanced RBAC (Hierarchical Roles)

**Database Schema:**
- `role_inheritances` - Efficient role hierarchy queries
- `contextual_permissions` - Context-aware permission system
- `role_templates` - Predefined role configurations
- Added `parent_role_id` and `hierarchy_level` to roles

**Key Capabilities:**
- ✅ Hierarchical role inheritance
- ✅ Context-aware permissions (organization, project, resource-specific)
- ✅ Role templates for quick role creation
- ✅ Permission inheritance through role hierarchy
- ✅ Organization-scoped permission checking

## 📊 Database Architecture

### Multi-Tenancy Tables

```sql
-- Organizations table
CREATE TABLE organizations (
    id SERIAL PRIMARY KEY,
    name VARCHAR(255) NOT NULL,
    slug VARCHAR(100) UNIQUE NOT NULL,
    domain VARCHAR(255),
    settings JSONB,
    is_active BOOLEAN DEFAULT TRUE,
    plan_type VARCHAR(50) DEFAULT 'free',
    created_at TIMESTAMP DEFAULT CURRENT_TIMESTAMP,
    updated_at TIMESTAMP DEFAULT CURRENT_TIMESTAMP,
    deleted_at TIMESTAMP
);

-- Organization Users (many-to-many)
CREATE TABLE organization_users (
    id SERIAL PRIMARY KEY,
    org_id INTEGER NOT NULL REFERENCES organizations(id) ON DELETE CASCADE,
    user_id INTEGER NOT NULL REFERENCES users(id) ON DELETE CASCADE,
    role VARCHAR(50) NOT NULL DEFAULT 'member',
    is_primary BOOLEAN DEFAULT FALSE,
    joined_at TIMESTAMP DEFAULT CURRENT_TIMESTAMP,
    created_at TIMESTAMP DEFAULT CURRENT_TIMESTAMP,
    updated_at TIMESTAMP DEFAULT CURRENT_TIMESTAMP,
    deleted_at TIMESTAMP,
    UNIQUE(org_id, user_id) WHERE deleted_at IS NULL
);

-- Organization Invitations
CREATE TABLE organization_invitations (
    id SERIAL PRIMARY KEY,
    org_id INTEGER NOT NULL REFERENCES organizations(id) ON DELETE CASCADE,
    email VARCHAR(255) NOT NULL,
    role VARCHAR(50) NOT NULL DEFAULT 'member',
    token VARCHAR(255) UNIQUE NOT NULL,
    expires_at TIMESTAMP NOT NULL,
    invited_by INTEGER NOT NULL REFERENCES users(id) ON DELETE CASCADE,
    accepted_at TIMESTAMP,
    created_at TIMESTAMP DEFAULT CURRENT_TIMESTAMP,
    updated_at TIMESTAMP DEFAULT CURRENT_TIMESTAMP,
    deleted_at TIMESTAMP
);
```

### Advanced RBAC Tables

```sql
-- Role Inheritance for hierarchy
CREATE TABLE role_inheritances (
    id SERIAL PRIMARY KEY,
    parent_role_id INTEGER NOT NULL REFERENCES roles(id) ON DELETE CASCADE,
    child_role_id INTEGER NOT NULL REFERENCES roles(id) ON DELETE CASCADE,
    depth INTEGER NOT NULL DEFAULT 1,
    created_at TIMESTAMP DEFAULT CURRENT_TIMESTAMP,
    updated_at TIMESTAMP DEFAULT CURRENT_TIMESTAMP,
    UNIQUE(parent_role_id, child_role_id)
);

-- Contextual Permissions
CREATE TABLE contextual_permissions (
    id SERIAL PRIMARY KEY,
    role_id INTEGER NOT NULL REFERENCES roles(id) ON DELETE CASCADE,
    resource VARCHAR(100) NOT NULL,
    action VARCHAR(50) NOT NULL,
    context_type VARCHAR(50), -- e.g., 'organization', 'project'
    context_value VARCHAR(255), -- specific ID or pattern
    is_granted BOOLEAN DEFAULT TRUE,
    created_at TIMESTAMP DEFAULT CURRENT_TIMESTAMP,
    updated_at TIMESTAMP DEFAULT CURRENT_TIMESTAMP,
    deleted_at TIMESTAMP,
    UNIQUE(role_id, resource, action, context_type, context_value) WHERE deleted_at IS NULL
);

-- Role Templates
CREATE TABLE role_templates (
    id SERIAL PRIMARY KEY,
    name VARCHAR(100) NOT NULL,
    display_name VARCHAR(255) NOT NULL,
    description VARCHAR(500),
    category VARCHAR(100), -- 'system', 'business', 'department', 'basic'
    config JSONB, -- JSON configuration for permissions
    is_active BOOLEAN DEFAULT TRUE,
    created_at TIMESTAMP DEFAULT CURRENT_TIMESTAMP,
    updated_at TIMESTAMP DEFAULT CURRENT_TIMESTAMP,
    deleted_at TIMESTAMP,
    UNIQUE(name) WHERE deleted_at IS NULL
);
```

## 🚀 API Endpoints

### Organization Management

```http
# Create Organization
POST /api/v1/organizations
{
  "name": "Acme Corporation",
  "slug": "acme-corp",
  "domain": "acme.com",
  "plan_type": "premium"
}

# Get Organization
GET /api/v1/organizations/{id}

# Update Organization
PUT /api/v1/organizations/{id}
{
  "name": "Updated Name",
  "plan_type": "enterprise"
}

# Get User Organizations
GET /api/v1/user/organizations

# Switch Organization Context
POST /api/v1/organizations/{id}/switch

# Invite User to Organization
POST /api/v1/organizations/{id}/invite
{
  "email": "user@example.com",
  "role": "member"
}

# Accept Invitation
POST /api/v1/organizations/invitations/{token}/accept

# Remove User from Organization
DELETE /api/v1/organizations/{id}/users/{userId}

# Update User Role in Organization
PUT /api/v1/organizations/{id}/users/{userId}/role
{
  "role": "admin"
}
```

### Advanced RBAC

```http
# Create Role from Template
POST /api/v1/rbac/roles/from-template
{
  "template_id": 1,
  "org_id": 123,
  "custom_name": "custom_role_name"
}

# Set Role Hierarchy
PUT /api/v1/rbac/roles/{id}/parent
{
  "parent_role_id": 456
}

# Get Roles by Organization
GET /api/v1/rbac/roles?org_id=123

# Create Contextual Permission
POST /api/v1/rbac/contextual-permissions
{
  "role_id": 789,
  "resource": "projects",
  "action": "read",
  "context_type": "organization",
  "context_value": "123"
}

# Get Effective Permissions
GET /api/v1/rbac/users/{userId}/effective-permissions?org_id=123
```

## 🛠️ Implementation Details

### Service Layer

**OrganizationService** (`internal/services/organization.go`):
- Organization CRUD operations
- User membership management
- Invitation system
- Context switching
- Permission validation

**Enhanced RBACService** (`internal/services/rbac.go`):
- Hierarchical permission checking
- Context-aware authorization
- Role template management
- Organization-scoped role queries

### Middleware

**TenantMiddleware** (`internal/middleware/tenant.go`):
- Organization context extraction from headers/params
- User organization validation
- Context injection for services

**OrganizationScope** (`internal/middleware/tenant.go`):
- Ensures user has access to requested organization
- Validates organization membership

### Models

**Organization Models** (`internal/models/organization.go`):
- Organization entity with validation
- OrganizationUser with role management
- OrganizationInvitation with expiration logic

**Advanced RBAC Models** (`internal/models/advanced_rbac.go`):
- RoleInheritance for hierarchy queries
- ContextualPermission for fine-grained access
- RoleTemplate for role standardization

## 🔐 Permission System

### Permission Hierarchy

1. **Global Permissions** - System-wide access (admin roles)
2. **Organization Permissions** - Organization-scoped access
3. **Resource Permissions** - Specific resource access
4. **Inherited Permissions** - From parent roles in hierarchy

### Permission Resolution Flow

```go
func CheckPermissionWithContext(userID uint, resource, action string, orgID *uint) bool {
    // 1. Get user's roles in organization context
    userRoles := GetUserRolesInOrganization(userID, orgID)
    
    for _, role := range userRoles {
        // 2. Check direct role permissions (Casbin)
        if enforcer.Enforce(role.Name, resource, action) {
            return true
        }
        
        // 3. Check contextual permissions
        if checkContextualPermissions(role.ID, resource, action, orgID) {
            return true
        }
        
        // 4. Check inherited permissions from parent roles
        if checkInheritedPermissions(role.ID, resource, action, orgID) {
            return true
        }
    }
    
    return false
}
```

### Context Types

- **organization** - Organization-scoped permissions
- **project** - Project-scoped permissions  
- **department** - Department-scoped permissions
- **resource** - Specific resource permissions

## 🎯 Usage Examples

### Organization Management

```go
// Create organization
orgReq := services.CreateOrganizationRequest{
    Name:     "My Company",
    Slug:     "my-company",
    PlanType: "premium",
}
org, err := orgService.CreateOrganization(userID, orgReq)

// Invite user
inviteReq := services.InviteUserRequest{
    Email: "colleague@company.com",
    Role:  "admin",
}
invitation, err := orgService.InviteUser(orgID, inviterID, inviteReq)

// Switch organization context
err := orgService.SwitchOrganization(userID, newOrgID)
```

### Advanced RBAC

```go
// Create hierarchical role
role := &models.Role{
    Name:           "team_lead",
    DisplayName:    "Team Lead",
    OrgID:          &orgID,
    ParentRoleID:   &managerRoleID,
    HierarchyLevel: 2,
}
err := rbacService.CreateRole(role)

// Check organization-scoped permission
hasPermission, err := rbacService.CheckPermissionWithContext(
    userID, "projects", "create", &orgID,
)

// Create contextual permission
permission := &models.ContextualPermission{
    RoleID:       roleID,
    Resource:     "users",
    Action:       "read",
    ContextType:  "organization",
    ContextValue: fmt.Sprintf("%d", orgID),
    IsGranted:    true,
}
err := contextualPermRepo.Create(permission)
```

## 🧪 Testing

### Unit Tests

Run organization and RBAC tests:
```bash
go test ./internal/services/organization_test.go
go test ./internal/services/rbac_test.go
go test ./internal/middleware/tenant_test.go
```

### Integration Tests

Test multi-tenant scenarios:
```bash
go test ./internal/handlers/organization_test.go
go test ./internal/handlers/rbac_test.go
```

### Migration Tests

Verify database migrations:
```bash
go run cmd/migrate/main.go up
go run cmd/migrate/main.go down
```

## 🔧 Configuration

### Environment Variables

```env
# Multi-tenancy settings
DEFAULT_ORG_PLAN=free
MAX_ORGS_PER_USER=10
INVITATION_EXPIRY_HOURS=168

# RBAC settings
MAX_ROLE_HIERARCHY_DEPTH=10
ENABLE_CONTEXTUAL_PERMISSIONS=true
```

### Database Indexes

Ensure these indexes exist for optimal performance:

```sql
-- Organization indexes
CREATE INDEX idx_organizations_slug ON organizations(slug);
CREATE INDEX idx_organizations_domain ON organizations(domain);
CREATE INDEX idx_organization_users_org_id ON organization_users(org_id);
CREATE INDEX idx_organization_users_user_id ON organization_users(user_id);

-- RBAC indexes
CREATE INDEX idx_roles_org_id ON roles(org_id);
CREATE INDEX idx_roles_hierarchy_level ON roles(hierarchy_level);
CREATE INDEX idx_contextual_permissions_role_id ON contextual_permissions(role_id);
CREATE INDEX idx_contextual_permissions_context ON contextual_permissions(context_type, context_value);
```

## 🚦 Migration Guide

### From Single-Tenant to Multi-Tenant

1. **Run Migrations**:
   ```bash
   go run cmd/migrate/main.go up
   ```

2. **Migrate Existing Users**:
   - All existing users are automatically added to a "Default Organization"
   - Existing roles become global roles (org_id = NULL)

3. **Update API Calls**:
   - Add organization context headers: `X-Organization-ID: 123`
   - Use organization-scoped endpoints where applicable

4. **Frontend Updates**:
   - Add organization selector component
   - Update permission checks to include organization context
   - Handle organization switching in navigation

### Role Migration

Existing roles are preserved and enhanced:
- System roles (admin, moderator, user) remain global
- New organization-specific roles inherit from global roles
- Permission hierarchy is automatically established

## 🔍 Monitoring & Analytics

### Key Metrics

- **Organizations**: Total count, active organizations, plan distribution
- **Users**: Cross-organization membership, role distribution
- **Invitations**: Success rate, expiration rate, response time
- **Permissions**: Permission check frequency, denial rate by resource

### Logging

Enhanced logging includes:
- Organization context in all operations
- Permission check results with hierarchy trace
- Role inheritance resolution paths
- Multi-tenant data access patterns

## 🛡️ Security Considerations

### Data Isolation

- **Row-Level Security**: Organization-scoped queries prevent data leakage
- **API Context Validation**: All endpoints validate organization membership
- **Role Scoping**: Organization-specific roles cannot access other organizations

### Permission Validation

- **Multi-Level Checks**: Global → Organization → Resource → Inherited permissions
- **Context Verification**: All permissions validated within correct organizational context
- **Audit Trail**: Complete permission resolution logging for compliance

### Invitation Security

- **Token Expiration**: Configurable invitation expiry (default: 7 days)
- **Email Verification**: Only verified email addresses can accept invitations
- **Role Restrictions**: Inviter cannot assign roles higher than their own

## 📈 Performance Optimizations

### Database Optimizations

- **Efficient Hierarchy Queries**: Recursive CTEs for role inheritance
- **Contextual Permission Indexing**: Optimized for organization-scoped lookups
- **Membership Caching**: Organization membership cached for frequent checks

### Service Optimizations

- **Permission Caching**: Role permissions cached with organization context
- **Batch Operations**: Bulk role assignments and permission updates
- **Lazy Loading**: On-demand loading of organization relationships

## 🔮 Future Enhancements

### Planned Features (v2.1+)

- **Custom Permission Types**: User-defined permission resources and actions
- **Dynamic Role Creation**: UI-based role builder with permission templates
- **Organization Analytics**: Usage metrics and activity dashboards
- **Advanced Auditing**: Detailed permission usage and change tracking
- **API Rate Limiting**: Per-organization rate limiting and quotas

### Integration Opportunities

- **Single Sign-On (SSO)**: SAML/OAuth2 integration for organization authentication
- **External Directory Sync**: LDAP/Active Directory synchronization
- **Webhook Events**: Organization and permission change notifications
- **Third-Party Integrations**: Slack, Microsoft Teams, Google Workspace

---

## 🏁 Conclusion

The Multi-Tenancy and Advanced RBAC implementation provides BezBase with enterprise-grade capabilities for managing complex organizational structures and sophisticated permission systems. This foundation supports scalable SaaS offerings with proper data isolation, flexible role management, and comprehensive access control.

For detailed API documentation, see the Swagger/OpenAPI specifications at `/docs` when running the application.

For implementation questions or feature requests, please refer to the main project documentation or create an issue in the repository.