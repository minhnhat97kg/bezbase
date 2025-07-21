import React, { useState, useEffect } from 'react';
import { rbacService, Role } from '../../services/rbacService';

interface EffectivePermissionsViewProps {
  userId?: number;
  roleId?: number;
  role?: Role;
  organizationId?: number;
}

interface PermissionItem {
  id: number;
  resource: string;
  action: string;
  context_type?: string;
  context_value?: string;
  is_granted: boolean;
  source: 'direct' | 'inherited';
  source_role?: string;
  role_id: number;
}

interface PermissionGroup {
  resource: string;
  permissions: PermissionItem[];
}

const EffectivePermissionsView: React.FC<EffectivePermissionsViewProps> = ({
  userId,
  roleId,
  role,
  organizationId
}) => {
  const [permissions, setPermissions] = useState<PermissionItem[]>([]);
  const [loading, setLoading] = useState(false);
  const [error, setError] = useState<string | null>(null);
  const [groupBy, setGroupBy] = useState<'resource' | 'source'>('resource');
  const [showInherited, setShowInherited] = useState(true);
  const [showDirect, setShowDirect] = useState(true);

  useEffect(() => {
    if (userId || roleId) {
      loadEffectivePermissions();
    }
  }, [userId, roleId, organizationId]);

  const loadEffectivePermissions = async () => {
    if (!userId && !roleId) return;

    try {
      setLoading(true);
      setError(null);
      
      let response;
      if (userId) {
        response = await rbacService.getEffectivePermissions(userId, organizationId);
      } else if (roleId) {
        // For role-based permissions, we'll need to implement this endpoint
        // For now, we'll use a placeholder approach
        response = { data: [] };
      }

      // Transform the data to include source information
      const transformedPermissions: PermissionItem[] = response.data.map((perm: any) => ({
        ...perm,
        source: perm.inherited_from ? 'inherited' : 'direct',
        source_role: perm.inherited_from || (role?.name)
      }));

      setPermissions(transformedPermissions);
    } catch (err: any) {
      console.error('Failed to load effective permissions:', err);
      setError(err.message || 'Failed to load permissions');
    } finally {
      setLoading(false);
    }
  };

  const getFilteredPermissions = () => {
    return permissions.filter(perm => {
      if (!showInherited && perm.source === 'inherited') return false;
      if (!showDirect && perm.source === 'direct') return false;
      return true;
    });
  };

  const getGroupedPermissions = (): PermissionGroup[] => {
    const filtered = getFilteredPermissions();
    
    if (groupBy === 'resource') {
      const groups = filtered.reduce((acc, perm) => {
        if (!acc[perm.resource]) {
          acc[perm.resource] = [];
        }
        acc[perm.resource].push(perm);
        return acc;
      }, {} as Record<string, PermissionItem[]>);

      return Object.entries(groups).map(([resource, perms]) => ({
        resource,
        permissions: perms.sort((a, b) => a.action.localeCompare(b.action))
      }));
    } else {
      const groups = filtered.reduce((acc, perm) => {
        const source = perm.source_role || perm.source;
        if (!acc[source]) {
          acc[source] = [];
        }
        acc[source].push(perm);
        return acc;
      }, {} as Record<string, PermissionItem[]>);

      return Object.entries(groups).map(([resource, perms]) => ({
        resource,
        permissions: perms.sort((a, b) => a.resource.localeCompare(b.resource) || a.action.localeCompare(b.action))
      }));
    }
  };

  const getPermissionIcon = (action: string) => {
    switch (action.toLowerCase()) {
      case 'create': return '➕';
      case 'read': return '👁️';
      case 'update': return '✏️';
      case 'delete': return '🗑️';
      case 'all': return '🔑';
      default: return '🔒';
    }
  };

  const getResourceIcon = (resource: string) => {
    switch (resource.toLowerCase()) {
      case 'users': return '👥';
      case 'roles': return '🎭';
      case 'permissions': return '🔐';
      case 'admin': return '⚙️';
      case 'posts': return '📝';
      case 'profile': return '👤';
      default: return '📄';
    }
  };

  if (loading) {
    return (
      <div className="bg-white shadow rounded-lg p-6">
        <div className="flex justify-center items-center py-8">
          <div className="animate-spin rounded-full h-8 w-8 border-b-2 border-indigo-500"></div>
          <span className="ml-3 text-gray-600">Loading effective permissions...</span>
        </div>
      </div>
    );
  }

  if (error) {
    return (
      <div className="bg-white shadow rounded-lg p-6">
        <div className="text-center py-8">
          <span className="text-red-500 text-4xl mb-4 block">⚠️</span>
          <p className="text-red-600 mb-4">{error}</p>
          <button
            onClick={loadEffectivePermissions}
            className="bg-indigo-600 hover:bg-indigo-700 text-white px-4 py-2 rounded-md text-sm"
          >
            Retry
          </button>
        </div>
      </div>
    );
  }

  const groupedPermissions = getGroupedPermissions();
  const filteredCount = getFilteredPermissions().length;
  const inheritedCount = permissions.filter(p => p.source === 'inherited').length;
  const directCount = permissions.filter(p => p.source === 'direct').length;

  return (
    <div className="bg-white shadow rounded-lg">
      <div className="px-4 py-5 sm:px-6 border-b border-gray-200">
        <div className="flex justify-between items-center">
          <div>
            <h3 className="text-lg font-medium text-gray-900">Effective Permissions</h3>
            <p className="mt-1 text-sm text-gray-500">
              {userId ? 'User permissions' : 'Role permissions'} including inherited permissions
            </p>
          </div>
          <button
            onClick={loadEffectivePermissions}
            className="bg-gray-100 hover:bg-gray-200 text-gray-700 px-3 py-1 rounded-md text-sm"
          >
            Refresh
          </button>
        </div>

        {/* Stats */}
        <div className="mt-4 grid grid-cols-3 gap-4">
          <div className="bg-blue-50 rounded-lg p-3">
            <div className="text-center">
              <div className="text-2xl font-bold text-blue-900">{filteredCount}</div>
              <div className="text-sm text-blue-700">Total Visible</div>
            </div>
          </div>
          <div className="bg-green-50 rounded-lg p-3">
            <div className="text-center">
              <div className="text-2xl font-bold text-green-900">{directCount}</div>
              <div className="text-sm text-green-700">Direct</div>
            </div>
          </div>
          <div className="bg-purple-50 rounded-lg p-3">
            <div className="text-center">
              <div className="text-2xl font-bold text-purple-900">{inheritedCount}</div>
              <div className="text-sm text-purple-700">Inherited</div>
            </div>
          </div>
        </div>
      </div>

      <div className="p-6">
        {/* Controls */}
        <div className="flex justify-between items-center mb-6">
          <div className="flex space-x-4">
            <div className="flex items-center space-x-2">
              <label className="text-sm font-medium text-gray-700">Show:</label>
              <label className="flex items-center">
                <input
                  type="checkbox"
                  checked={showDirect}
                  onChange={(e) => setShowDirect(e.target.checked)}
                  className="mr-1"
                />
                <span className="text-sm text-gray-600">Direct</span>
              </label>
              <label className="flex items-center">
                <input
                  type="checkbox"
                  checked={showInherited}
                  onChange={(e) => setShowInherited(e.target.checked)}
                  className="mr-1"
                />
                <span className="text-sm text-gray-600">Inherited</span>
              </label>
            </div>
          </div>

          <div className="flex items-center space-x-2">
            <label className="text-sm font-medium text-gray-700">Group by:</label>
            <select
              value={groupBy}
              onChange={(e) => setGroupBy(e.target.value as 'resource' | 'source')}
              className="border border-gray-300 rounded-md px-2 py-1 text-sm"
            >
              <option value="resource">Resource</option>
              <option value="source">Source Role</option>
            </select>
          </div>
        </div>

        {/* Permissions List */}
        {groupedPermissions.length === 0 ? (
          <div className="text-center py-8">
            <span className="text-gray-400 text-4xl mb-4 block">🔒</span>
            <p className="text-gray-500">No permissions found with current filters</p>
          </div>
        ) : (
          <div className="space-y-6">
            {groupedPermissions.map((group) => (
              <div key={group.resource} className="border border-gray-200 rounded-lg">
                <div className="bg-gray-50 px-4 py-3 border-b border-gray-200">
                  <div className="flex items-center">
                    <span className="mr-3 text-lg">
                      {groupBy === 'resource' ? getResourceIcon(group.resource) : '🎭'}
                    </span>
                    <h4 className="text-sm font-medium text-gray-900 capitalize">
                      {group.resource}
                    </h4>
                    <span className="ml-2 bg-gray-200 text-gray-700 px-2 py-1 rounded-full text-xs">
                      {group.permissions.length} permissions
                    </span>
                  </div>
                </div>

                <div className="p-4">
                  <div className="grid grid-cols-1 gap-3">
                    {group.permissions.map((permission, index) => (
                      <div
                        key={`${permission.resource}-${permission.action}-${index}`}
                        className={`flex items-center justify-between p-3 rounded-lg border ${
                          permission.source === 'inherited' 
                            ? 'bg-purple-50 border-purple-200' 
                            : 'bg-green-50 border-green-200'
                        }`}
                      >
                        <div className="flex items-center">
                          <span className="mr-3 text-lg">
                            {getPermissionIcon(permission.action)}
                          </span>
                          <div>
                            <div className="flex items-center space-x-2">
                              <span className="font-medium text-gray-900">
                                {permission.action}
                              </span>
                              {groupBy === 'source' && (
                                <span className="text-sm text-gray-600">
                                  on {permission.resource}
                                </span>
                              )}
                              {permission.context_type && (
                                <span className="inline-flex items-center px-2 py-0.5 rounded-full text-xs font-medium bg-gray-100 text-gray-800">
                                  {permission.context_type}: {permission.context_value}
                                </span>
                              )}
                            </div>
                            
                            <div className="flex items-center mt-1 space-x-2">
                              <span className={`inline-flex items-center px-2 py-0.5 rounded-full text-xs font-medium ${
                                permission.source === 'inherited'
                                  ? 'bg-purple-100 text-purple-800'
                                  : 'bg-green-100 text-green-800'
                              }`}>
                                {permission.source === 'inherited' ? '↗ Inherited' : '● Direct'}
                              </span>
                              
                              {permission.source === 'inherited' && permission.source_role && (
                                <span className="text-xs text-gray-500">
                                  from {permission.source_role}
                                </span>
                              )}
                            </div>
                          </div>
                        </div>

                        <div className="flex items-center">
                          <span className={`inline-flex items-center px-2 py-1 rounded-full text-xs font-medium ${
                            permission.is_granted
                              ? 'bg-green-100 text-green-800'
                              : 'bg-red-100 text-red-800'
                          }`}>
                            {permission.is_granted ? '✓ Granted' : '✗ Denied'}
                          </span>
                        </div>
                      </div>
                    ))}
                  </div>
                </div>
              </div>
            ))}
          </div>
        )}
      </div>
    </div>
  );
};

export default EffectivePermissionsView;