import axios, { AxiosResponse } from 'axios';

const API_BASE_URL = process.env.REACT_APP_API_URL || 'http://localhost:8080/api';

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

const api = axios.create({
  baseURL: API_BASE_URL,
  headers: {
    'Content-Type': 'application/json',
  },
});

// Request interceptor to add auth token
api.interceptors.request.use(
  (config) => {
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

// Response interceptor to handle auth errors
api.interceptors.response.use(
  (response) => response,
  (error) => {
    if (error.response?.status === 401) {
      localStorage.removeItem('token');
      localStorage.removeItem('user');
      window.location.href = '/login';
    }
    return Promise.reject(error);
  }
);

export const authService = {
  login: (username: string, password: string): Promise<AxiosResponse<any>> => {
    return api.post('/auth/login', { username, password });
  },
  
  register: (userData: RegisterRequest): Promise<AxiosResponse<any>> => {
    return api.post('/auth/register', userData);
  },
  
  logout: (): void => {
    localStorage.removeItem('token');
    localStorage.removeItem('user');
  },
};

export const userService = {
  getPermissions: (): Promise<AxiosResponse<any>> => {
    return api.get('/v1/me/permissions');
  },
  getProfile: (): Promise<AxiosResponse<any>> => {
    return api.get('/v1/profile');
  },
  
  updateProfile: (userData: Partial<UserData>): Promise<AxiosResponse<any>> => {
    return api.put('/v1/profile', userData);
  },
  
  changePassword: (passwordData: PasswordData): Promise<AxiosResponse<any>> => {
    return api.put('/v1/profile/password', passwordData);
  },
  
  getUsers: (searchTerm: string = ''): Promise<AxiosResponse<any>> => {
    const params = searchTerm ? { search: searchTerm } : {};
    return api.get('/v1/users', { params });
  },
  
  getUser: (userId: string): Promise<AxiosResponse<any>> => {
    return api.get(`/v1/users/${userId}`);
  },
  
  createUser: (userData: UserData): Promise<AxiosResponse<any>> => {
    return api.post('/v1/users', userData);
  },
  
  updateUser: (userId: string, userData: Partial<UserData>): Promise<AxiosResponse<any>> => {
    return api.put(`/v1/users/${userId}`, userData);
  },
  
  deleteUser: (userId: string): Promise<AxiosResponse<any>> => {
    return api.delete(`/v1/users/${userId}`);
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
    
    const url = queryParams.toString() ? `/v1/rbac/roles?${queryParams.toString()}` : '/v1/rbac/roles';
    return api.get(url);
  },
  
  getRole: (roleId: string): Promise<AxiosResponse<any>> => {
    return api.get(`/v1/rbac/roles/${roleId}`);
  },
  
  createRole: (roleData: RoleData): Promise<AxiosResponse<any>> => {
    return api.post('/v1/rbac/roles', roleData);
  },
  
  updateRole: (roleId: string, roleData: RoleData): Promise<AxiosResponse<any>> => {
    return api.put(`/v1/rbac/roles/${roleId}`, roleData);
  },
  
  deleteRole: (roleName: string): Promise<AxiosResponse<any>> => {
    return api.delete(`/v1/rbac/roles/${roleName}`);
  },
  
  // Permission management
  getPermissions: (params = {}) => {
    return api.get('/v1/rbac/permissions', { params });
  },
  
  getAvailablePermissions: () => {
    return api.get('/v1/rbac/permissions/available');
  },
  
  getRolePermissions: (roleName) => {
    return api.get(`/v1/rbac/roles/${roleName}/permissions`);
  },
  
  addPermission: (permissionData) => {
    return api.post('/v1/rbac/permissions', permissionData);
  },
  
  removePermission: (permissionData) => {
    return api.delete('/v1/rbac/permissions', { data: permissionData });
  },
  
  // User role management
  getUserRoles: (userId) => {
    return api.get(`/v1/rbac/users/${userId}/roles`);
  },
  
  assignRole: (assignmentData) => {
    return api.post('/v1/rbac/users/assign-role', assignmentData);
  },
  
  removeRole: (removalData) => {
    return api.post('/v1/rbac/users/remove-role', removalData);
  },
  
  getUsersWithRole: (roleName) => {
    return api.get(`/v1/rbac/roles/${roleName}/users`);
  },
  
  checkPermission: (userId: string, resource: string, action: string): Promise<AxiosResponse<any>> => {
    return api.get(`/v1/rbac/users/${userId}/check-permission`, {
      params: { resource, action }
    });
  },

  // Resource and action endpoints
  getResources: (params: any = {}): Promise<AxiosResponse<any>> => {
    return api.get('/v1/rbac/resources', { params });
  },

  getActions: (params: any = {}): Promise<AxiosResponse<any>> => {
    return api.get('/v1/rbac/actions', { params });
  },

};

export default api;

