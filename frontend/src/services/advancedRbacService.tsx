import axios, { AxiosResponse } from 'axios';

const API_BASE_URL = process.env.REACT_APP_API_URL || 'http://localhost:8080/api';
const API_VERSION = 'v1';

// Type definitions for advanced RBAC
export interface Role {
  id: number;
  name: string;
  display_name: string;
  description: string;
  is_system: boolean;
  is_active: boolean;
  org_id?: number;
  parent_role_id?: number;
  hierarchy_level: number;
  organization?: Organization;
  parent_role?: Role;
  child_roles?: Role[];
  contextual_permissions?: ContextualPermission[];
}

export interface Organization {
  id: number;
  name: string;
  slug: string;
  domain?: string;
  plan_type: string;
}

export interface ContextualPermission {
  id: number;
  role_id: number;
  resource: string;
  action: string;
  context_type?: string;
  context_value?: string;
  is_granted: boolean;
  created_at: string;
  updated_at: string;
  role?: Role;
}

export interface RoleTemplate {
  id: number;
  name: string;
  display_name: string;
  description: string;
  category: 'system' | 'business' | 'department' | 'basic';
  config: string; // JSON configuration
  is_active: boolean;
  created_at: string;
  updated_at: string;
}

export interface CreateRoleFromTemplateRequest {
  template_id: number;
  custom_name?: string;
}

export interface SetRoleParentRequest {
  parent_role_id?: number;
}

export interface CreateContextualPermissionRequest {
  role_id: number;
  resource: string;
  action: string;
  context_type?: string;
  context_value?: string;
  is_granted?: boolean;
}

export interface RoleHierarchyResponse {
  role_id: number;
  parent_roles: Role[];
  child_roles: Role[];
}

export interface ApiResponse<T = any> {
  data?: T;
  message?: string;
  success?: boolean;
}

export interface ErrorResponse {
  message: string;
  details?: string;
  code?: string;
}

// Create API instance
const createApiInstance = () => {
  const instance = axios.create({
    baseURL: `${API_BASE_URL}/${API_VERSION}`,
    headers: {
      'Content-Type': 'application/json',
    },
  });

  // Add auth token and organization context to requests
  instance.interceptors.request.use(
    (config) => {
      const token = localStorage.getItem('token');
      if (token) {
        config.headers.Authorization = `Bearer ${token}`;
      }
      
      // Add organization context if available
      const currentOrgId = localStorage.getItem('currentOrgId');
      if (currentOrgId) {
        config.headers['X-Organization-ID'] = currentOrgId;
      }

      return config;
    },
    (error) => {
      return Promise.reject(error);
    }
  );

  // Handle responses and errors
  instance.interceptors.response.use(
    (response) => response,
    (error) => {
      if (error.response?.status === 401) {
        localStorage.removeItem('token');
        localStorage.removeItem('user');
        localStorage.removeItem('currentOrgId');
        window.location.href = '/login';
      }
      return Promise.reject(error);
    }
  );

  return instance;
};

const api = createApiInstance();

