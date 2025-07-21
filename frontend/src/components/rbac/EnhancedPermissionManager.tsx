import React, { useState, useEffect, useCallback } from 'react';
import { rbacService } from '../../services/rbacService';
import Table from '../common/Table';
import Icon from '../common/Icons';

interface EnhancedPermissionManagerProps {
  roles: any[];
  onRefresh: () => void;
}

const EnhancedPermissionManager: React.FC<EnhancedPermissionManagerProps> = ({ roles, onRefresh }) => {
  const [permissions, setPermissions] = useState([]);
  const [availablePermissions, setAvailablePermissions] = useState([]);
  const [loading, setLoading] = useState(false);
  const [error, setError] = useState(null);
  const [showAddForm, setShowAddForm] = useState(false);

  
  const [newPermission, setNewPermission] = useState({
    role: '',
    permission: '',
    context_type: '',
    context_value: '',
    is_granted: true,
  });

  const [currentPage, setCurrentPage] = useState(1);
  const [pageSize, setPageSize] = useState(10);
  const [totalPages, setTotalPages] = useState(0);
  const [totalPermissions, setTotalPermissions] = useState(0);

  const [filters, setFilters] = useState({
    role: '',
    resource: '',
    action: '',
    permission: '',
    context_type: '',
    is_granted: 'all',
  });

  const [sortField] = useState('');
  const [sortOrder] = useState<'asc' | 'desc'>('asc');

  const fetchAvailablePermissions = useCallback(async () => {
    try {
      const response = await rbacService.getAvailablePermissions();
      setAvailablePermissions(response.data);
    } catch (err) {
      console.error('Failed to fetch available permissions:', err);
      setAvailablePermissions([]);
    }
  }, []);

  const fetchPermissions = useCallback(async () => {
    try {
      setLoading(true);
      setError(null);

      const params = {
        page: currentPage,
        page_size: pageSize,
        ...(filters.role && { role: filters.role }),
        ...(filters.resource && { resource: filters.resource }),
        ...(filters.action && { action: filters.action }),
        ...(filters.permission && { permission: filters.permission }),
        ...(filters.context_type && { context_type: filters.context_type }),
        ...(filters.is_granted !== 'all' && { is_granted: filters.is_granted === 'true' }),
        ...(sortField && { sort: sortField, order: sortOrder }),
      };

      const response = await rbacService.getPermissions(params);
      const data = response.data;

      setPermissions(data.data);
      setTotalPermissions(data.total_items || data.total || 0);
      setTotalPages(data.total_pages || 0);
    } catch (err) {
      setError(err.response?.data?.message || 'Failed to fetch permissions');
      setPermissions([]);
    } finally {
      setLoading(false);
    }
  }, [currentPage, pageSize, filters, sortField, sortOrder]);

  useEffect(() => {
    fetchPermissions();
    fetchAvailablePermissions();
  }, [fetchPermissions, fetchAvailablePermissions]);

  const handleAddPermission = async (e) => {
    if (e) e.preventDefault();
    if (!newPermission.role || !newPermission.permission) {
      setError('Please select both role and permission');
      return;
    }

    const selectedPermission = availablePermissions.find(p => p.permission === newPermission.permission);
    if (!selectedPermission) {
      setError('Invalid permission selected');
      return;
    }

    try {
      setLoading(true);
      setError(null);
      
      const permissionData = {
        role: newPermission.role,
        resource: selectedPermission.resource,
        action: selectedPermission.action,
        ...(newPermission.context_type && {
          context_type: newPermission.context_type,
          context_value: newPermission.context_value,
        }),
        is_granted: newPermission.is_granted,
      };

      await rbacService.addPermission(permissionData);
      setNewPermission({ role: '', permission: '', context_type: '', context_value: '', is_granted: true });
      setShowAddForm(false);
      await fetchPermissions();
      onRefresh();
    } catch (err) {
      setError(err.response?.data?.message || 'Failed to add permission');
    } finally {
      setLoading(false);
    }
  };


  const columns = [
    {
      key: 'role',
      header: 'Role',
      sortable: true,
      render: (value) => (
        <div className="flex items-center space-x-2">
          <span className="text-sm font-medium text-gray-900 dark:text-white">{value}</span>
          {roles.find(r => r.name === value)?.is_system && (
            <span className="inline-flex items-center px-2 py-0.5 rounded text-xs font-medium bg-red-100 text-red-800">
              System
            </span>
          )}
        </div>
      )
    },
    {
      key: 'resource',
      header: 'Resource',
      sortable: true,
      render: (value) => (
        <span className={`inline-flex items-center px-2.5 py-0.5 rounded-full text-xs font-medium ${
          value === '*' ? 'bg-yellow-100 text-yellow-800' : 'bg-blue-100 text-blue-800'
        }`}>
          {value}
        </span>
      )
    },
    {
      key: 'action',
      header: 'Action',
      sortable: true,
      render: (value) => (
        <span className={`inline-flex items-center px-2.5 py-0.5 rounded-full text-xs font-medium ${
          value === '*' ? 'bg-yellow-100 text-yellow-800' : 'bg-green-100 text-green-800'
        }`}>
          {value}
        </span>
      )
    },
    {
      key: 'context',
      header: 'Context',
      render: (value, row) => (
        row.context_type ? (
          <div className="text-sm">
            <span className="font-medium text-gray-700">{row.context_type}:</span>
            <span className="ml-1 text-gray-500">{row.context_value || '*'}</span>
          </div>
        ) : (
          <span className="text-sm text-gray-400">Global</span>
        )
      )
    },
    {
      key: 'is_granted',
      header: 'Status',
      sortable: true,
      render: (value) => (
        <span className={`inline-flex items-center px-2 py-1 rounded-full text-xs font-medium ${
          value ? 'bg-green-100 text-green-800' : 'bg-red-100 text-red-800'
        }`}>
          {value ? 'Granted' : 'Denied'}
        </span>
      )
    },
    {
      key: 'actions',
      header: 'Actions',
      align: 'right',
      render: (value, row) => (
        <div className="flex space-x-2">
          <button
            onClick={() => {/* Handle edit */}}
            className="text-blue-600 hover:text-blue-900"
            title="Edit"
          >
            <Icon name="edit" className="w-4 h-4" />
          </button>
          <button
            onClick={() => handleRemovePermission(row)}
            className="text-red-600 hover:text-red-900"
            title="Delete"
          >
            <Icon name="delete" className="w-4 h-4" />
          </button>
        </div>
      )
    }
  ];

  const handleRemovePermission = async (permission) => {
    if (!window.confirm(`Remove permission "${permission.action}" on "${permission.resource}" from role "${permission.role}"?`)) {
      return;
    }

    try {
      setLoading(true);
      await rbacService.removePermission({
        role: permission.role,
        resource: permission.resource,
        action: permission.action,
      });
      await fetchPermissions();
      onRefresh();
    } catch (err) {
      setError(err.response?.data?.message || 'Failed to remove permission');
    } finally {
      setLoading(false);
    }
  };


  const paginationConfig = {
    currentPage,
    totalPages,
    pageSize,
    total: totalPermissions,
    pageSizeOptions: [5, 10, 25, 50, 100]
  };

  return (
    <div className="space-y-6">
      {error && (
        <div className="bg-red-50 border border-red-200 text-red-700 px-4 py-3 rounded-md">
          {error}
        </div>
      )}

      {/* Header with Actions */}
      <div className="flex flex-col sm:flex-row justify-between items-start sm:items-center space-y-4 sm:space-y-0">
        <div>
          <h2 className="text-2xl font-bold text-gray-900 dark:text-white">Permission Management</h2>
          <p className="text-sm text-gray-500 dark:text-gray-400">
            Manage role permissions and access controls
          </p>
        </div>

        <div className="flex justify-end">
          <button
            onClick={() => setShowAddForm(true)}
            className="bg-blue-600 hover:bg-blue-700 text-white px-4 py-2 rounded-md text-sm font-medium flex items-center"
          >
            <Icon name="plus" className="w-4 h-4 mr-2" />
            Add Permission
          </button>
        </div>
      </div>

      {/* Advanced Filters */}
      <div className="bg-white dark:bg-gray-800 rounded-lg border p-4">
        <div className="grid grid-cols-1 md:grid-cols-2 lg:grid-cols-6 gap-4">
          <div>
            <label className="block text-sm font-medium text-gray-700 dark:text-gray-300 mb-1">Role</label>
            <input
              type="text"
              value={filters.role}
              onChange={(e) => setFilters({ ...filters, role: e.target.value })}
              placeholder="Filter by role..."
              className="w-full px-3 py-2 border border-gray-300 dark:border-gray-600 rounded-md text-sm focus:outline-none focus:ring-2 focus:ring-blue-500 dark:bg-gray-700 dark:text-white"
            />
          </div>
          <div>
            <label className="block text-sm font-medium text-gray-700 dark:text-gray-300 mb-1">Resource</label>
            <input
              type="text"
              value={filters.resource}
              onChange={(e) => setFilters({ ...filters, resource: e.target.value })}
              placeholder="Filter by resource..."
              className="w-full px-3 py-2 border border-gray-300 dark:border-gray-600 rounded-md text-sm focus:outline-none focus:ring-2 focus:ring-blue-500 dark:bg-gray-700 dark:text-white"
            />
          </div>
          <div>
            <label className="block text-sm font-medium text-gray-700 dark:text-gray-300 mb-1">Action</label>
            <input
              type="text"
              value={filters.action}
              onChange={(e) => setFilters({ ...filters, action: e.target.value })}
              placeholder="Filter by action..."
              className="w-full px-3 py-2 border border-gray-300 dark:border-gray-600 rounded-md text-sm focus:outline-none focus:ring-2 focus:ring-blue-500 dark:bg-gray-700 dark:text-white"
            />
          </div>
          <div>
            <label className="block text-sm font-medium text-gray-700 dark:text-gray-300 mb-1">Context</label>
            <select
              value={filters.context_type}
              onChange={(e) => setFilters({ ...filters, context_type: e.target.value })}
              className="w-full px-3 py-2 border border-gray-300 dark:border-gray-600 rounded-md text-sm focus:outline-none focus:ring-2 focus:ring-blue-500 dark:bg-gray-700 dark:text-white"
            >
              <option value="">All contexts</option>
              <option value="project">Project</option>
              <option value="department">Department</option>
              <option value="team">Team</option>
              <option value="resource">Resource</option>
            </select>
          </div>
          <div>
            <label className="block text-sm font-medium text-gray-700 dark:text-gray-300 mb-1">Status</label>
            <select
              value={filters.is_granted}
              onChange={(e) => setFilters({ ...filters, is_granted: e.target.value })}
              className="w-full px-3 py-2 border border-gray-300 dark:border-gray-600 rounded-md text-sm focus:outline-none focus:ring-2 focus:ring-blue-500 dark:bg-gray-700 dark:text-white"
            >
              <option value="all">All</option>
              <option value="true">Granted</option>
              <option value="false">Denied</option>
            </select>
          </div>
          <div className="flex items-end">
            <button
              onClick={() => setFilters({ role: '', resource: '', action: '', permission: '', context_type: '', is_granted: 'all' })}
              className="w-full px-3 py-2 border border-gray-300 dark:border-gray-600 rounded-md text-sm hover:bg-gray-50 dark:hover:bg-gray-700"
            >
              Clear Filters
            </button>
          </div>
        </div>
      </div>

      {/* Permissions Table */}
      <div className="bg-white dark:bg-gray-800 shadow overflow-hidden rounded-lg">
        <div className="px-4 py-5 sm:p-6">
          <div className="flex justify-between items-center mb-4">
            <h3 className="text-lg font-medium text-gray-900 dark:text-white">
              All Permissions ({totalPermissions})
            </h3>
          </div>

          <Table
            columns={columns}
            data={permissions}
            loading={loading}
            emptyMessage="No permissions found"
            pagination={paginationConfig}
            onPageChange={setCurrentPage}
            onPageSizeChange={setPageSize}
          />
        </div>
      </div>

      {/* Enhanced Add Permission Modal */}
      {showAddForm && (
        <div className="fixed inset-0 z-50 overflow-y-auto">
          <div className="flex items-center justify-center min-h-screen pt-4 px-4 pb-20 text-center sm:block sm:p-0">
            <div 
              className="fixed inset-0 bg-gray-500 bg-opacity-75 transition-opacity"
              onClick={() => setShowAddForm(false)}
            />
            
            <div className="inline-block align-bottom bg-white dark:bg-gray-800 rounded-lg text-left overflow-hidden shadow-xl transform transition-all sm:my-8 sm:align-middle sm:max-w-2xl sm:w-full">
              <div className="bg-white dark:bg-gray-800 px-4 pt-5 pb-4 sm:p-6 sm:pb-4">
                <h3 className="text-lg leading-6 font-medium text-gray-900 dark:text-white mb-4">
                  Add New Permission
                </h3>
                <form onSubmit={handleAddPermission} className="space-y-6">
                  <div className="grid grid-cols-1 md:grid-cols-2 gap-4">
                    <div>
                      <label className="block text-sm font-medium text-gray-700 dark:text-gray-300 mb-1">
                        Role *
                      </label>
                      <select
                        value={newPermission.role}
                        onChange={(e) => setNewPermission({ ...newPermission, role: e.target.value })}
                        className="block w-full px-3 py-2 border border-gray-300 dark:border-gray-600 rounded-md shadow-sm focus:outline-none focus:ring-blue-500 focus:border-blue-500 dark:bg-gray-700 dark:text-white"
                        required
                      >
                        <option value="">Select role...</option>
                        {Array.isArray(roles) && roles.map((role) => (
                          <option key={role.id} value={role.name}>
                            {role.display_name} ({role.name})
                          </option>
                        ))}
                      </select>
                    </div>
                    
                    <div>
                      <label className="block text-sm font-medium text-gray-700 dark:text-gray-300 mb-1">
                        Permission *
                      </label>
                      <select
                        value={newPermission.permission}
                        onChange={(e) => setNewPermission({ ...newPermission, permission: e.target.value })}
                        className="block w-full px-3 py-2 border border-gray-300 dark:border-gray-600 rounded-md shadow-sm focus:outline-none focus:ring-blue-500 focus:border-blue-500 dark:bg-gray-700 dark:text-white"
                        required
                      >
                        <option value="">Select permission...</option>
                        {availablePermissions.map(permission => (
                          <option key={permission.permission} value={permission.permission}>
                            {permission.permission}
                          </option>
                        ))}
                      </select>
                    </div>
                  </div>

                  <div className="grid grid-cols-1 md:grid-cols-2 gap-4">
                    <div>
                      <label className="block text-sm font-medium text-gray-700 dark:text-gray-300 mb-1">
                        Context Type
                      </label>
                      <select
                        value={newPermission.context_type}
                        onChange={(e) => setNewPermission({ ...newPermission, context_type: e.target.value })}
                        className="block w-full px-3 py-2 border border-gray-300 dark:border-gray-600 rounded-md shadow-sm focus:outline-none focus:ring-blue-500 focus:border-blue-500 dark:bg-gray-700 dark:text-white"
                      >
                        <option value="">Global (no context)</option>
                        <option value="project">Project</option>
                        <option value="department">Department</option>
                        <option value="team">Team</option>
                        <option value="resource">Resource</option>
                      </select>
                    </div>
                    
                    <div>
                      <label className="block text-sm font-medium text-gray-700 dark:text-gray-300 mb-1">
                        Context Value
                      </label>
                      <input
                        type="text"
                        value={newPermission.context_value}
                        onChange={(e) => setNewPermission({ ...newPermission, context_value: e.target.value })}
                        placeholder="e.g., project-123 or * for all"
                        disabled={!newPermission.context_type}
                        className="block w-full px-3 py-2 border border-gray-300 dark:border-gray-600 rounded-md shadow-sm focus:outline-none focus:ring-blue-500 focus:border-blue-500 dark:bg-gray-700 dark:text-white disabled:opacity-50"
                      />
                    </div>
                  </div>

                  <div className="flex items-center">
                    <input
                      type="checkbox"
                      id="is_granted"
                      checked={newPermission.is_granted}
                      onChange={(e) => setNewPermission({ ...newPermission, is_granted: e.target.checked })}
                      className="h-4 w-4 text-blue-600 focus:ring-blue-500 border-gray-300 rounded"
                    />
                    <label htmlFor="is_granted" className="ml-2 block text-sm text-gray-900 dark:text-white">
                      Grant Permission
                    </label>
                    <span className="ml-2 text-xs text-gray-500">(uncheck to deny)</span>
                  </div>
                </form>
              </div>
              
              <div className="bg-gray-50 dark:bg-gray-700 px-4 py-3 sm:px-6 sm:flex sm:flex-row-reverse">
                <button
                  type="button"
                  onClick={handleAddPermission}
                  disabled={loading}
                  className="w-full inline-flex justify-center rounded-md border border-transparent shadow-sm px-4 py-2 bg-blue-600 text-base font-medium text-white hover:bg-blue-700 focus:outline-none focus:ring-2 focus:ring-offset-2 focus:ring-blue-500 sm:ml-3 sm:w-auto sm:text-sm disabled:opacity-50"
                >
                  {loading ? 'Adding...' : 'Add Permission'}
                </button>
                <button
                  type="button"
                  onClick={() => setShowAddForm(false)}
                  className="mt-3 w-full inline-flex justify-center rounded-md border border-gray-300 shadow-sm px-4 py-2 bg-white text-base font-medium text-gray-700 hover:bg-gray-50 focus:outline-none focus:ring-2 focus:ring-offset-2 focus:ring-indigo-500 sm:mt-0 sm:ml-3 sm:w-auto sm:text-sm dark:bg-gray-600 dark:text-white dark:border-gray-500 dark:hover:bg-gray-500"
                >
                  Cancel
                </button>
              </div>
            </div>
          </div>
        </div>
      )}
    </div>
  );
};

export default EnhancedPermissionManager;