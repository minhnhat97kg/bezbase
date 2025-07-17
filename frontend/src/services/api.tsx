import axios, { AxiosResponse } from 'axios';
import i18n from '../i18n';

const API_BASE_URL = process.env.REACT_APP_API_URL || 'http://localhost:8080/api';

// API Configuration - hardcoded to v1
const API_VERSION = 'v1';

// Type definitions
interface LoginRequest {
  username: string;
  password: string;
}

interface RegisterRequest {
  username: string;
  password: string;
  email: string;
  [key: string]: any;
}

interface UserData {
  id?: string;
  username?: string;
  email: string;
  [key: string]: any;
}

interface PasswordData {
  current_password: string;
  new_password: string;
  confirm_password?: string;
}

interface RoleData {
  name: string;
  description?: string;
  permissions?: string[];
  [key: string]: any;
}

interface PermissionData {
  resource: string;
  action: string;
  [key: string]: any;
}

interface RoleQueryParams {
  page?: number;
  page_size?: number;
  search?: string;
  status?: string;
  is_system?: boolean;
  sort?: string;
  order?: string;
}

interface AssignmentData {
  userId: string;
  roleName: string;
  [key: string]: any;
}

// API Response Types
interface ApiResponse<T = any> {
  data: T;
  success: boolean;
  message?: string;
  version?: string;
  timestamp?: string;
}

interface ApiError {
  error: string;
  message: string;
  status_code: number;
  version?: string;
  timestamp?: string;
}

interface PaginatedResponse<T = any> {
  data: T[];
  pagination: {
    page: number;
    page_size: number;
    total_items: number;
    total_pages: number;
    has_next: boolean;
    has_prev: boolean;
  };
  success: boolean;
  message?: string;
  version?: string;
}

// Version-specific response types
interface VersionedApiResponse<T = any> extends ApiResponse<T> {
  api_version: string;
  supported_versions: string[];
  deprecation_notice?: string;
}

// Export types for use in components
export type {
  LoginRequest,
  RegisterRequest,
  UserData,
  PasswordData,
  RoleData,
  PermissionData,
  RoleQueryParams,
  AssignmentData,
  ApiResponse,
  ApiError,
  PaginatedResponse,
  VersionedApiResponse
};

// Create API instance with language support and hardcoded v1 versioning
const createApiInstance = (useVersioning: boolean = true) => {
  const baseURL = useVersioning ? `${API_BASE_URL}/${API_VERSION}` : API_BASE_URL;

  const instance = axios.create({
    baseURL: baseURL,
    headers: {
      'Content-Type': 'application/json',
    },
  });

  // Add language headers to all requests
  instance.interceptors.request.use(
    (config) => {
      // Add language headers
      const currentLanguage = i18n.language || 'en';
      config.headers['Accept-Language'] = currentLanguage;
      config.headers['Content-Language'] = currentLanguage;
      
      // Add auth token
      const token = localStorage.getItem('token');
      if (token) {
        config.headers.Authorization = `Bearer ${token}`;
      }

      return config;
    },
    (error) => {
      return Promise.reject(error);
    }
  );

  return instance;
};

// Create API instances
const api = createApiInstance(false); // Unversioned for health endpoint
const apiV1 = createApiInstance(true);  // Versioned for all other endpoints

// Add response interceptor to all API instances
const addResponseInterceptor = (instance: any) => {
  instance.interceptors.response.use(
    (response) => {
      // Add version info to successful responses
      if (response.headers['api-version']) {
        response.data.api_version = response.headers['api-version'];
      }
      if (response.headers['api-version-supported']) {
        response.data.api_version_supported = response.headers['api-version-supported'] === 'true';
      }
      return response;
    },
    (error) => {
      // Handle auth errors
      if (error.response?.status === 401) {
        localStorage.removeItem('token');
        localStorage.removeItem('user');
        window.location.href = '/login';
      }
      
      // Handle API version errors
      if (error.response?.status === 406) {
        console.warn('API Version not supported:', error.response.data);
      }
      
      // Handle rate limit errors
      if (error.response?.status === 429) {
        console.warn('Rate limit exceeded:', error.response.data);
        // Could show a toast notification here
      }
      
      return Promise.reject(error);
    }
  );
};

// Apply response interceptor to all instances
addResponseInterceptor(api);
addResponseInterceptor(apiV1);

export const authService = {
  login: (username: string, password: string): Promise<AxiosResponse<any>> => {
    return apiV1.post('/auth/login', { username, password });
  },
  
  register: (userData: RegisterRequest): Promise<AxiosResponse<any>> => {
    return apiV1.post('/auth/register', userData);
  },
  
  logout: (): void => {
    localStorage.removeItem('token');
    localStorage.removeItem('user');
  },
};

