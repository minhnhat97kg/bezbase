import axios, { AxiosResponse } from 'axios';

const API_BASE_URL = process.env.REACT_APP_API_URL || 'http://localhost:8080/api';
const API_VERSION = 'v1';

// Type definitions for organization management
export interface Organization {
  id: number;
  name: string;
  slug: string;
  domain?: string;
  settings?: string;
  is_active: boolean;
  plan_type: 'free' | 'basic' | 'premium' | 'enterprise';
  created_at: string;
  updated_at: string;
  users?: OrganizationUser[];
  invitations?: OrganizationInvitation[];
  roles?: Role[];
}

export interface OrganizationUser {
  id: number;
  org_id: number;
  user_id: number;
  role: 'owner' | 'admin' | 'member' | 'viewer';
  is_primary: boolean;
  joined_at: string;
  created_at: string;
  updated_at: string;
  organization?: Organization;
  user?: User;
}

export interface OrganizationInvitation {
  id: number;
  org_id: number;
  email: string;
  role: 'owner' | 'admin' | 'member' | 'viewer';
  token: string;
  expires_at: string;
  invited_by: number;
  accepted_at?: string;
  created_at: string;
  updated_at: string;
  organization?: Organization;
  invited_by_user?: User;
}

export interface User {
  id: number;
  username: string;
  first_name: string;
  last_name: string;
  email: string;
  avatar_url?: string;
  current_org_id?: number;
}

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
}

export interface CreateOrganizationRequest {
  name: string;
  slug: string;
  domain?: string;
  plan_type?: 'free' | 'basic' | 'premium' | 'enterprise';
}

export interface UpdateOrganizationRequest {
  name?: string;
  domain?: string;
  plan_type?: 'free' | 'basic' | 'premium' | 'enterprise';
  settings?: string;
}

export interface InviteUserRequest {
  email: string;
  role: 'admin' | 'member' | 'viewer';
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

  // Add auth token to requests
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

// Organization Service
export const organizationService = {
  // Organization CRUD
  createOrganization: (data: CreateOrganizationRequest): Promise<AxiosResponse<Organization>> => {
    return api.post('/organizations', data);
  },

  getOrganization: (orgId: number): Promise<AxiosResponse<Organization>> => {
    return api.get(`/organizations/${orgId}`);
  },

  getUserOrganizations: (): Promise<AxiosResponse<OrganizationUser[]>> => {
    return api.get('/user/organizations');
  },

  updateOrganization: (orgId: number, data: UpdateOrganizationRequest): Promise<AxiosResponse<Organization>> => {
    return api.put(`/organizations/${orgId}`, data);
  },

  // User management
  inviteUser: (orgId: number, data: InviteUserRequest): Promise<AxiosResponse<OrganizationInvitation>> => {
    return api.post(`/organizations/${orgId}/invite`, data);
  },

  acceptInvitation: (token: string): Promise<AxiosResponse<ApiResponse>> => {
    return api.post(`/organizations/invitations/${token}/accept`);
  },

  removeUser: (orgId: number, userId: number): Promise<AxiosResponse<ApiResponse>> => {
    return api.delete(`/organizations/${orgId}/users/${userId}`);
  },

  updateUserRole: (orgId: number, userId: number, role: string): Promise<AxiosResponse<ApiResponse>> => {
    return api.put(`/organizations/${orgId}/users/${userId}/role`, { role });
  },

  // Organization context
  switchOrganization: (orgId: number): Promise<AxiosResponse<ApiResponse>> => {
    return api.post(`/organizations/${orgId}/switch`);
  },

  // Utility functions
  setCurrentOrganization: (orgId: number): void => {
    localStorage.setItem('currentOrgId', orgId.toString());
  },

  getCurrentOrganization: (): number | null => {
    const orgId = localStorage.getItem('currentOrgId');
    return orgId ? parseInt(orgId, 10) : null;
  },

  clearCurrentOrganization: (): void => {
    localStorage.removeItem('currentOrgId');
  },

  // Generate slug from name
  generateSlug: (name: string): string => {
    return name
      .toLowerCase()
      .replace(/[^a-z0-9]+/g, '-')
      .replace(/^-+|-+$/g, '')
      .substring(0, 100);
  },

  // Validate slug
  isValidSlug: (slug: string): boolean => {
    const slugRegex = /^[a-z0-9]([a-z0-9-]*[a-z0-9])?$/;
    return slugRegex.test(slug) && slug.length >= 2 && slug.length <= 100;
  },

  // Check if user can manage organization
  canManageOrganization: (orgUser: OrganizationUser): boolean => {
    return orgUser.role === 'owner';
  },

  // Check if user can manage members
  canManageMembers: (orgUser: OrganizationUser): boolean => {
    return orgUser.role === 'owner' || orgUser.role === 'admin';
  },

  // Check if user can invite members
  canInviteMembers: (orgUser: OrganizationUser): boolean => {
    return orgUser.role === 'owner' || orgUser.role === 'admin';
  },

  // Format role display
  formatRole: (role: string): string => {
    return role.charAt(0).toUpperCase() + role.slice(1);
  },

  // Format plan type
  formatPlanType: (planType: string): string => {
    return planType.charAt(0).toUpperCase() + planType.slice(1);
  },

  // Get role color for UI
  getRoleColor: (role: string): string => {
    switch (role) {
      case 'owner':
        return 'bg-purple-100 text-purple-800';
      case 'admin':
        return 'bg-red-100 text-red-800';
      case 'member':
        return 'bg-blue-100 text-blue-800';
      case 'viewer':
        return 'bg-gray-100 text-gray-800';
      default:
        return 'bg-gray-100 text-gray-800';
    }
  },

  // Get plan color for UI
  getPlanColor: (planType: string): string => {
    switch (planType) {
      case 'free':
        return 'bg-gray-100 text-gray-800';
      case 'basic':
        return 'bg-blue-100 text-blue-800';
      case 'premium':
        return 'bg-purple-100 text-purple-800';
      case 'enterprise':
        return 'bg-yellow-100 text-yellow-800';
      default:
        return 'bg-gray-100 text-gray-800';
    }
  },
};

export default organizationService;