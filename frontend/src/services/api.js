import axios from 'axios';

const API_BASE_URL = process.env.REACT_APP_API_URL || 'http://localhost:8080/api';

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
  login: (username, password) => {
    return api.post('/auth/login', { username, password });
  },
  
  register: (userData) => {
    return api.post('/auth/register', userData);
  },
  
  logout: () => {
    localStorage.removeItem('token');
    localStorage.removeItem('user');
  },
};

export const userService = {
  getPermissions: () => {
    return api.get('/v1/me/permissions');
  },
  getProfile: () => {
    return api.get('/v1/profile');
  },
  
  updateProfile: (userData) => {
    return api.put('/v1/profile', userData);
  },
  
  getUsers: (searchTerm = '') => {
    const params = searchTerm ? { search: searchTerm } : {};
    return api.get('/v1/users', { params });
  },
  
  getUser: (userId) => {
    return api.get(`/v1/users/${userId}`);
  },
  
  createUser: (userData) => {
    return api.post('/v1/users', userData);
  },
  
  updateUser: (userId, userData) => {
    return api.put(`/v1/users/${userId}`, userData);
  },
  
  deleteUser: (userId) => {
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
  getRoles: (params = {}) => {
    const queryParams = new URLSearchParams();
    if (params.page) queryParams.append('page', params.page);
    if (params.page_size) queryParams.append('page_size', params.page_size);
    if (params.search) queryParams.append('search', params.search);
    if (params.status) queryParams.append('status', params.status);
    if (params.is_system !== undefined) queryParams.append('is_system', params.is_system);
    if (params.sort) queryParams.append('sort', params.sort);
    if (params.order) queryParams.append('order', params.order);
    
    const url = queryParams.toString() ? `/v1/rbac/roles?${queryParams.toString()}` : '/v1/rbac/roles';
    return api.get(url);
  },
  
  getRole: (roleId) => {
    return api.get(`/v1/rbac/roles/${roleId}`);
  },
  
  createRole: (roleData) => {
    return api.post('/v1/rbac/roles', roleData);
  },
  
  updateRole: (roleId, roleData) => {
    return api.put(`/v1/rbac/roles/${roleId}`, roleData);
  },
  
  deleteRole: (roleName) => {
    return api.delete(`/v1/rbac/roles/${roleName}`);
  },
  
  // Permission management
  getPermissions: (params = {}) => {
    return api.get('/v1/rbac/permissions', { params });
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
  
  checkPermission: (userId, resource, action) => {
    return api.get(`/v1/rbac/users/${userId}/check-permission`, {
      params: { resource, action }
    });
  },

  // Resource management
  getResources: (params = {}) => {
    const queryParams = new URLSearchParams();
    if (params.page) queryParams.append('page', params.page);
    if (params.page_size) queryParams.append('page_size', params.page_size);
    if (params.search) queryParams.append('search', params.search);
    
    const url = queryParams.toString() ? `/v1/rbac/resources?${queryParams.toString()}` : '/v1/rbac/resources';
    return api.get(url);
  },

  // Action management
  getActions: (params = {}) => {
    const queryParams = new URLSearchParams();
    if (params.page) queryParams.append('page', params.page);
    if (params.page_size) queryParams.append('page_size', params.page_size);
    if (params.search) queryParams.append('search', params.search);
    
    const url = queryParams.toString() ? `/v1/rbac/actions?${queryParams.toString()}` : '/v1/rbac/actions';
    return api.get(url);
  },
};

export default api;

