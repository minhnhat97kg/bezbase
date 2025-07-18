import React, { useState, useEffect } from 'react';
import { useOrganization } from '../../context/OrganizationContext';
import { 
  advancedRbacService, 
  ContextualPermission, 
  Role, 
  CreateContextualPermissionRequest 
} from '../../services/advancedRbacService';
import Table from '../common/Table';

const ContextualPermissionManager: React.FC = () => {
  const { currentOrganization, getCurrentOrgId } = useOrganization();
  const [permissions, setPermissions] = useState<ContextualPermission[]>([]);
  const [roles, setRoles] = useState<Role[]>([]);
  const [selectedUserId, setSelectedUserId] = useState<number | null>(null);
  const [showCreateModal, setShowCreateModal] = useState(false);
  const [isLoading, setIsLoading] = useState(false);

  useEffect(() => {
    loadRoles();
  }, [currentOrganization]);

  const loadRoles = async () => {
    try {
      const orgId = getCurrentOrgId();
      const response = await advancedRbacService.getRolesByOrganization(orgId || undefined);
      setRoles(Array.isArray(response.data) ? response.data : []);
    } catch (error) {
      console.error('Failed to load roles:', error);
      setRoles([]);
    }
  };

  const loadUserPermissions = async (userId: number) => {
    setIsLoading(true);
    try {
      const orgId = getCurrentOrgId();
      const response = await advancedRbacService.getEffectivePermissions(userId, orgId || undefined);
      setPermissions(response.data);
      setSelectedUserId(userId);
    } catch (error) {
      console.error('Failed to load user permissions:', error);
    } finally {
      setIsLoading(false);
    }
  };

  const handleCreatePermission = async (data: CreateContextualPermissionRequest) => {
    try {
      await advancedRbacService.createContextualPermission(data);
      setShowCreateModal(false);
      if (selectedUserId) {
        await loadUserPermissions(selectedUserId);
      }
    } catch (error) {
      console.error('Failed to create contextual permission:', error);
    }
  };

  const permissionColumns = [
    {
      key: 'permission_key',
      label: 'Permission',
      render: (permission: ContextualPermission) => (
        <div>
          <div className="font-medium text-gray-900">
            {advancedRbacService.formatPermissionKey(permission)}
          </div>
          <div className="text-sm text-gray-500">
            {permission.resource}:{permission.action}
          </div>
        </div>
      ),
    },
    {
      key: 'context',
      label: 'Context',
      render: (permission: ContextualPermission) => (
        permission.context_type && permission.context_value ? (
          <div className="text-sm">
            <span className="font-medium">{permission.context_type}:</span>
            <span className="ml-1">{permission.context_value}</span>
          </div>
        ) : (
          <span className="text-sm text-gray-400">Global</span>
        )
      ),
    },
    {
      key: 'status',
      label: 'Status',
      render: (permission: ContextualPermission) => (
        <span className={`inline-flex px-2 py-1 text-xs font-semibold rounded-full ${advancedRbacService.getPermissionColor(permission)}`}>
          {permission.is_granted ? 'Granted' : 'Denied'}
        </span>
      ),
    },
    {
      key: 'role',
      label: 'Role',
      render: (permission: ContextualPermission) => (
        permission?.role ? (
          <div className="text-sm">
            <div className="font-medium">{permission.role?.display_name || permission.role?.name || 'Unknown Role'}</div>
            <div className="text-gray-500">Level {permission.role?.hierarchy_level ?? 0}</div>
          </div>
        ) : (
          <span className="text-sm text-gray-400">Unknown</span>
        )
      ),
    },
  ];

  return (
    <div className="space-y-6">
      {/* Header */}
      <div className="bg-white shadow rounded-lg">
        <div className="px-4 py-5 sm:px-6 border-b border-gray-200">
          <div className="flex items-center justify-between">
            <div>
              <h3 className="text-lg font-medium text-gray-900">Contextual Permission Manager</h3>
              <p className="mt-1 text-sm text-gray-500">
                Manage fine-grained permissions with organizational context
              </p>
            </div>
            <button
              onClick={() => setShowCreateModal(true)}
              className="inline-flex items-center px-4 py-2 border border-transparent text-sm font-medium rounded-md shadow-sm text-white bg-indigo-600 hover:bg-indigo-700"
              disabled={roles.length === 0}
            >
              Add Permission
            </button>
          </div>
        </div>

        {/* User Selection */}
        <div className="px-4 py-5 sm:p-6">
          <div className="max-w-lg">
            <label className="block text-sm font-medium text-gray-700 mb-2">
              Select User to View Permissions
            </label>
            <div className="flex space-x-3">
              <input
                type="number"
                placeholder="Enter User ID"
                className="flex-1 border-gray-300 rounded-md shadow-sm focus:ring-indigo-500 focus:border-indigo-500"
                onKeyPress={(e) => {
                  if (e.key === 'Enter') {
                    const value = parseInt((e.target as HTMLInputElement).value);
                    if (value) {
                      loadUserPermissions(value);
                    }
                  }
                }}
              />
              <button
                onClick={() => {
                  const input = document.querySelector('input[type="number"]') as HTMLInputElement;
                  const value = parseInt(input.value);
                  if (value) {
                    loadUserPermissions(value);
                  }
                }}
                className="inline-flex items-center px-4 py-2 border border-gray-300 text-sm font-medium rounded-md text-gray-700 bg-white hover:bg-gray-50"
              >
                Load Permissions
              </button>
            </div>
            <p className="mt-1 text-xs text-gray-500">
              Enter a user ID to view their effective permissions in the current organization context
            </p>
          </div>
        </div>
      </div>

      {/* Context Information */}
      {currentOrganization && (
        <div className="bg-blue-50 border border-blue-200 rounded-lg p-4">
          <div className="flex items-center">
            <div className="flex-shrink-0">
              <svg className="h-5 w-5 text-blue-400" fill="currentColor" viewBox="0 0 20 20">
                <path fillRule="evenodd" d="M18 10a8 8 0 11-16 0 8 8 0 0116 0zm-7-4a1 1 0 11-2 0 1 1 0 012 0zM9 9a1 1 0 000 2v3a1 1 0 001 1h1a1 1 0 100-2v-3a1 1 0 00-1-1H9z" clipRule="evenodd" />
              </svg>
            </div>
            <div className="ml-3">
              <h3 className="text-sm font-medium text-blue-800">Organization Context</h3>
              <div className="mt-2 text-sm text-blue-700">
                <p>
                  Currently viewing permissions for <strong>{currentOrganization.name}</strong> (ID: {currentOrganization.id})
                </p>
                <p className="text-xs mt-1">
                  Permissions shown are scoped to this organization and include inherited global permissions.
                </p>
              </div>
            </div>
          </div>
        </div>
      )}

      {/* Permissions Table */}
      {selectedUserId && (
        <div className="bg-white shadow rounded-lg">
          <div className="px-4 py-5 sm:px-6 border-b border-gray-200">
            <h3 className="text-lg font-medium text-gray-900">
              Effective Permissions for User ID: {selectedUserId}
            </h3>
            <p className="mt-1 text-sm text-gray-500">
              {permissions.length} permissions found (including inherited)
            </p>
          </div>
          
          <Table
            data={permissions}
            columns={permissionColumns}
            loading={isLoading}
            emptyMessage="No permissions found for this user"
          />
        </div>
      )}

      {/* Available Roles */}
      <div className="bg-white shadow rounded-lg">
        <div className="px-4 py-5 sm:px-6 border-b border-gray-200">
          <h3 className="text-lg font-medium text-gray-900">Available Roles</h3>
          <p className="mt-1 text-sm text-gray-500">
            Roles available in the current organization context
          </p>
        </div>
        
        <div className="px-4 py-5 sm:p-6">
          <div className="grid grid-cols-1 gap-4 sm:grid-cols-2 lg:grid-cols-3">
            {Array.isArray(roles) && roles.map((role) => (
              <div key={role.id} className="border border-gray-200 rounded-lg p-4">
                <div className="flex items-center space-x-3">
                  <span className="text-lg">{role?.is_system ? '🔒' : '👤'}</span>
                  <div className="flex-1 min-w-0">
                    <div className="font-medium text-gray-900 truncate">{role?.display_name || role?.name || 'Unknown Role'}</div>
                    <div className="text-sm text-gray-500 truncate">{role?.name || 'unknown'}</div>
                  </div>
                  <span className={`inline-flex px-2 py-1 text-xs font-semibold rounded-full ${role?.is_system ? 'bg-red-100 text-red-800' : 'bg-blue-100 text-blue-800'}`}>
                    L{role?.hierarchy_level ?? 0}
                  </span>
                </div>
                {role?.description && (
                  <p className="mt-2 text-sm text-gray-600">{role.description}</p>
                )}
              </div>
            ))}
          </div>
          
          {(!Array.isArray(roles) || roles.length === 0) && (
            <div className="text-center py-6">
              <p className="text-sm text-gray-500">No roles found in current context</p>
            </div>
          )}
        </div>
      </div>

      {/* Create Permission Modal */}
      {showCreateModal && (
        <CreateContextualPermissionModal
          roles={roles}
          onSubmit={handleCreatePermission}
          onClose={() => setShowCreateModal(false)}
        />
      )}
    </div>
  );
};

