import React, { useState, useEffect } from 'react';
import { rbacService } from '../services/api';
import RolesList from '../components/rbac/RolesList';
import RoleForm from '../components/rbac/RoleForm';
import PermissionManager from '../components/rbac/PermissionManager';
import UserRoleAssignment from '../components/rbac/UserRoleAssignment';

const RoleManagement = () => {
  // Set page title
  useEffect(() => {
    document.title = 'Role Management - BezBase';
  }, []);
  
  const [roles, setRoles] = useState([]);
  const [loading, setLoading] = useState(true);
  const [error, setError] = useState(null);
  const [showForm, setShowForm] = useState(false);
  const [selectedRole, setSelectedRole] = useState(null);
  const [activeTab, setActiveTab] = useState('roles');
  
  const fetchRoles = async () => {
    try {
      setLoading(true);
      const response = await rbacService.getRoles();
      setRoles(response.data);
    } catch (err) {
      setError(err.response?.data?.message || 'Failed to fetch roles');
    } finally {
      setLoading(false);
    }
  };

  useEffect(() => {
    fetchRoles();
  }, []);

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

  const tabs = [
    { id: 'roles', name: 'Roles', icon: 'shield' },
    { id: 'permissions', name: 'Permissions', icon: 'key' },
    { id: 'assignments', name: 'User Assignments', icon: 'users' },
  ];

  const getTabIcon = (iconName) => {
    const icons = {
      shield: (
        <svg className="w-5 h-5" fill="none" stroke="currentColor" viewBox="0 0 24 24">
          <path strokeLinecap="round" strokeLinejoin="round" strokeWidth={2} d="M9 12l2 2 4-4m5.618-4.016A11.955 11.955 0 0112 2.944a11.955 11.955 0 01-8.618 3.04A12.02 12.02 0 003 9c0 5.591 3.824 10.29 9 11.622 5.176-1.332 9-6.03 9-11.622 0-1.042-.133-2.052-.382-3.016z" />
        </svg>
      ),
      key: (
        <svg className="w-5 h-5" fill="none" stroke="currentColor" viewBox="0 0 24 24">
          <path strokeLinecap="round" strokeLinejoin="round" strokeWidth={2} d="M15 7a2 2 0 012 2m4 0a6 6 0 01-7.743 5.743L11 17H9v2H7v2H4a1 1 0 01-1-1v-2.586a1 1 0 01.293-.707l5.964-5.964A6 6 0 1121 9z" />
        </svg>
      ),
      users: (
        <svg className="w-5 h-5" fill="none" stroke="currentColor" viewBox="0 0 24 24">
          <path strokeLinecap="round" strokeLinejoin="round" strokeWidth={2} d="M12 4.354a4 4 0 110 5.292M15 21H3v-1a6 6 0 0112 0v1zm0 0h6v-1a6 6 0 00-9-5.197m13.5-9a2.25 2.25 0 11-4.5 0 2.25 2.25 0 014.5 0z" />
        </svg>
      ),
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
              className={`py-2 px-1 border-b-2 font-medium text-sm flex items-center space-x-2 ${
                activeTab === tab.id
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
            <RolesList
              roles={roles}
              onEdit={handleEditRole}
              onDelete={handleDeleteRole}
              loading={loading}
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