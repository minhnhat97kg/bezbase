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
};

export const healthService = {
  check: () => {
    return api.get('/health');
  },
};

export const rbacService = {
  // Role management
  getRoles: () => {
    return api.get('/v1/rbac/roles');
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
};

export default api;

