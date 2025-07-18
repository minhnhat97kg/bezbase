import React, { useState } from 'react';
import { useOrganization } from '../context/OrganizationContext';
import TabLayout from '../components/common/TabLayout';
import AdvancedRoleManager from '../components/rbac/AdvancedRoleManager';
import ContextualPermissionManager from '../components/rbac/ContextualPermissionManager';

const AdvancedRBAC: React.FC = () => {
  const { currentOrganization } = useOrganization();
  const [activeTab, setActiveTab] = useState('roles');

  const tabs = [
    { id: 'roles', label: 'Hierarchical Roles', icon: '🎭' },
    { id: 'permissions', label: 'Contextual Permissions', icon: '🔑' },
    { id: 'templates', label: 'Role Templates', icon: '📋' },
  ];

  return (
    <div className="space-y-6">
      {/* Header */}
      <div className="bg-white shadow">
        <div className="px-4 py-5 sm:p-6">
          <div className="flex items-center justify-between">
            <div>
              <h1 className="text-2xl font-bold text-gray-900">Advanced RBAC</h1>
              <p className="mt-1 text-sm text-gray-500">
                Manage hierarchical roles and contextual permissions for{' '}
                {currentOrganization ? currentOrganization.name : 'Global System'}
              </p>
            </div>
            {currentOrganization && (
              <div className="text-right">
                <div className="text-sm text-gray-500">Current Organization</div>
                <div className="font-medium text-gray-900">{currentOrganization.name}</div>
                <div className="text-xs text-gray-500">ID: {currentOrganization.id}</div>
              </div>
            )}
          </div>
        </div>
      </div>

      <TabLayout
        tabs={tabs}
        activeTab={activeTab}
        onTabChange={setActiveTab}
      >
        {/* Hierarchical Roles Tab */}
        {activeTab === 'roles' && (
          <div className="space-y-6">
            <AdvancedRoleManager />
          </div>
        )}

        {/* Contextual Permissions Tab */}
        {activeTab === 'permissions' && (
          <div className="space-y-6">
            <ContextualPermissionManager />
          </div>
        )}

        {/* Role Templates Tab */}
        {activeTab === 'templates' && (
          <div className="space-y-6">
            <RoleTemplatesView />
          </div>
        )}
      </TabLayout>
    </div>
  );
};

// Role Templates View Component
const RoleTemplatesView: React.FC = () => {
  return (
    <div className="bg-white shadow rounded-lg">
      <div className="px-4 py-5 sm:px-6 border-b border-gray-200">
        <h3 className="text-lg font-medium text-gray-900">Role Templates</h3>
        <p className="mt-1 text-sm text-gray-500">
          Predefined role configurations for quick role creation
        </p>
      </div>
      
      <div className="px-4 py-5 sm:p-6">
        <div className="grid grid-cols-1 gap-6 sm:grid-cols-2 lg:grid-cols-3">
          {/* System Templates */}
          <div className="border border-gray-200 rounded-lg p-6">
            <div className="flex items-center space-x-3 mb-4">
              <div className="flex-shrink-0">
                <div className="w-10 h-10 bg-red-100 rounded-full flex items-center justify-center">
                  <span className="text-red-600 font-semibold">🔒</span>
                </div>
              </div>
              <div>
                <h4 className="text-lg font-medium text-gray-900">System Templates</h4>
                <p className="text-sm text-gray-500">Global system roles</p>
              </div>
            </div>
            <div className="space-y-3">
              <TemplateCard
                name="Organization Admin"
                description="Full administrative access within an organization"
                permissions={['*:*']}
                context="organization"
              />
              <TemplateCard
                name="Global Admin"
                description="System-wide administrative access"
                permissions={['*:*']}
                context="global"
              />
            </div>
          </div>

          {/* Business Templates */}
          <div className="border border-gray-200 rounded-lg p-6">
            <div className="flex items-center space-x-3 mb-4">
              <div className="flex-shrink-0">
                <div className="w-10 h-10 bg-blue-100 rounded-full flex items-center justify-center">
                  <span className="text-blue-600 font-semibold">💼</span>
                </div>
              </div>
              <div>
                <h4 className="text-lg font-medium text-gray-900">Business Templates</h4>
                <p className="text-sm text-gray-500">Business-focused roles</p>
              </div>
            </div>
            <div className="space-y-3">
              <TemplateCard
                name="Team Lead"
                description="Team management and project oversight"
                permissions={['users:read', 'users:update', 'projects:*']}
                context="team"
              />
              <TemplateCard
                name="Project Manager"
                description="Project management and coordination"
                permissions={['projects:*', 'reports:read']}
                context="project"
              />
            </div>
          </div>

          {/* Basic Templates */}
          <div className="border border-gray-200 rounded-lg p-6">
            <div className="flex items-center space-x-3 mb-4">
              <div className="flex-shrink-0">
                <div className="w-10 h-10 bg-gray-100 rounded-full flex items-center justify-center">
                  <span className="text-gray-600 font-semibold">👤</span>
                </div>
              </div>
              <div>
                <h4 className="text-lg font-medium text-gray-900">Basic Templates</h4>
                <p className="text-sm text-gray-500">Standard user roles</p>
              </div>
            </div>
            <div className="space-y-3">
              <TemplateCard
                name="Viewer"
                description="Read-only access to assigned resources"
                permissions={['*:read']}
                context="resource"
              />
              <TemplateCard
                name="Contributor"
                description="Create and edit assigned resources"
                permissions={['*:read', '*:create', '*:update']}
                context="resource"
              />
            </div>
          </div>
        </div>

        <div className="mt-8 text-center">
          <div className="bg-gray-50 rounded-lg p-6">
            <h4 className="text-lg font-medium text-gray-900 mb-2">How to Use Templates</h4>
            <div className="text-sm text-gray-600 space-y-2">
              <p>1. Go to the "Hierarchical Roles" tab</p>
              <p>2. Click "Create Role from Template"</p>
              <p>3. Select a template and optionally customize the name</p>
              <p>4. The role will be created with predefined permissions and context</p>
            </div>
          </div>
        </div>
      </div>
    </div>
  );
};

// Template Card Component
interface TemplateCardProps {
  name: string;
  description: string;
  permissions: string[];
  context: string;
}

const TemplateCard: React.FC<TemplateCardProps> = ({ name, description, permissions, context }) => {
  const getContextColor = (context: string) => {
    switch (context) {
      case 'global':
        return 'bg-red-100 text-red-800';
      case 'organization':
        return 'bg-blue-100 text-blue-800';
      case 'team':
      case 'project':
        return 'bg-green-100 text-green-800';
      default:
        return 'bg-gray-100 text-gray-800';
    }
  };

  return (
    <div className="bg-gray-50 rounded-lg p-4">
      <div className="flex items-center justify-between mb-2">
        <h5 className="font-medium text-gray-900">{name}</h5>
        <span className={`inline-flex px-2 py-1 text-xs font-semibold rounded-full ${getContextColor(context)}`}>
          {context}
        </span>
      </div>
      <p className="text-sm text-gray-600 mb-3">{description}</p>
      <div className="space-y-1">
        <div className="text-xs font-medium text-gray-700">Permissions:</div>
        <div className="flex flex-wrap gap-1">
          {permissions.map((permission, index) => (
            <span
              key={index}
              className="inline-flex px-2 py-1 text-xs bg-white text-gray-700 rounded border"
            >
              {permission}
            </span>
          ))}
        </div>
      </div>
    </div>
  );
};

export default AdvancedRBAC;