// Create Contextual Permission Modal
interface CreateContextualPermissionModalProps {
  roles: Role[];
  onSubmit: (data: CreateContextualPermissionRequest) => void;
  onClose: () => void;
}

const CreateContextualPermissionModal: React.FC<CreateContextualPermissionModalProps> = ({
  roles,
  onSubmit,
  onClose,
}) => {
  const [formData, setFormData] = useState<CreateContextualPermissionRequest>({
    role_id: 0,
    resource: '',
    action: '',
    context_type: '',
    context_value: '',
    is_granted: true,
  });

  const handleSubmit = (e: React.FormEvent) => {
    e.preventDefault();
    onSubmit(formData);
  };

  const commonResources = ['users', 'projects', 'reports', 'settings', 'admin', 'dashboard'];
  const commonActions = ['create', 'read', 'update', 'delete', 'export', 'all'];
  const contextTypes = ['organization', 'project', 'department', 'team'];

  return (
    <div className="fixed inset-0 z-50 overflow-y-auto">
      <div className="flex items-center justify-center min-h-screen pt-4 px-4 pb-20 text-center sm:block sm:p-0">
        <div className="fixed inset-0 bg-gray-500 bg-opacity-75 transition-opacity" onClick={onClose}></div>

        <div className="inline-block align-bottom bg-white rounded-lg text-left overflow-hidden shadow-xl transform transition-all sm:my-8 sm:align-middle sm:max-w-lg sm:w-full">
          <form onSubmit={handleSubmit}>
            <div className="bg-white px-4 pt-5 pb-4 sm:p-6 sm:pb-4">
              <h3 className="text-lg font-medium text-gray-900 mb-4">Create Contextual Permission</h3>
              
              <div className="space-y-4">
                <div>
                  <label className="block text-sm font-medium text-gray-700">Role</label>
                  <select
                    value={formData.role_id}
                    onChange={(e) => setFormData({ ...formData, role_id: parseInt(e.target.value) })}
                    className="mt-1 block w-full border-gray-300 rounded-md shadow-sm focus:ring-indigo-500 focus:border-indigo-500"
                    required
                  >
                    <option value={0}>Select a role</option>
                    {Array.isArray(roles) && roles.map((role) => (
                      <option key={role.id} value={role.id}>
                        {role?.display_name || role?.name || 'Unknown Role'} (Level {role?.hierarchy_level ?? 0})
                      </option>
                    ))}
                  </select>
                </div>

                <div className="grid grid-cols-2 gap-4">
                  <div>
                    <label className="block text-sm font-medium text-gray-700">Resource</label>
                    <input
                      type="text"
                      list="resources"
                      value={formData.resource}
                      onChange={(e) => setFormData({ ...formData, resource: e.target.value })}
                      className="mt-1 block w-full border-gray-300 rounded-md shadow-sm focus:ring-indigo-500 focus:border-indigo-500"
                      required
                    />
                    <datalist id="resources">
                      {commonResources.map(resource => (
                        <option key={resource} value={resource} />
                      ))}
                    </datalist>
                  </div>

                  <div>
                    <label className="block text-sm font-medium text-gray-700">Action</label>
                    <input
                      type="text"
                      list="actions"
                      value={formData.action}
                      onChange={(e) => setFormData({ ...formData, action: e.target.value })}
                      className="mt-1 block w-full border-gray-300 rounded-md shadow-sm focus:ring-indigo-500 focus:border-indigo-500"
                      required
                    />
                    <datalist id="actions">
                      {commonActions.map(action => (
                        <option key={action} value={action} />
                      ))}
                    </datalist>
                  </div>
                </div>

                <div className="grid grid-cols-2 gap-4">
                  <div>
                    <label className="block text-sm font-medium text-gray-700">Context Type (Optional)</label>
                    <select
                      value={formData.context_type}
                      onChange={(e) => setFormData({ ...formData, context_type: e.target.value })}
                      className="mt-1 block w-full border-gray-300 rounded-md shadow-sm focus:ring-indigo-500 focus:border-indigo-500"
                    >
                      <option value="">No specific context</option>
                      {contextTypes.map(type => (
                        <option key={type} value={type}>{type}</option>
                      ))}
                    </select>
                  </div>

                  <div>
                    <label className="block text-sm font-medium text-gray-700">Context Value</label>
                    <input
                      type="text"
                      value={formData.context_value}
                      onChange={(e) => setFormData({ ...formData, context_value: e.target.value })}
                      className="mt-1 block w-full border-gray-300 rounded-md shadow-sm focus:ring-indigo-500 focus:border-indigo-500"
                      placeholder="e.g., organization ID or '*' for all"
                      disabled={!formData.context_type}
                    />
                  </div>
                </div>

                <div>
                  <label className="flex items-center">
                    <input
                      type="checkbox"
                      checked={formData.is_granted}
                      onChange={(e) => setFormData({ ...formData, is_granted: e.target.checked })}
                      className="rounded border-gray-300 text-indigo-600 shadow-sm focus:border-indigo-300 focus:ring focus:ring-indigo-200 focus:ring-opacity-50"
                    />
                    <span className="ml-2 text-sm text-gray-700">Grant Permission</span>
                  </label>
                  <p className="mt-1 text-xs text-gray-500">
                    Uncheck to explicitly deny this permission
                  </p>
                </div>
              </div>
            </div>

            <div className="bg-gray-50 px-4 py-3 sm:px-6 sm:flex sm:flex-row-reverse">
              <button
                type="submit"
                disabled={formData.role_id === 0 || !formData.resource || !formData.action}
                className="w-full inline-flex justify-center rounded-md border border-transparent shadow-sm px-4 py-2 bg-indigo-600 text-base font-medium text-white hover:bg-indigo-700 focus:outline-none focus:ring-2 focus:ring-offset-2 focus:ring-indigo-500 sm:ml-3 sm:w-auto sm:text-sm disabled:opacity-50"
              >
                Create Permission
              </button>
              <button
                type="button"
                onClick={onClose}
                className="mt-3 w-full inline-flex justify-center rounded-md border border-gray-300 shadow-sm px-4 py-2 bg-white text-base font-medium text-gray-700 hover:bg-gray-50 focus:outline-none focus:ring-2 focus:ring-offset-2 focus:ring-indigo-500 sm:mt-0 sm:ml-3 sm:w-auto sm:text-sm"
              >
                Cancel
              </button>
            </div>
          </form>
        </div>
      </div>
    </div>
  );
};

export default ContextualPermissionManager;