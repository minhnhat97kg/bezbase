import axios, { AxiosResponse } from 'axios';

const API_BASE_URL = process.env.REACT_APP_API_URL || 'http://localhost:8080/api';
const API_VERSION = 'v1';

// Unified type definitions for RBAC
export interface Role {
  id: number;
  name: string;
  display_name: string;
  description: string;
  is_system: boolean;
  is_active: boolean;
  parent_role_id?: number;
  hierarchy_level: number;
  parent_role?: Role;
  child_roles?: Role[];
  user_count?: number;
}

export interface Permission {
  resource: string;
  action: string;
  permission: string;
}

export interface UserRoleAssignment {
  user_id: number;
  role: string;
}

export interface RoleQueryParams {
  page?: number;
  page_size?: number;
  search?: string;
  status?: string;
  is_system?: boolean;
  sort?: string;
  order?: 'asc' | 'desc';
}

// Create unified API instance
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

      const orgId = localStorage.getItem('organization_id');
      if (orgId) {
        config.headers['X-Organization-ID'] = orgId;
      }

      const currentLanguage = localStorage.getItem('i18nextLng') || 'en';
      config.headers['Accept-Language'] = currentLanguage;
      config.headers['Content-Language'] = currentLanguage;

      return config;
    },
    (error) => {
      return Promise.reject(error);
    }
  );

  // Add response interceptor for error handling
  instance.interceptors.response.use(
    (response) => response,
    (error) => {
      if (error.response?.status === 401) {
        // Token expired or invalid
        localStorage.removeItem('token');
        localStorage.removeItem('user');
        window.location.href = '/login';
      }
      return Promise.reject(error);
    }
  );

  return instance;
};

const api = createApiInstance();

// Unified RBAC Service
export const rbacService = {
  // Role Management
  getRoles: (params: RoleQueryParams = {}): Promise<AxiosResponse<Role[]>> => {
    const queryParams = new URLSearchParams();
    if (params.page) queryParams.append('page', params.page.toString());
    if (params.page_size) queryParams.append('page_size', params.page_size.toString());
    if (params.search) queryParams.append('search', params.search);
    if (params.status) queryParams.append('status', params.status);
    if (params.is_system !== undefined) queryParams.append('is_system', params.is_system.toString());
    if (params.sort) queryParams.append('sort', params.sort);
    if (params.order) queryParams.append('order', params.order);
    
    const url = queryParams.toString() ? `/rbac/roles?${queryParams.toString()}` : '/rbac/roles';
    return api.get(url);
  },

  getRole: (roleId: number): Promise<AxiosResponse<Role>> => {
    return api.get(`/rbac/roles/${roleId}`);
  },

  createRole: (data: { name: string; display_name: string; description: string }): Promise<AxiosResponse<Role>> => {
    return api.post('/rbac/roles', data);
  },

  updateRole: (roleId: number, data: { name?: string; display_name?: string; description?: string }): Promise<AxiosResponse<Role>> => {
    return api.put(`/rbac/roles/${roleId}`, data);
  },

  deleteRole: (roleId: number): Promise<AxiosResponse<void>> => {
    return api.delete(`/rbac/roles/${roleId}`);
  },

  // User Role Assignment
  assignRole: (data: UserRoleAssignment): Promise<AxiosResponse<any>> => {
    return api.post('/rbac/users/assign-role', data);
  },

  removeRole: (data: UserRoleAssignment): Promise<AxiosResponse<any>> => {
    return api.post('/rbac/users/remove-role', data);
  },

  getUserRoles: (userId: number): Promise<AxiosResponse<string[]>> => {
    return api.get(`/rbac/users/${userId}/roles`);
  },

  // Permission Management
  getPermissions: (params: any = {}): Promise<AxiosResponse<any>> => {
    const queryParams = new URLSearchParams();
    Object.keys(params).forEach(key => {
      if (params[key] !== undefined && params[key] !== null && params[key] !== '') {
        queryParams.append(key, params[key].toString());
      }
    });
    
    const url = queryParams.toString() ? `/rbac/permissions?${queryParams.toString()}` : '/rbac/permissions';
    return api.get(url);
  },

  addPermission: (permissionData: any): Promise<AxiosResponse<any>> => {
    return api.post('/rbac/permissions', permissionData);
  },

  removePermission: (permissionData: any): Promise<AxiosResponse<any>> => {
    return api.delete('/rbac/permissions', { data: permissionData });
  },

  getAvailablePermissions: (): Promise<AxiosResponse<Permission[]>> => {
    return api.get('/rbac/permissions/available');
  },

  // User Permissions (for permission checking)
  getUserPermissions: (userId: number): Promise<AxiosResponse<string[]>> => {
    return api.get(`/rbac/users/${userId}/permissions`);
  },

  checkUserPermission: (userId: number, resource: string, action: string): Promise<AxiosResponse<{ allowed: boolean }>> => {
    return api.post(`/rbac/users/${userId}/check-permission`, { resource, action });
  },

  // Get current user's permissions (for AuthContext)
  getMyPermissions: (): Promise<AxiosResponse<{user_id: number, permissions: string[]}>> => {
    return api.get('/rbac/me/permissions');
  },

  // Role Inheritance Management
  setRoleParent: (roleId: number, parentRoleId: number | null): Promise<AxiosResponse<any>> => {
    return api.put(`/rbac/roles/${roleId}/parent`, { parent_role_id: parentRoleId });
  },

  getRoleHierarchy: (roleId: number): Promise<AxiosResponse<{
    role_id: number;
    parent_roles: Role[];
    child_roles: Role[];
  }>> => {
    return api.get(`/rbac/roles/${roleId}/hierarchy`);
  },

  getEffectivePermissions: (userId: number, orgId?: number): Promise<AxiosResponse<any[]>> => {
    const params = orgId ? `?org_id=${orgId}` : '';
    return api.get(`/rbac/users/${userId}/effective-permissions${params}`);
  },

  // Get roles that can be set as parent (prevents circular references)
  getEligibleParentRoles: (roleId: number): Promise<AxiosResponse<Role[]>> => {
    return api.get(`/rbac/roles/${roleId}/eligible-parents`);
  }
};

export default rbacService;