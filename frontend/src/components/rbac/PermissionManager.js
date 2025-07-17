import React, { useState, useEffect, useCallback } from 'react';
import useResourceActionOptions from '../../hooks/useResourceActionOptions';
import { rbacService } from '../../services/api';
import Table from '../common/Table';
import Icon from '../common/Icons';

const PermissionManager = ({ roles, onRefresh }) => {
  const { resources: resourceOptions, actions: actionOptions, loading: optionsLoading } = useResourceActionOptions();
  const [permissions, setPermissions] = useState([]);
  const [loading, setLoading] = useState(false);
  const [error, setError] = useState(null);
  const [showAddForm, setShowAddForm] = useState(false);
  const [newPermission, setNewPermission] = useState({
    role: '',
    resource: '',
    action: '',
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
  });

  // Sort state
  const [sortField, setSortField] = useState('');
  const [sortOrder, setSortOrder] = useState('asc');

  // Common resources and actions for quick selection
  // Deprecated: commonResources and commonActions

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
  }, [fetchPermissions]);

  const handleAddPermission = async (e) => {
    e.preventDefault();
    if (!newPermission.role || !newPermission.resource || !newPermission.action) {
      return;
    }

    try {
      setLoading(true);
      await rbacService.addPermission({
        role: newPermission.role,
        resource: newPermission.resource,
        action: newPermission.action,
      });

      setNewPermission({ role: '', resource: '', action: '' });
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
      key: 'permission',
      header: 'Permission',
      nowrap: true,
      render: (value, row) => getPermissionBadge(row.resource, row.action)
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

  return (
    <div className="space-y-6">
      {error && (
        <div className="bg-red-50 border border-red-200 text-red-700 px-4 py-3 rounded-md">
          {error}
        </div>
      )}

      {/* Add Permission Form */}
      <div className="bg-white dark:bg-gray-800 shadow rounded-lg">
        <div className="px-4 py-5 sm:p-6">
          <div className="flex justify-between items-center mb-4">
            <h3 className="text-lg font-medium text-gray-900 dark:text-white">
              Add New Permission
            </h3>
            <button
              onClick={() => setShowAddForm(!showAddForm)}
              className="bg-blue-600 hover:bg-blue-700 text-white px-4 py-2 rounded-md text-sm font-medium"
            >
              {showAddForm ? 'Cancel' : 'Add Permission'}
            </button>
          </div>

          {showAddForm && (
            <form onSubmit={handleAddPermission} className="space-y-4">
              <div className="grid grid-cols-3 gap-4">
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
                    Resource
                  </label>
                <select
                  value={newPermission.resource}
                  onChange={(e) => setNewPermission({ ...newPermission, resource: e.target.value })}
                  className="block w-full mt-2 px-3 py-2 border border-gray-300 dark:border-gray-600 rounded-md shadow-sm focus:outline-none focus:ring-blue-500 focus:border-blue-500 dark:bg-gray-700 dark:text-white"
                  required
                >
                  <option value="">Select resource...</option>
                  {resourceOptions.map(resource => (
                    <option key={resource.name} value={resource.name}>{resource.name}</option>
                  ))}
                </select>
                </div>
                <div>
                  <label className="block text-sm font-medium text-gray-700 dark:text-gray-300 mb-1">
                    Action
                  </label>
                <select
                  value={newPermission.action}
                  onChange={(e) => setNewPermission({ ...newPermission, action: e.target.value })}
                  className="block w-full mt-2 px-3 py-2 border border-gray-300 dark:border-gray-600 rounded-md shadow-sm focus:outline-none focus:ring-blue-500 focus:border-blue-500 dark:bg-gray-700 dark:text-white"
                  required
                >
                  <option value="">Select action...</option>
                  {actionOptions.map(action => (
                    <option key={action.name} value={action.name}>{action.name}</option>
                  ))}
                </select>
                </div>
              </div>
              <div className="flex justify-end">
                <button
                  type="submit"
                  disabled={loading}
                  className="px-4 py-2 text-sm font-medium text-white bg-blue-600 border border-transparent rounded-md hover:bg-blue-700 focus:outline-none focus:ring-2 focus:ring-offset-2 focus:ring-blue-500 disabled:opacity-50"
                >
                  {loading ? 'Adding...' : 'Add Permission'}
                </button>
              </div>
            </form>
          )}
        </div>
      </div>

      {/* Permissions Table */}
      <div className="bg-white dark:bg-gray-800 shadow overflow-hidden sm:rounded-md">
        <div className="px-4 py-5 sm:p-6">
          <div className="flex justify-between items-center mb-4">
            <h3 className="text-lg font-medium text-gray-900 dark:text-white">
              All Permissions ({totalPermissions})
            </h3>
          </div>

          {/* Filters */}
          <div className="grid grid-cols-3 gap-4 mb-4">
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
