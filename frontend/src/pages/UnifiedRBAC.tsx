import React, { useState, useEffect } from 'react';
import { useTranslation } from 'react-i18next';
import TabLayout from '../components/common/TabLayout';

// Import simplified components
import SimplifiedRoleManager from '../components/rbac/SimplifiedRoleManager';
import EnhancedPermissionManager from '../components/rbac/EnhancedPermissionManager';
import UserRoleAssignment from '../components/rbac/UserRoleAssignment';
import SimpleRoleHierarchy from '../components/rbac/SimpleRoleHierarchy';
import EffectivePermissionsView from '../components/rbac/EffectivePermissionsView';
import { rbacService } from '../services/rbacService';

const UnifiedRBAC: React.FC = () => {
  const { t } = useTranslation();
  const [activeTab, setActiveTab] = useState('roles');
  const [roles, setRoles] = useState([]);
  const [selectedRole, setSelectedRole] = useState(null);

  const loadRoles = async () => {
    try {
      const response = await rbacService.getRoles();
      const rolesData = Array.isArray(response.data) ? response.data : [];
      setRoles(rolesData);
    } catch (error) {
      console.error('Failed to load roles:', error);
      setRoles([]);
    }
  };

  const handleRefresh = () => {
    loadRoles();
  };

  useEffect(() => {
    loadRoles();
  }, []);

  useEffect(() => {
    console.log('Roles updated:', roles);
  }, [roles]);

  const tabs = [
    { id: 'roles', label: 'Roles & Permissions', icon: '🎭' },
    { id: 'hierarchy', label: 'Role Hierarchy', icon: '🌳' },
    { id: 'assignments', label: 'User Assignments', icon: '👥' },
    { id: 'permissions', label: 'Effective Permissions', icon: '🔐' },
    { id: 'settings', label: 'Settings', icon: '⚙️' },
  ];

  return (
    <div className="space-y-6">
      {/* Header */}
      <div className="bg-white shadow rounded-lg">
        <div className="px-4 py-5 sm:p-6">
          <div className="flex items-center justify-between">
            <div>
              <h1 className="text-2xl font-bold text-gray-900">
                {t('rbac.title', 'Role-Based Access Control')}
              </h1>
              <p className="mt-1 text-sm text-gray-500">
                Comprehensive role and permission management for the system
              </p>
            </div>
          </div>

          {/* Quick Stats */}
          <div className="mt-6 grid grid-cols-1 gap-5 sm:grid-cols-5">
            <div className="bg-blue-50 rounded-lg p-4">
              <div className="flex items-center">
                <div className="flex-shrink-0">
                  <span className="text-blue-600 text-xl">🎭</span>
                </div>
                <div className="ml-3">
                  <p className="text-sm font-medium text-blue-800">Active Roles</p>
                  <p className="text-lg font-semibold text-blue-900">{roles.length}</p>
                </div>
              </div>
            </div>

            <div className="bg-green-50 rounded-lg p-4">
              <div className="flex items-center">
                <div className="flex-shrink-0">
                  <span className="text-green-600 text-xl">🔑</span>
                </div>
                <div className="ml-3">
                  <p className="text-sm font-medium text-green-800">Permissions</p>
                  <p className="text-lg font-semibold text-green-900">-</p>
                </div>
              </div>
            </div>

            <div className="bg-orange-50 rounded-lg p-4">
              <div className="flex items-center">
                <div className="flex-shrink-0">
                  <span className="text-orange-600 text-xl">🌳</span>
                </div>
                <div className="ml-3">
                  <p className="text-sm font-medium text-orange-800">Root Roles</p>
                  <p className="text-lg font-semibold text-orange-900">
                    {roles.filter(r => !r.parent_role_id).length}
                  </p>
                </div>
              </div>
            </div>

            <div className="bg-purple-50 rounded-lg p-4">
              <div className="flex items-center">
                <div className="flex-shrink-0">
                  <span className="text-purple-600 text-xl">👥</span>
                </div>
                <div className="ml-3">
                  <p className="text-sm font-medium text-purple-800">Users Assigned</p>
                  <p className="text-lg font-semibold text-purple-900">-</p>
                </div>
              </div>
            </div>

            <div className="bg-yellow-50 rounded-lg p-4">
              <div className="flex items-center">
                <div className="flex-shrink-0">
                  <span className="text-yellow-600 text-xl">📊</span>
                </div>
                <div className="ml-3">
                  <p className="text-sm font-medium text-yellow-800">Inherited Roles</p>
                  <p className="text-lg font-semibold text-yellow-900">
                    {roles.filter(r => r.parent_role_id).length}
                  </p>
                </div>
              </div>
            </div>
          </div>
        </div>
      </div>

      <TabLayout
        tabs={tabs}
        activeTab={activeTab}
        onTabChange={setActiveTab}
      >
        {/* Roles & Permissions Tab */}
        {activeTab === 'roles' && (
          <div className="space-y-6">
            <SimplifiedRoleManager />
            <EnhancedPermissionManager roles={roles} onRefresh={handleRefresh} />
          </div>
        )}

        {/* Role Hierarchy Tab */}
        {activeTab === 'hierarchy' && (
          <div className="space-y-6">
            <SimpleRoleHierarchy 
              roles={roles}
              onRoleUpdate={loadRoles}
            />
          </div>
        )}

        {/* User Assignments Tab */}
        {activeTab === 'assignments' && (
          <div className="space-y-6">
            <UserRoleAssignment roles={roles} onRefresh={handleRefresh} />
          </div>
        )}

        {/* Effective Permissions Tab */}
        {activeTab === 'permissions' && (
          <div className="space-y-6">
            {selectedRole ? (
              <EffectivePermissionsView 
                roleId={selectedRole.id}
                role={selectedRole}
              />
            ) : (
              <div className="bg-white shadow rounded-lg p-6">
                <div className="text-center py-8">
                  <span className="text-gray-400 text-4xl mb-4 block">🔐</span>
                  <h3 className="text-lg font-medium text-gray-900 mb-2">
                    Select a Role
                  </h3>
                  <p className="text-gray-500 mb-4">
                    Choose a role from the hierarchy view to see its effective permissions
                  </p>
                  <button
                    onClick={() => setActiveTab('hierarchy')}
                    className="bg-indigo-600 hover:bg-indigo-700 text-white px-4 py-2 rounded-md text-sm"
                  >
                    View Role Hierarchy
                  </button>
                </div>
              </div>
            )}
          </div>
        )}

        {/* Settings Tab */}
        {activeTab === 'settings' && (
          <div className="space-y-6">
            <BasicSettingsView />
          </div>
        )}
      </TabLayout>
    </div>
  );
};

