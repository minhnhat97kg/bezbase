import React, { createContext, useContext, useState, useEffect, ReactNode } from 'react';
import { organizationService, Organization, OrganizationUser } from '../services/organizationService';

interface OrganizationContextType {
  // Current organization state
  currentOrganization: Organization | null;
  currentOrgUser: OrganizationUser | null;
  userOrganizations: OrganizationUser[];
  
  // Loading states
  isLoading: boolean;
  isLoadingOrganizations: boolean;
  
  // Actions
  switchOrganization: (orgId: number) => Promise<void>;
  refreshOrganizations: () => Promise<void>;
  refreshCurrentOrganization: () => Promise<void>;
  
  // Permissions
  canManageOrganization: () => boolean;
  canManageMembers: () => boolean;
  canInviteMembers: () => boolean;
  
  // Utilities
  getCurrentOrgId: () => number | null;
  setCurrentOrgId: (orgId: number | null) => void;
}

const OrganizationContext = createContext<OrganizationContextType | undefined>(undefined);

interface OrganizationProviderProps {
  children: ReactNode;
}

export const OrganizationProvider: React.FC<OrganizationProviderProps> = ({ children }) => {
  const [currentOrganization, setCurrentOrganization] = useState<Organization | null>(null);
  const [currentOrgUser, setCurrentOrgUser] = useState<OrganizationUser | null>(null);
  const [userOrganizations, setUserOrganizations] = useState<OrganizationUser[]>([]);
  const [isLoading, setIsLoading] = useState(false);
  const [isLoadingOrganizations, setIsLoadingOrganizations] = useState(false);

  // Initialize organization context
  useEffect(() => {
    const initializeOrganizations = async () => {
      const token = localStorage.getItem('token');
      if (token) {
        await refreshOrganizations();
      }
    };

    initializeOrganizations();
  }, []);

  // Load current organization when currentOrgId changes
  useEffect(() => {
    const currentOrgId = organizationService.getCurrentOrganization();
    if (currentOrgId && userOrganizations.length > 0) {
      loadCurrentOrganization(currentOrgId);
    }
  }, [userOrganizations]);

  const refreshOrganizations = async (): Promise<void> => {
    setIsLoadingOrganizations(true);
    try {
      console.log('Fetching user organizations...');
      const response = await organizationService.getUserOrganizations();
      const orgs = Array.isArray(response.data) ? response.data : [];
      console.log('User organizations loaded:', orgs);
      setUserOrganizations(orgs);

      // If no current organization is set and user has organizations, set the first one
      const currentOrgId = organizationService.getCurrentOrganization();
      if (!currentOrgId && orgs.length > 0) {
        const primaryOrg = orgs.find(org => org.is_primary) || orgs[0];
        await switchOrganization(primaryOrg.org_id);
      } else if (currentOrgId) {
        // Verify current organization is still valid
        const currentOrgUser = orgs.find(org => org.org_id === currentOrgId);
        if (!currentOrgUser) {
          // Current org is no longer valid, switch to first available
          if (orgs.length > 0) {
            await switchOrganization(orgs[0].org_id);
          } else {
            // No organizations available
            organizationService.clearCurrentOrganization();
            setCurrentOrganization(null);
            setCurrentOrgUser(null);
          }
        }
      }
    } catch (error) {
      console.error('Failed to load organizations:', error);
      setUserOrganizations([]);
      setCurrentOrganization(null);
      setCurrentOrgUser(null);
    } finally {
      setIsLoadingOrganizations(false);
    }
  };

  const loadCurrentOrganization = async (orgId: number): Promise<void> => {
    try {
      console.log('Loading organization:', orgId);
      const orgUser = userOrganizations.find(org => org.org_id === orgId);
      if (!orgUser) {
        console.warn('Organization not found in user organizations:', orgId);
        throw new Error('Organization not found in user organizations');
      }

      setCurrentOrgUser(orgUser);
      console.log('Set current org user:', orgUser);

      // Load full organization details
      const response = await organizationService.getOrganization(orgId);
      console.log('Loaded organization details:', response.data);
      setCurrentOrganization(response.data);
    } catch (error) {
      console.error('Failed to load current organization:', error);
      setCurrentOrganization(null);
      setCurrentOrgUser(null);
    }
  };

  const refreshCurrentOrganization = async (): Promise<void> => {
    const currentOrgId = organizationService.getCurrentOrganization();
    if (currentOrgId) {
      await loadCurrentOrganization(currentOrgId);
    }
  };

  const switchOrganization = async (orgId: number): Promise<void> => {
    setIsLoading(true);
    try {
      // Call API to switch organization context
      await organizationService.switchOrganization(orgId);
      
      // Update local storage
      organizationService.setCurrentOrganization(orgId);
      
      // Load the new organization
      await loadCurrentOrganization(orgId);
      
      // Refresh the page to update all components with new context
      window.location.reload();
    } catch (error) {
      console.error('Failed to switch organization:', error);
      throw error;
    } finally {
      setIsLoading(false);
    }
  };

  const canManageOrganization = (): boolean => {
    return currentOrgUser ? organizationService.canManageOrganization(currentOrgUser) : false;
  };

  const canManageMembers = (): boolean => {
    return currentOrgUser ? organizationService.canManageMembers(currentOrgUser) : false;
  };

  const canInviteMembers = (): boolean => {
    return currentOrgUser ? organizationService.canInviteMembers(currentOrgUser) : false;
  };

  const getCurrentOrgId = (): number | null => {
    return organizationService.getCurrentOrganization();
  };

  const setCurrentOrgId = (orgId: number | null): void => {
    if (orgId) {
      organizationService.setCurrentOrganization(orgId);
    } else {
      organizationService.clearCurrentOrganization();
    }
  };

  const value: OrganizationContextType = {
    currentOrganization,
    currentOrgUser,
    userOrganizations,
    isLoading,
    isLoadingOrganizations,
    switchOrganization,
    refreshOrganizations,
    refreshCurrentOrganization,
    canManageOrganization,
    canManageMembers,
    canInviteMembers,
    getCurrentOrgId,
    setCurrentOrgId,
  };

  return (
    <OrganizationContext.Provider value={value}>
      {children}
    </OrganizationContext.Provider>
  );
};

export const useOrganization = (): OrganizationContextType => {
  const context = useContext(OrganizationContext);
  if (context === undefined) {
    throw new Error('useOrganization must be used within an OrganizationProvider');
  }
  return context;
};

export default OrganizationContext;