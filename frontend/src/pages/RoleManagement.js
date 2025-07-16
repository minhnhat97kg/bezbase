import React, { useState, useEffect } from 'react';
import { rbacService } from '../services/api';
import RolesList from '../components/rbac/RolesList';
import RoleForm from '../components/rbac/RoleForm';
import PermissionManager from '../components/rbac/PermissionManager';
import UserRoleAssignment from '../components/rbac/UserRoleAssignment';
import Icon from '../components/common/Icons';

const RoleManagement = () => {
  // Set page title
  useEffect(() => {
    document.title = 'Role Management - BezBase';
  }, []);

  const [roles, setRoles] = useState([]);
  const [pagination, setPagination] = useState({
    currentPage: 1,
    pageSize: 10,
    total: 0,
    totalPages: 0
  });
  const [filters, setFilters] = useState({
    search: '',
    status: '',
    is_system: '',
    sort: 'created_at',
    order: 'desc'
  });
  const [loading, setLoading] = useState(true);
  const [error, setError] = useState(null);
  const [showForm, setShowForm] = useState(false);
  const [selectedRole, setSelectedRole] = useState(null);
  const [activeTab, setActiveTab] = useState('roles');

  const fetchRoles = async (page = pagination.currentPage) => {
    try {
      setLoading(true);
      const params = {
        page,
        page_size: pagination.pageSize,
        ...filters
      };

      // Remove empty filters
      Object.keys(params).forEach(key => {
        if (params[key] === '' || params[key] === null || params[key] === undefined) {
          delete params[key];
        }
      });

      const response = await rbacService.getRoles(params);

      setRoles(response.data.data);
      setPagination({
        currentPage: response.data.page,
        pageSize: response.data.page_size,
        total: response.data.total_items,
        totalPages: response.data.total_pages
      });
    } catch (err) {
      setError(err.response?.data?.message || 'Failed to fetch roles');
    } finally {
      setLoading(false);
    }
  };

  useEffect(() => {
    fetchRoles();
  }, []);

  useEffect(() => {
    fetchRoles(1);
  }, [filters]);

  const handleCreateRole = () => {
    setSelectedRole(null);
    setShowForm(true);
  };

  const handleEditRole = (role) => {
    setSelectedRole(role);
    setShowForm(true);
  };

  const handleFormClose = () => {
    setShowForm(false);
    setSelectedRole(null);
  };

  const handleFormSuccess = () => {
    fetchRoles();
    handleFormClose();
  };

  const handleDeleteRole = async (roleName) => {
    if (window.confirm(`Are you sure you want to delete the role "${roleName}"?`)) {
      try {
        await rbacService.deleteRole(roleName);
        fetchRoles();
      } catch (err) {
        alert(err.response?.data?.message || 'Failed to delete role');
      }
    }
  };

  const handlePageChange = (page) => {
    fetchRoles(page);
  };

  const handlePageSizeChange = (pageSize) => {
    setPagination(prev => ({ ...prev, pageSize, currentPage: 1 }));
    fetchRoles(1);
  };

  const handleSort = (field) => {
    const newOrder = filters.sort === field && filters.order === 'asc' ? 'desc' : 'asc';
    setFilters(prev => ({ ...prev, sort: field, order: newOrder }));
  };

  const handleFilterChange = (key, value) => {
    setFilters(prev => ({ ...prev, [key]: value }));
  };

  const tabs = [
    { id: 'roles', name: 'Roles', icon: 'shield' },
    { id: 'permissions', name: 'Permissions', icon: 'key' },
    { id: 'assignments', name: 'User Assignments', icon: 'users' },
  ];

  const getTabIcon = (iconName) => {
    const icons = {
      shield: <Icon name="shield" />,
      key: <Icon name="key" />,
      users: <Icon name="users" />,
    };
    return icons[iconName];
  };

  if (loading) {
    return (
      <div className="flex justify-center items-center h-64">
        <div className="animate-spin rounded-full h-12 w-12 border-b-2 border-blue-600"></div>
      </div>
    );
  }

  return (
    <div className="max-w-7xl mx-auto">
      <div className="mb-8">
        <h1 className="text-3xl font-bold text-gray-900 dark:text-white">Role Management</h1>
        <p className="mt-2 text-gray-600 dark:text-gray-400">
          Manage roles, permissions, and user assignments
        </p>
      </div>

      {error && (
        <div className="mb-6 bg-red-50 border border-red-200 text-red-700 px-4 py-3 rounded-md">
          {error}
        </div>
      )}

      {/* Tab Navigation */}
      <div className="mb-6 border-b border-gray-200 dark:border-gray-700">
        <nav className="-mb-px flex space-x-8">
          {tabs.map((tab) => (
            <button
              key={tab.id}
              onClick={() => setActiveTab(tab.id)}
              className={`py-2 px-1 border-b-2 font-medium text-sm flex items-center space-x-2 ${activeTab === tab.id
                ? 'border-blue-500 text-blue-600 dark:text-blue-400'
                : 'border-transparent text-gray-500 hover:text-gray-700 hover:border-gray-300 dark:text-gray-400 dark:hover:text-gray-300'
                }`}
            >
              {getTabIcon(tab.icon)}
              <span>{tab.name}</span>
            </button>
          ))}
        </nav>
      </div>

      {/* Tab Content */}
      <div className="bg-white dark:bg-gray-800 rounded-lg shadow">
        {activeTab === 'roles' && (
          <div className="p-6">
            <div className="flex justify-between items-center mb-6">
              <h2 className="text-lg font-medium text-gray-900 dark:text-white">Roles</h2>
              <button
                onClick={handleCreateRole}
                className="bg-blue-600 hover:bg-blue-700 text-white px-4 py-2 rounded-md text-sm font-medium"
              >
                Create Role
              </button>
            </div>

            {/* Search and Filter Controls */}
            <div className="mb-6 space-y-4">
              <div className="flex flex-col sm:flex-row gap-4">
                <div className="flex-1">
                  <input
                    type="text"
                    placeholder="Search roles..."
                    value={filters.search}
                    onChange={(e) => handleFilterChange('search', e.target.value)}
                    className="w-full px-3 py-2 border border-gray-300 dark:border-gray-600 rounded-md dark:bg-gray-700 dark:text-white"
                  />
                </div>
                <div className="flex gap-2">
                  <select
                    value={filters.status}
                    onChange={(e) => handleFilterChange('status', e.target.value)}
                    className="px-3 py-2 border border-gray-300 dark:border-gray-600 rounded-md dark:bg-gray-700 dark:text-white"
                  >
                    <option value="">All Status</option>
                    <option value="active">Active</option>
                    <option value="inactive">Inactive</option>
                  </select>
                  <select
                    value={filters.is_system}
                    onChange={(e) => handleFilterChange('is_system', e.target.value)}
                    className="px-3 py-2 border border-gray-300 dark:border-gray-600 rounded-md dark:bg-gray-700 dark:text-white"
                  >
                    <option value="">All Types</option>
                    <option value="true">System</option>
                    <option value="false">Custom</option>
                  </select>
                </div>
              </div>
            </div>

            <RolesList
              roles={roles}
              onEdit={handleEditRole}
              onDelete={handleDeleteRole}
              loading={loading}
              pagination={pagination}
              onPageChange={handlePageChange}
              onPageSizeChange={handlePageSizeChange}
              onSort={handleSort}
              sortField={filters.sort}
              sortOrder={filters.order}
            />
          </div>
        )}

        {activeTab === 'permissions' && (
          <div className="p-6">
            <h2 className="text-lg font-medium text-gray-900 dark:text-white mb-6">Permission Management</h2>
            <PermissionManager roles={roles} onRefresh={fetchRoles} />
          </div>
        )}

        {activeTab === 'assignments' && (
          <div className="p-6">
            <h2 className="text-lg font-medium text-gray-900 dark:text-white mb-6">User Role Assignments</h2>
            <UserRoleAssignment roles={roles} onRefresh={fetchRoles} />
          </div>
        )}
      </div>

      {/* Role Form Modal */}
      {showForm && (
        <RoleForm
          role={selectedRole}
          onClose={handleFormClose}
          onSuccess={handleFormSuccess}
        />
      )}
    </div>
  );
};

export default RoleManagement;
