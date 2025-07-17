import React, { useState, useEffect } from 'react';
import { rbacService } from '../services/api';
import RolesList from '../components/rbac/RolesList';
import RoleForm from '../components/rbac/RoleForm';
import PermissionManager from '../components/rbac/PermissionManager';
import UserRoleAssignment from '../components/rbac/UserRoleAssignment';
import ResourcesList from '../components/rbac/ResourcesList';
import ActionsList from '../components/rbac/ActionsList';
import Icon from '../components/common/Icons';
import TabLayout from '../components/common/TabLayout';

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
    { id: 'resources', name: 'Resources', icon: 'settings' },
    { id: 'actions', name: 'Actions', icon: 'star' },
  ];

  if (loading) {
    return (
      <div className="flex justify-center items-center h-64">
        <div className="animate-spin rounded-full h-12 w-12 border-b-2 border-primary-600"></div>
      </div>
    );
  }

  return (
    <>
    <TabLayout
      title="Role Management"
      subtitle="Manage roles, permissions, and user assignments"
      showTabs={true}
      tabs={tabs}
      activeTab={activeTab}
      onTabChange={setActiveTab}
      error={error}
      onErrorDismiss={() => setError('')}
      className="max-w-7xl mx-auto"
    >
        {activeTab === 'roles' && (
          <div>
            <div className="flex justify-between items-center mb-6">
              <h2 className="text-lg font-medium text-gray-900">Roles</h2>
              <button
                onClick={handleCreateRole}
                className="bg-primary-600 hover:bg-primary-700 text-white px-4 py-2 rounded-md text-sm font-medium"
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
                    className="w-full px-4 py-3 border border-gray-300 rounded-lg focus:ring-2 focus:ring-primary-500 focus:border-transparent transition-colors"
                  />
                </div>
                <div className="flex gap-2">
                  <select
                    value={filters.status}
                    onChange={(e) => handleFilterChange('status', e.target.value)}
                    className="px-4 py-3 border border-gray-300 rounded-lg focus:ring-2 focus:ring-primary-500 focus:border-transparent transition-colors"
                  >
                    <option value="">All Status</option>
                    <option value="active">Active</option>
                    <option value="inactive">Inactive</option>
                  </select>
                  <select
                    value={filters.is_system}
                    onChange={(e) => handleFilterChange('is_system', e.target.value)}
                    className="px-4 py-3 border border-gray-300 rounded-lg focus:ring-2 focus:ring-primary-500 focus:border-transparent transition-colors"
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
          <div>
            <h2 className="text-lg font-medium text-gray-900 mb-6">Permission Management</h2>
            <PermissionManager roles={roles} onRefresh={fetchRoles} />
          </div>
        )}

        {activeTab === 'assignments' && (
          <div>
            <h2 className="text-lg font-medium text-gray-900 mb-6">User Role Assignments</h2>
            <UserRoleAssignment roles={roles} onRefresh={fetchRoles} />
          </div>
        )}

        {activeTab === 'resources' && (
          <div>
            <h2 className="text-lg font-medium text-gray-900 mb-6">Resource Management</h2>
            <ResourcesList />
          </div>
        )}

        {activeTab === 'actions' && (
          <div>
            <h2 className="text-lg font-medium text-gray-900 mb-6">Action Management</h2>
            <ActionsList />
          </div>
        )}
    </TabLayout>

    {/* Role Form Modal */}
    {showForm && (
      <RoleForm
        role={selectedRole}
        onClose={handleFormClose}
        onSuccess={handleFormSuccess}
      />
    )}
    </>
  );
};

export default RoleManagement;