// Advanced RBAC Service
export const advancedRbacService = {
  // Role management with hierarchy
  createRoleFromTemplate: (data: CreateRoleFromTemplateRequest): Promise<AxiosResponse<Role>> => {
    return api.post('/rbac/roles/from-template', data);
  },

  setRoleParent: (roleId: number, data: SetRoleParentRequest): Promise<AxiosResponse<ApiResponse>> => {
    return api.put(`/rbac/roles/${roleId}/parent`, data);
  },

  getRolesByOrganization: (orgId?: number): Promise<AxiosResponse<Role[]>> => {
    const params = orgId ? { org_id: orgId } : {};
    return api.get('/rbac/roles', { params });
  },

  getRoleHierarchy: (roleId: number): Promise<AxiosResponse<RoleHierarchyResponse>> => {
    return api.get(`/rbac/roles/${roleId}/hierarchy`);
  },

  // Role templates
  getRoleTemplates: (category?: string): Promise<AxiosResponse<RoleTemplate[]>> => {
    const params = category ? { category } : {};
    return api.get('/rbac/role-templates', { params });
  },

  // Contextual permissions
  createContextualPermission: (data: CreateContextualPermissionRequest): Promise<AxiosResponse<ContextualPermission>> => {
    return api.post('/rbac/contextual-permissions', data);
  },

  getEffectivePermissions: (userId: number, orgId?: number): Promise<AxiosResponse<ContextualPermission[]>> => {
    const params = orgId ? { org_id: orgId } : {};
    return api.get(`/rbac/users/${userId}/effective-permissions`, { params });
  },

  // Utility functions
  formatRoleHierarchy: (roles: Role[]): Role[] => {
    // Ensure roles is an array before sorting
    if (!Array.isArray(roles)) {
      console.warn('formatRoleHierarchy: roles is not an array, returning empty array');
      return [];
    }
    
    // Sort roles by hierarchy level and display name
    return roles.sort((a, b) => {
      if (a.hierarchy_level !== b.hierarchy_level) {
        return a.hierarchy_level - b.hierarchy_level;
      }
      return a.display_name.localeCompare(b.display_name);
    });
  },

  getRoleColor: (role: Role): string => {
    if (role.is_system) {
      return 'bg-red-100 text-red-800';
    }
    
    switch (role.hierarchy_level) {
      case 0:
        return 'bg-purple-100 text-purple-800'; // Top level
      case 1:
        return 'bg-blue-100 text-blue-800'; // Organization level
      case 2:
        return 'bg-green-100 text-green-800'; // Department level
      default:
        return 'bg-gray-100 text-gray-800'; // Lower levels
    }
  },

  getRoleIcon: (role: Role): string => {
    if (role.is_system) return '🔒';
    
    switch (role.hierarchy_level) {
      case 0: return '👑'; // Global admin
      case 1: return '🏢'; // Org admin
      case 2: return '👥'; // Team lead
      default: return '👤'; // Regular role
    }
  },

  getPermissionColor: (permission: ContextualPermission): string => {
    if (!permission.is_granted) {
      return 'bg-red-100 text-red-800';
    }
    
    if (permission.context_type) {
      switch (permission.context_type) {
        case 'organization':
          return 'bg-blue-100 text-blue-800';
        case 'project':
          return 'bg-green-100 text-green-800';
        case 'department':
          return 'bg-purple-100 text-purple-800';
        default:
          return 'bg-gray-100 text-gray-800';
      }
    }
    
    return 'bg-indigo-100 text-indigo-800'; // Global permission
  },

  formatPermissionKey: (permission: ContextualPermission): string => {
    let key = `${permission.resource}:${permission.action}`;
    if (permission.context_type && permission.context_value) {
      key += `:${permission.context_type}:${permission.context_value}`;
    }
    return key;
  },

  buildRoleTree: (roles: Role[]): Role[] => {
    if (!Array.isArray(roles)) {
      console.warn('buildRoleTree: roles is not an array, returning empty array');
      return [];
    }
    
    const roleMap = new Map<number, Role>();
    const rootRoles: Role[] = [];

    // Create a map of all roles
    roles.forEach(role => {
      roleMap.set(role.id, { ...role, child_roles: [] });
    });

    // Build the tree structure
    roles.forEach(role => {
      const roleWithChildren = roleMap.get(role.id)!;
      
      if (role.parent_role_id) {
        const parent = roleMap.get(role.parent_role_id);
        if (parent) {
          parent.child_roles = parent.child_roles || [];
          parent.child_roles.push(roleWithChildren);
        } else {
          // Parent not found, treat as root
          rootRoles.push(roleWithChildren);
        }
      } else {
        rootRoles.push(roleWithChildren);
      }
    });

    return rootRoles;
  },

  flattenRoleTree: (roles: Role[], level: number = 0): Array<Role & { indentLevel: number }> => {
    if (!Array.isArray(roles)) {
      console.warn('flattenRoleTree: roles is not an array, returning empty array');
      return [];
    }
    
    const flattened: Array<Role & { indentLevel: number }> = [];

    roles.forEach(role => {
      flattened.push({ ...role, indentLevel: level });
      
      if (role.child_roles && role.child_roles.length > 0) {
        const childrenFlattened = advancedRbacService.flattenRoleTree(role.child_roles, level + 1);
        flattened.push(...childrenFlattened);
      }
    });

    return flattened;
  },

  getAvailableParentRoles: (roles: Role[], currentRoleId: number): Role[] => {
    if (!Array.isArray(roles)) {
      console.warn('getAvailableParentRoles: roles is not an array, returning empty array');
      return [];
    }
    
    // Filter out the current role and its descendants to prevent circular dependencies
    const isDescendant = (role: Role, ancestorId: number): boolean => {
      if (role.id === ancestorId) return true;
      if (role.child_roles && Array.isArray(role.child_roles)) {
        return role.child_roles.some(child => isDescendant(child, ancestorId));
      }
      return false;
    };

    return roles.filter(role => 
      role.id !== currentRoleId && 
      !isDescendant(role, currentRoleId)
    );
  },

  // Permission utilities
  canManageRole: (role: Role, userRoles: Role[]): boolean => {
    // System roles can only be managed by system admins
    if (role.is_system) {
      return userRoles.some(userRole => userRole.is_system && userRole.name === 'admin');
    }
    
    // Organization roles can be managed by organization admins
    if (role.org_id) {
      return userRoles.some(userRole => 
        (userRole.org_id === role.org_id && userRole.hierarchy_level <= role.hierarchy_level) ||
        (userRole.is_system && userRole.name === 'admin')
      );
    }
    
    // Global roles require system admin
    return userRoles.some(userRole => userRole.is_system && userRole.name === 'admin');
  },

  canCreateRole: (orgId: number | null, userRoles: Role[]): boolean => {
    if (orgId) {
      // Creating organization role
      return userRoles.some(userRole => 
        (userRole.org_id === orgId && (userRole.name.includes('admin') || userRole.hierarchy_level <= 1)) ||
        (userRole.is_system && userRole.name === 'admin')
      );
    }
    
    // Creating global role requires system admin
    return userRoles.some(userRole => userRole.is_system && userRole.name === 'admin');
  },
};

export default advancedRbacService;