export const userService = {
  getPermissions: (): Promise<AxiosResponse<any>> => {
    return apiV1.get('/me/permissions');
  },
  getProfile: (): Promise<AxiosResponse<any>> => {
    return apiV1.get('/profile');
  },
  
  updateProfile: (userData: Partial<UserData>): Promise<AxiosResponse<any>> => {
    return apiV1.put('/profile', userData);
  },
  
  changePassword: (passwordData: PasswordData): Promise<AxiosResponse<any>> => {
    return apiV1.put('/profile/password', passwordData);
  },
  
  getUsers: (searchTerm: string = ''): Promise<AxiosResponse<any>> => {
    const params = searchTerm ? { search: searchTerm } : {};
    return apiV1.get('/users', { params });
  },
  
  getUser: (userId: string): Promise<AxiosResponse<any>> => {
    return apiV1.get(`/users/${userId}`);
  },
  
  createUser: (userData: UserData): Promise<AxiosResponse<any>> => {
    return apiV1.post('/users', userData);
  },
  
  updateUser: (userId: string, userData: Partial<UserData>): Promise<AxiosResponse<any>> => {
    return apiV1.put(`/users/${userId}`, userData);
  },
  
  deleteUser: (userId: string): Promise<AxiosResponse<any>> => {
    return apiV1.delete(`/users/${userId}`);
  },
};

export const healthService = {
  check: () => {
    return api.get('/health');
  },
};

export const rbacService = {
  // Role management
  getRoles: (params: RoleQueryParams = {}): Promise<AxiosResponse<any>> => {
    const queryParams = new URLSearchParams();
    if (params.page) queryParams.append('page', params.page.toString());
    if (params.page_size) queryParams.append('page_size', params.page_size.toString());
    if (params.search) queryParams.append('search', params.search);
    if (params.status) queryParams.append('status', params.status);
    if (params.is_system !== undefined) queryParams.append('is_system', params.is_system.toString());
    if (params.sort) queryParams.append('sort', params.sort);
    if (params.order) queryParams.append('order', params.order);
    
    const url = queryParams.toString() ? `/rbac/roles?${queryParams.toString()}` : '/rbac/roles';
    return apiV1.get(url);
  },
  
  getRole: (roleId: string): Promise<AxiosResponse<any>> => {
    return apiV1.get(`/rbac/roles/${roleId}`);
  },
  
  createRole: (roleData: RoleData): Promise<AxiosResponse<any>> => {
    return apiV1.post('/rbac/roles', roleData);
  },
  
  updateRole: (roleId: string, roleData: RoleData): Promise<AxiosResponse<any>> => {
    return apiV1.put(`/rbac/roles/${roleId}`, roleData);
  },
  
  deleteRole: (roleName: string): Promise<AxiosResponse<any>> => {
    return apiV1.delete(`/rbac/roles/${roleName}`);
  },
  
  // Permission management
  getPermissions: (params = {}) => {
    return apiV1.get('/rbac/permissions', { params });
  },
  
  getAvailablePermissions: () => {
    return apiV1.get('/rbac/permissions/available');
  },
  
  getRolePermissions: (roleName) => {
    return apiV1.get(`/rbac/roles/${roleName}/permissions`);
  },
  
  addPermission: (permissionData) => {
    return apiV1.post('/rbac/permissions', permissionData);
  },
  
  removePermission: (permissionData) => {
    return apiV1.delete('/rbac/permissions', { data: permissionData });
  },
  
  // User role management
  getUserRoles: (userId) => {
    return apiV1.get(`/rbac/users/${userId}/roles`);
  },
  
  assignRole: (assignmentData) => {
    return apiV1.post('/rbac/users/assign-role', assignmentData);
  },
  
  removeRole: (removalData) => {
    return apiV1.post('/rbac/users/remove-role', removalData);
  },
  
  getUsersWithRole: (roleName) => {
    return apiV1.get(`/rbac/roles/${roleName}/users`);
  },
  
  checkPermission: (userId: string, resource: string, action: string): Promise<AxiosResponse<any>> => {
    return apiV1.get(`/rbac/users/${userId}/check-permission`, {
      params: { resource, action }
    });
  },

  // Resource and action endpoints
  getResources: (params: any = {}): Promise<AxiosResponse<any>> => {
    return apiV1.get('/rbac/resources', { params });
  },

  getActions: (params: any = {}): Promise<AxiosResponse<any>> => {
    return apiV1.get('/rbac/actions', { params });
  },

};

export default api;

