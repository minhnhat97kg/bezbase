import React, { useState, useRef, useEffect } from 'react';
import { useOrganization } from '../../context/OrganizationContext';
import { organizationService } from '../../services/organizationService';

interface OrganizationSelectorProps {
  className?: string;
}

const OrganizationSelector: React.FC<OrganizationSelectorProps> = ({ className = '' }) => {
  const {
    currentOrganization,
    userOrganizations,
    switchOrganization,
    isLoading,
    isLoadingOrganizations,
  } = useOrganization();

  const [isOpen, setIsOpen] = useState(false);
  const dropdownRef = useRef<HTMLDivElement>(null);

  // Close dropdown when clicking outside
  useEffect(() => {
    const handleClickOutside = (event: MouseEvent) => {
      if (dropdownRef.current && !dropdownRef.current.contains(event.target as Node)) {
        setIsOpen(false);
      }
    };

    document.addEventListener('mousedown', handleClickOutside);
    return () => document.removeEventListener('mousedown', handleClickOutside);
  }, []);

  const handleOrganizationSwitch = async (orgId: number) => {
    if (orgId === currentOrganization?.id) {
      setIsOpen(false);
      return;
    }

    try {
      await switchOrganization(orgId);
      setIsOpen(false);
    } catch (error) {
      console.error('Failed to switch organization:', error);
      // Could show a toast notification here
    }
  };

  if (isLoadingOrganizations) {
    return (
      <div className={`animate-pulse ${className}`}>
        <div className="h-10 bg-gray-200 rounded-md"></div>
      </div>
    );
  }

  if (userOrganizations.length === 0) {
    return (
      <div className={`relative ${className}`}>
        <button
          type="button"
          className="flex items-center w-full px-3 py-2 text-left bg-white border border-gray-300 rounded-md shadow-sm hover:bg-gray-50 focus:outline-none focus:ring-2 focus:ring-indigo-500 focus:border-indigo-500"
          onClick={() => window.location.href = '/organizations/create'}
        >
          <div className="flex-1 min-w-0">
            <div className="flex items-center">
              <div className="flex-shrink-0">
                <div className="w-8 h-8 bg-gray-100 rounded-full flex items-center justify-center">
                  <svg className="w-5 h-5 text-gray-400" fill="none" stroke="currentColor" viewBox="0 0 24 24">
                    <path strokeLinecap="round" strokeLinejoin="round" strokeWidth={2} d="M12 6v6m0 0v6m0-6h6m-6 0H6" />
                  </svg>
                </div>
              </div>
              <div className="ml-3 min-w-0 flex-1">
                <p className="text-sm font-medium text-gray-700 truncate">
                  Create Organization
                </p>
                <p className="text-xs text-gray-500 truncate">
                  No organizations found
                </p>
              </div>
            </div>
          </div>
        </button>
      </div>
    );
  }

  return (
    <div className={`relative ${className}`} ref={dropdownRef}>
      <button
        type="button"
        className="flex items-center w-full px-3 py-2 text-left bg-white border border-gray-300 rounded-md shadow-sm hover:bg-gray-50 focus:outline-none focus:ring-2 focus:ring-indigo-500 focus:border-indigo-500"
        onClick={() => setIsOpen(!isOpen)}
        disabled={isLoading}
      >
        <div className="flex-1 min-w-0">
          <div className="flex items-center">
            <div className="flex-shrink-0">
              <div className="w-8 h-8 bg-indigo-100 rounded-full flex items-center justify-center">
                <span className="text-sm font-medium text-indigo-600">
                  {currentOrganization?.name.charAt(0).toUpperCase()}
                </span>
              </div>
            </div>
            <div className="ml-3 min-w-0 flex-1">
              <p className="text-sm font-medium text-gray-900 truncate">
                {currentOrganization?.name || 'Select Organization'}
              </p>
              <p className="text-xs text-gray-500 truncate">
                {currentOrganization && (
                  <span className={`inline-block px-2 py-1 rounded-full text-xs ${organizationService.getPlanColor(currentOrganization.plan_type)}`}>
                    {organizationService.formatPlanType(currentOrganization.plan_type)}
                  </span>
                )}
              </p>
            </div>
          </div>
        </div>
        <div className="flex-shrink-0 ml-2">
          {isLoading ? (
            <div className="w-5 h-5 animate-spin border-2 border-indigo-500 border-t-transparent rounded-full"></div>
          ) : (
            <svg
              className={`w-5 h-5 text-gray-400 transition-transform ${isOpen ? 'transform rotate-180' : ''}`}
              fill="none"
              stroke="currentColor"
              viewBox="0 0 24 24"
            >
              <path strokeLinecap="round" strokeLinejoin="round" strokeWidth={2} d="M19 9l-7 7-7-7" />
            </svg>
          )}
        </div>
      </button>

      {isOpen && (
        <div className="absolute z-50 w-full mt-1 bg-white border border-gray-300 rounded-md shadow-lg">
          <div className="py-1 max-h-60 overflow-auto">
            {userOrganizations.map((orgUser) => {
              const isSelected = orgUser.org_id === currentOrganization?.id;
              return (
                <button
                  key={orgUser.org_id}
                  type="button"
                  className={`w-full px-4 py-2 text-left hover:bg-gray-50 focus:outline-none focus:bg-gray-50 ${
                    isSelected ? 'bg-indigo-50' : ''
                  }`}
                  onClick={() => handleOrganizationSwitch(orgUser.org_id)}
                >
                  <div className="flex items-center">
                    <div className="flex-shrink-0">
                      <div className={`w-8 h-8 rounded-full flex items-center justify-center ${
                        isSelected ? 'bg-indigo-100' : 'bg-gray-100'
                      }`}>
                        <span className={`text-sm font-medium ${
                          isSelected ? 'text-indigo-600' : 'text-gray-600'
                        }`}>
                          {orgUser.organization?.name.charAt(0).toUpperCase()}
                        </span>
                      </div>
                    </div>
                    <div className="ml-3 min-w-0 flex-1">
                      <div className="flex items-center justify-between">
                        <p className={`text-sm font-medium truncate ${
                          isSelected ? 'text-indigo-900' : 'text-gray-900'
                        }`}>
                          {orgUser.organization?.name}
                        </p>
                        {isSelected && (
                          <svg className="w-4 h-4 text-indigo-600 flex-shrink-0" fill="currentColor" viewBox="0 0 20 20">
                            <path fillRule="evenodd" d="M16.707 5.293a1 1 0 010 1.414l-8 8a1 1 0 01-1.414 0l-4-4a1 1 0 011.414-1.414L8 12.586l7.293-7.293a1 1 0 011.414 0z" clipRule="evenodd" />
                          </svg>
                        )}
                      </div>
                      <div className="flex items-center mt-1 space-x-2">
                        <span className={`inline-block px-2 py-1 rounded-full text-xs ${organizationService.getRoleColor(orgUser.role)}`}>
                          {organizationService.formatRole(orgUser.role)}
                        </span>
                        {orgUser.organization && (
                          <span className={`inline-block px-2 py-1 rounded-full text-xs ${organizationService.getPlanColor(orgUser.organization.plan_type)}`}>
                            {organizationService.formatPlanType(orgUser.organization.plan_type)}
                          </span>
                        )}
                        {orgUser.is_primary && (
                          <span className="inline-block px-2 py-1 bg-green-100 text-green-800 rounded-full text-xs">
                            Primary
                          </span>
                        )}
                      </div>
                    </div>
                  </div>
                </button>
              );
            })}
          </div>
          
          <div className="border-t border-gray-200">
            <button
              type="button"
              className="w-full px-4 py-2 text-left text-sm text-indigo-600 hover:bg-gray-50 focus:outline-none focus:bg-gray-50"
              onClick={() => {
                setIsOpen(false);
                // Navigate to create organization page
                window.location.href = '/organizations/create';
              }}
            >
              <div className="flex items-center">
                <svg className="w-5 h-5 mr-3" fill="none" stroke="currentColor" viewBox="0 0 24 24">
                  <path strokeLinecap="round" strokeLinejoin="round" strokeWidth={2} d="M12 6v6m0 0v6m0-6h6m-6 0H6" />
                </svg>
                Create New Organization
              </div>
            </button>
          </div>
        </div>
      )}
    </div>
  );
};

export default OrganizationSelector;