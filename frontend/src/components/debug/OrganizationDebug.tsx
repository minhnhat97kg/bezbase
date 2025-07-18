import React from 'react';
import { useOrganization } from '../../context/OrganizationContext';

const OrganizationDebug: React.FC = () => {
  const {
    currentOrganization,
    currentOrgUser,
    userOrganizations,
    isLoading,
    isLoadingOrganizations,
    getCurrentOrgId,
  } = useOrganization();

  // Only show in development
  if (process.env.NODE_ENV !== 'development') {
    return null;
  }

  const currentOrgId = getCurrentOrgId();

  return (
    <div className="fixed bottom-4 right-4 bg-white border border-gray-300 rounded-lg shadow-lg p-4 max-w-sm text-xs z-50">
      <h3 className="font-bold text-gray-900 mb-2">Organization Debug</h3>
      
      <div className="space-y-2">
        <div>
          <span className="font-medium">Loading States:</span>
          <div className="ml-2">
            <div>isLoading: {isLoading ? '✅' : '❌'}</div>
            <div>isLoadingOrganizations: {isLoadingOrganizations ? '✅' : '❌'}</div>
          </div>
        </div>

        <div>
          <span className="font-medium">Current Org ID (localStorage):</span>
          <div className="ml-2 font-mono bg-gray-100 px-1 rounded">
            {currentOrgId || 'null'}
          </div>
        </div>

        <div>
          <span className="font-medium">User Organizations ({userOrganizations.length}):</span>
          <div className="ml-2 max-h-20 overflow-y-auto">
            {userOrganizations.map((org) => (
              <div key={org.org_id} className="text-xs">
                #{org.org_id}: {org.organization?.name || 'Unknown'} 
                {org.is_primary && ' (Primary)'}
                <span className="text-gray-500"> - {org.role}</span>
              </div>
            ))}
            {userOrganizations.length === 0 && (
              <div className="text-gray-500 italic">No organizations</div>
            )}
          </div>
        </div>

        <div>
          <span className="font-medium">Current Organization:</span>
          <div className="ml-2">
            {currentOrganization ? (
              <div className="font-mono bg-gray-100 px-1 rounded">
                #{currentOrganization.id}: {currentOrganization.name}
              </div>
            ) : (
              <div className="text-red-500 italic">None selected</div>
            )}
          </div>
        </div>

        <div>
          <span className="font-medium">Current Org User:</span>
          <div className="ml-2">
            {currentOrgUser ? (
              <div className="font-mono bg-gray-100 px-1 rounded">
                {currentOrgUser.role} in #{currentOrgUser.org_id}
              </div>
            ) : (
              <div className="text-red-500 italic">None</div>
            )}
          </div>
        </div>

        <div>
          <span className="font-medium">Auth Token:</span>
          <div className="ml-2 font-mono bg-gray-100 px-1 rounded">
            {localStorage.getItem('token') ? '✅ Present' : '❌ Missing'}
          </div>
        </div>
      </div>
    </div>
  );
};

export default OrganizationDebug;