// Basic Settings View Component
const BasicSettingsView: React.FC = () => {
  const handleCreateTemplateRole = async (template: string) => {
    const templates = {
      admin: {
        name: 'admin',
        display_name: 'Administrator',
        description: 'Full system administration access'
      },
      manager: {
        name: 'manager',
        display_name: 'Manager',
        description: 'Team management and oversight capabilities'
      },
      user: {
        name: 'user',
        display_name: 'User',
        description: 'Basic user access to the system'
      }
    };

    try {
      const roleData = templates[template];
      await rbacService.createRole(roleData);
      alert(`${roleData.display_name} role created successfully!`);
    } catch (error: any) {
      console.error('Failed to create role from template:', error);
      alert(`Failed to create role: ${error.response?.data?.message || error.message || 'Unknown error'}`);
    }
  };

  return (
    <div className="space-y-6">
      {/* Quick Role Templates */}
      <div className="bg-white shadow rounded-lg">
        <div className="px-4 py-5 sm:px-6 border-b border-gray-200">
          <h3 className="text-lg font-medium text-gray-900">Quick Role Templates</h3>
          <p className="mt-1 text-sm text-gray-500">
            Create common roles quickly with pre-defined permissions
          </p>
        </div>
        <div className="p-6">
          <div className="grid grid-cols-1 md:grid-cols-3 gap-4">
            <button 
              onClick={() => handleCreateTemplateRole('admin')}
              className="p-4 border-2 border-dashed border-gray-300 rounded-lg hover:border-indigo-500 hover:bg-indigo-50 transition-colors"
            >
              <div className="text-center">
                <span className="text-2xl mb-2 block">👑</span>
                <h4 className="font-medium text-gray-900">Admin</h4>
                <p className="text-sm text-gray-500">Full system access</p>
              </div>
            </button>
            <button 
              onClick={() => handleCreateTemplateRole('manager')}
              className="p-4 border-2 border-dashed border-gray-300 rounded-lg hover:border-indigo-500 hover:bg-indigo-50 transition-colors"
            >
              <div className="text-center">
                <span className="text-2xl mb-2 block">👨‍💼</span>
                <h4 className="font-medium text-gray-900">Manager</h4>
                <p className="text-sm text-gray-500">Team management</p>
              </div>
            </button>
            <button 
              onClick={() => handleCreateTemplateRole('user')}
              className="p-4 border-2 border-dashed border-gray-300 rounded-lg hover:border-indigo-500 hover:bg-indigo-50 transition-colors"
            >
              <div className="text-center">
                <span className="text-2xl mb-2 block">👤</span>
                <h4 className="font-medium text-gray-900">User</h4>
                <p className="text-sm text-gray-500">Basic access</p>
              </div>
            </button>
          </div>
        </div>
      </div>

      {/* System Settings */}
      <div className="bg-white shadow rounded-lg">
        <div className="px-4 py-5 sm:px-6 border-b border-gray-200">
          <h3 className="text-lg font-medium text-gray-900">System Settings</h3>
          <p className="mt-1 text-sm text-gray-500">
            Configure system-wide RBAC behavior
          </p>
        </div>
        <div className="p-6 space-y-4">
          <div className="flex items-center justify-between">
            <div>
              <h4 className="text-sm font-medium text-gray-900">Require role for new users</h4>
              <p className="text-sm text-gray-500">Automatically assign a default role to new users</p>
            </div>
            <input type="checkbox" className="h-4 w-4 text-indigo-600 border-gray-300 rounded" />
          </div>
          <div className="flex items-center justify-between">
            <div>
              <h4 className="text-sm font-medium text-gray-900">Enable role hierarchy</h4>
              <p className="text-sm text-gray-500">Allow roles to inherit permissions from parent roles</p>
            </div>
            <div className="flex items-center">
              <input type="checkbox" className="h-4 w-4 text-indigo-600 border-gray-300 rounded" defaultChecked disabled />
              <span className="ml-2 inline-flex items-center px-2 py-0.5 rounded-full text-xs font-medium bg-green-100 text-green-800">
                Active
              </span>
            </div>
          </div>
          <div className="flex items-center justify-between">
            <div>
              <h4 className="text-sm font-medium text-gray-900">Audit permission changes</h4>
              <p className="text-sm text-gray-500">Log all permission and role modifications</p>
            </div>
            <input type="checkbox" className="h-4 w-4 text-indigo-600 border-gray-300 rounded" defaultChecked />
          </div>
        </div>
      </div>
    </div>
  );
};

export default UnifiedRBAC;
