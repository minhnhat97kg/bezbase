import React, { useState, useEffect, useCallback } from 'react';
import { rbacService } from '../../services/api';
import Table from '../common/Table';
import Icon from '../common/Icons';

interface PermissionManagerProps {
  roles: any[];
  onRefresh: () => void;
}

const PermissionManager: React.FC<PermissionManagerProps> = ({ roles, onRefresh }) => {
  const [permissions, setPermissions] = useState([]);
  const [availablePermissions, setAvailablePermissions] = useState([]);
  const [loading, setLoading] = useState(false);
  const [error, setError] = useState(null);
  const [showAddForm, setShowAddForm] = useState(false);
  const [newPermission, setNewPermission] = useState({
    role: '',
    permission: '',
  });

  // Pagination state
  const [currentPage, setCurrentPage] = useState(1);
  const [pageSize, setPageSize] = useState(10);
  const [totalPages, setTotalPages] = useState(0);
  const [totalPermissions, setTotalPermissions] = useState(0);

  // Filter state
  const [filters, setFilters] = useState({
    role: '',
    resource: '',
    action: '',
    permission: '',
  });

  // Sort state
  const [sortField, setSortField] = useState('');
  const [sortOrder, setSortOrder] = useState<'asc' | 'desc'>('asc');

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

  // Handle escape key for modal
  useEffect(() => {
    const handleEscapeKey = (event) => {
      if (event.key === 'Escape' && showAddForm) {
        setShowAddForm(false);
      }
    };

    if (showAddForm) {
      document.addEventListener('keydown', handleEscapeKey);
    }

    return () => {
      document.removeEventListener('keydown', handleEscapeKey);
    };
  }, [showAddForm]);

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
      await rbacService.addPermission({
        role: newPermission.role,
        resource: selectedPermission.resource,
        action: selectedPermission.action,
      });

      setNewPermission({ role: '', permission: '' });
      setShowAddForm(false);
      await fetchPermissions();
      onRefresh();
    } catch (err) {
      setError(err.response?.data?.message || 'Failed to add permission');
    } finally {
      setLoading(false);
    }
  };

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

  const handleFilterChange = (field, value) => {
    setFilters(prev => ({ ...prev, [field]: value }));
    setCurrentPage(1); // Reset to first page when filtering
  };

  const handleSort = (field) => {
    if (sortField === field) {
      setSortOrder(sortOrder === 'asc' ? 'desc' : 'asc');
    } else {
      setSortField(field);
      setSortOrder('asc');
    }
    setCurrentPage(1);
  };

  const handlePageChange = (newPage) => {
    setCurrentPage(newPage);
  };

  const handlePageSizeChange = (newPageSize) => {
    setPageSize(newPageSize);
    setCurrentPage(1);
  };

  // Table configuration
  const columns = [
    {
      key: 'role',
      header: 'Role',
      sortable: true,
      nowrap: true,
      render: (value) => (
        <span className="text-sm font-medium text-gray-900 dark:text-white">
          {value}
        </span>
      )
    },
    {
      key: 'permission',
      header: 'Permission',
      sortable: true,
      nowrap: true,
      render: (value, row) => (
        <span className="text-sm text-gray-500 dark:text-gray-300">
          {value}
        </span>
      )
    },
    {
      key: 'resource',
      header: 'Resource',
      sortable: true,
      nowrap: true,
      render: (value) => (
        <span className="text-sm text-gray-500 dark:text-gray-300">
          {value}
        </span>
      )
    },
    {
      key: 'action',
      header: 'Action',
      sortable: true,
      nowrap: true,
      render: (value) => (
        <span className="text-sm text-gray-500 dark:text-gray-300">
          {value}
        </span>
      )
    },
    {
      key: 'actions',
      header: 'Actions',
      align: 'right',
      nowrap: true,
      render: (value, row) => (
        <button
          onClick={() => handleRemovePermission(row)}
          className="text-red-600 hover:text-red-900 dark:text-red-400 dark:hover:text-red-300"
          disabled={loading}
        >
          <Icon name="delete" />
        </button>
      )
    }
  ];

  const getPermissionBadge = (resource, action) => {
    const isWildcard = resource === '*' || action === '*';
    return (
      <span className={`inline-flex items-center px-2.5 py-0.5 rounded-full text-xs font-medium ${isWildcard
        ? 'bg-yellow-100 text-yellow-800 dark:bg-yellow-800 dark:text-yellow-100'
        : 'bg-green-100 text-green-800 dark:bg-green-800 dark:text-green-100'
        }`}>
        {action} on {resource}
      </span>
    );
  };

  // Pagination configuration
  const paginationConfig = {
    currentPage,
    totalPages,
    pageSize,
    total: totalPermissions,
    pageSizeOptions: [5, 10, 25, 50]
  };

  // Group permissions by resource for dropdown
  const renderGroupedPermissions = () => {
    // Group permissions by resource
    const groupedPermissions = availablePermissions.reduce((groups, permission) => {
      const resource = permission.resource;
      if (!groups[resource]) {
        groups[resource] = [];
      }
      groups[resource].push(permission);
      return groups;
    }, {});

    // Sort resources alphabetically
    const sortedResources = Object.keys(groupedPermissions).sort();

    // Action priority for sorting (common CRUD order)
    const actionPriority = {
      'create': 1,
      'read': 2,
      'update': 3,
      'delete': 4,
      'export': 5,
      'restore': 6
    };

    return sortedResources.map(resource => (
      <optgroup key={resource} label={resource.charAt(0).toUpperCase() + resource.slice(1)}>
        {groupedPermissions[resource]
          .sort((a, b) => {
            const aPriority = actionPriority[a.action] || 999;
            const bPriority = actionPriority[b.action] || 999;
            return aPriority - bPriority;
          })
          .map(permission => (
            <option key={permission.permission} value={permission.permission}>
              {permission.permission}
            </option>
          ))}
      </optgroup>
    ));
  };

  return (
    <div className="space-y-6">
      {error && (
        <div className="bg-red-50 border border-red-200 text-red-700 px-4 py-3 rounded-md">
          {error}
        </div>
      )}

      {/* Add Permission Button */}
      <div className="flex justify-end">
        <button
          onClick={() => setShowAddForm(true)}
          className="bg-blue-600 hover:bg-blue-700 text-white px-4 py-2 rounded-md text-sm font-medium inline-flex items-center"
        >
          <Icon name="plus" className="w-4 h-4 mr-2" />
          Add Permission
        </button>
      </div>

      {/* Add Permission Modal */}
      {showAddForm && (
        <div className="fixed inset-0 z-50 overflow-y-auto">
          <div className="flex items-end justify-center min-h-screen pt-4 px-4 pb-20 text-center sm:block sm:p-0">
            {/* Background overlay */}
            <div 
              className="fixed inset-0 bg-gray-500 bg-opacity-75 transition-opacity"
              onClick={() => setShowAddForm(false)}
            ></div>

            {/* Modal positioning */}
            <span className="hidden sm:inline-block sm:align-middle sm:h-screen">&#8203;</span>

            {/* Modal content */}
            <div className="inline-block align-bottom bg-white dark:bg-gray-800 rounded-lg text-left overflow-hidden shadow-xl transform transition-all sm:my-8 sm:align-middle sm:max-w-lg sm:w-full">
              <div className="bg-white dark:bg-gray-800 px-4 pt-5 pb-4 sm:p-6 sm:pb-4">
                <div className="sm:flex sm:items-start">
                  <div className="mt-3 text-center sm:mt-0 sm:text-left w-full">
                    <h3 className="text-lg leading-6 font-medium text-gray-900 dark:text-white mb-4">
                      Add New Permission
                    </h3>
                    <form onSubmit={handleAddPermission} className="space-y-4">
                      <div>
                        <label className="block text-sm font-medium text-gray-700 dark:text-gray-300 mb-1">
                          Role
                        </label>
                        <select
                          value={newPermission.role}
                          onChange={(e) => setNewPermission({ ...newPermission, role: e.target.value })}
                          className="block w-full px-3 py-2 border border-gray-300 dark:border-gray-600 rounded-md shadow-sm focus:outline-none focus:ring-blue-500 focus:border-blue-500 dark:bg-gray-700 dark:text-white"
                          required
                        >
                          <option value="">Select role...</option>
                          {roles.map((role) => (
                            <option key={role.id} value={role.name}>
                              {role.display_name} ({role.name})
                            </option>
                          ))}
                        </select>
                      </div>
                      <div>
                        <label className="block text-sm font-medium text-gray-700 dark:text-gray-300 mb-1">
                          Permission
                        </label>
                        <select
                          value={newPermission.permission}
                          onChange={(e) => setNewPermission({ ...newPermission, permission: e.target.value })}
                          className="block w-full px-3 py-2 border border-gray-300 dark:border-gray-600 rounded-md shadow-sm focus:outline-none focus:ring-blue-500 focus:border-blue-500 dark:bg-gray-700 dark:text-white"
                          required
                        >
                          <option value="">Select permission...</option>
                          {renderGroupedPermissions()}
                        </select>
                      </div>
                    </form>
                  </div>
                </div>
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

      {/* Permissions Table */}
      <div className="bg-white dark:bg-gray-800 shadow overflow-hidden sm:rounded-md">
        <div className="px-4 py-5 sm:p-6">
          <div className="flex justify-between items-center mb-4">
            <h3 className="text-lg font-medium text-gray-900 dark:text-white">
              All Permissions ({totalPermissions})
            </h3>
          </div>

          {/* Filters */}
          <div className="grid grid-cols-4 gap-4 mb-4">
            <div>
              <label className="block text-sm font-medium text-gray-700 dark:text-gray-300 mb-1">
                Filter by Role
              </label>
              <input
                type="text"
                value={filters.role}
                onChange={(e) => handleFilterChange('role', e.target.value)}
                placeholder="Filter roles..."
                className="block w-full px-3 py-2 border border-gray-300 dark:border-gray-600 rounded-md shadow-sm focus:outline-none focus:ring-blue-500 focus:border-blue-500 dark:bg-gray-700 dark:text-white"
              />
            </div>
            <div>
              <label className="block text-sm font-medium text-gray-700 dark:text-gray-300 mb-1">
                Filter by Permission
              </label>
              <input
                type="text"
                value={filters.permission}
                onChange={(e) => handleFilterChange('permission', e.target.value)}
                placeholder="Filter permissions..."
                className="block w-full px-3 py-2 border border-gray-300 dark:border-gray-600 rounded-md shadow-sm focus:outline-none focus:ring-blue-500 focus:border-blue-500 dark:bg-gray-700 dark:text-white"
              />
            </div>
            <div>
              <label className="block text-sm font-medium text-gray-700 dark:text-gray-300 mb-1">
                Filter by Resource
              </label>
              <input
                type="text"
                value={filters.resource}
                onChange={(e) => handleFilterChange('resource', e.target.value)}
                placeholder="Filter resources..."
                className="block w-full px-3 py-2 border border-gray-300 dark:border-gray-600 rounded-md shadow-sm focus:outline-none focus:ring-blue-500 focus:border-blue-500 dark:bg-gray-700 dark:text-white"
              />
            </div>
            <div>
              <label className="block text-sm font-medium text-gray-700 dark:text-gray-300 mb-1">
                Filter by Action
              </label>
              <input
                type="text"
                value={filters.action}
                onChange={(e) => handleFilterChange('action', e.target.value)}
                placeholder="Filter actions..."
                className="block w-full px-3 py-2 border border-gray-300 dark:border-gray-600 rounded-md shadow-sm focus:outline-none focus:ring-blue-500 focus:border-blue-500 dark:bg-gray-700 dark:text-white"
              />
            </div>
          </div>

          {/* Table */}
          <Table
            columns={columns}
            data={permissions}
            loading={loading}
            emptyMessage="No permissions found"
            onSort={handleSort}
            sortField={sortField}
            sortOrder={sortOrder}
            pagination={paginationConfig}
            onPageChange={handlePageChange}
            onPageSizeChange={handlePageSizeChange}
          />
        </div>
      </div>
    </div>
  );
};

export default PermissionManager;
