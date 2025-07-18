import React, { useState, useEffect } from 'react';
import { useTranslation } from 'react-i18next';
import { useOrganization } from '../context/OrganizationContext';
import TabLayout from '../components/common/TabLayout';
import { rbacService } from '../services/api';

// Import existing components
import AdvancedRoleManager from '../components/rbac/AdvancedRoleManager';
import ContextualPermissionManager from '../components/rbac/ContextualPermissionManager';
import PermissionManager from '../components/rbac/PermissionManager';
import UserRoleAssignment from '../components/rbac/UserRoleAssignment';

const UnifiedRBAC: React.FC = () => {
  const { t } = useTranslation();
  const { currentOrganization } = useOrganization();
  const [activeTab, setActiveTab] = useState('roles');
  const [roles, setRoles] = useState([]);

  const loadRoles = async () => {
    try {
      const response = await rbacService.getRoles();
      setRoles(response.data.data || []);
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

  const tabs = [
    { id: 'roles', label: 'Roles & Hierarchy', icon: '🎭' },
    { id: 'permissions', label: 'Permission Management', icon: '🔑' },
    { id: 'contextual', label: 'Contextual Permissions', icon: '🌐' },
    { id: 'assignments', label: 'User Assignments', icon: '👥' },
    { id: 'templates', label: 'Role Templates', icon: '📋' },
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
                Comprehensive role and permission management for{' '}
                {currentOrganization ? currentOrganization.name : 'Global System'}
              </p>
            </div>
            {currentOrganization && (
              <div className="text-right">
                <div className="text-sm text-gray-500">Current Context</div>
                <div className="font-medium text-gray-900">{currentOrganization.name}</div>
                <div className="text-xs text-gray-500">
                  Organization ID: {currentOrganization.id}
                </div>
              </div>
            )}
          </div>

          {/* Quick Stats */}
          <div className="mt-6 grid grid-cols-1 gap-5 sm:grid-cols-4">
            <div className="bg-blue-50 rounded-lg p-4">
              <div className="flex items-center">
                <div className="flex-shrink-0">
                  <span className="text-blue-600 text-xl">🎭</span>
                </div>
                <div className="ml-3">
                  <p className="text-sm font-medium text-blue-800">Active Roles</p>
                  <p className="text-lg font-semibold text-blue-900">-</p>
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
                  <p className="text-sm font-medium text-yellow-800">Templates</p>
                  <p className="text-lg font-semibold text-yellow-900">6</p>
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
        {/* Roles & Hierarchy Tab */}
        {activeTab === 'roles' && (
          <div className="space-y-6">
            <AdvancedRoleManager />
          </div>
        )}

        {/* Permission Management Tab */}
        {activeTab === 'permissions' && (
          <div className="space-y-6">
            <PermissionManager roles={roles} onRefresh={handleRefresh} />
          </div>
        )}

        {/* Contextual Permissions Tab */}
        {activeTab === 'contextual' && (
          <div className="space-y-6">
            <ContextualPermissionManager />
          </div>
        )}

        {/* User Assignments Tab */}
        {activeTab === 'assignments' && (
          <div className="space-y-6">
            <UserRoleAssignment roles={roles} onRefresh={handleRefresh} />
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

// Role Templates View Component (Enhanced)
const RoleTemplatesView: React.FC = () => {
  const [selectedCategory, setSelectedCategory] = useState<string>('all');

  const categories = [
    { id: 'all', name: 'All Templates', count: 6 },
    { id: 'system', name: 'System', count: 2 },
    { id: 'business', name: 'Business', count: 2 },
    { id: 'basic', name: 'Basic', count: 2 },
  ];

  const templates = [
    {
      id: 1,
      name: 'Global Administrator',
      category: 'system',
      description: 'Complete system administration access across all organizations',
      permissions: ['*:*'],
      context: 'global',
      users: 2,
      icon: '🔒',
      color: 'red',
    },
    {
      id: 2,
      name: 'Organization Administrator',
      category: 'system',
      description: 'Full administrative access within a specific organization',
      permissions: ['org:*', 'users:*', 'roles:*'],
      context: 'organization',
      users: 5,
      icon: '🏢',
      color: 'blue',
    },
    {
      id: 3,
      name: 'Team Lead',
      category: 'business',
      description: 'Team management and project oversight capabilities',
      permissions: ['team:manage', 'users:read', 'projects:*'],
      context: 'team',
      users: 12,
      icon: '👨‍💼',
      color: 'green',
    },
    {
      id: 4,
      name: 'Project Manager',
      category: 'business',
      description: 'Project coordination and resource management',
      permissions: ['projects:*', 'reports:read', 'resources:manage'],
      context: 'project',
      users: 8,
      icon: '📊',
      color: 'purple',
    },
    {
      id: 5,
      name: 'Content Editor',
      category: 'basic',
      description: 'Create, edit, and publish content within assigned areas',
      permissions: ['content:read', 'content:create', 'content:update'],
      context: 'resource',
      users: 25,
      icon: '✏️',
      color: 'indigo',
    },
    {
      id: 6,
      name: 'Viewer',
      category: 'basic',
      description: 'Read-only access to assigned resources and data',
      permissions: ['*:read'],
      context: 'resource',
      users: 45,
      icon: '👁️',
      color: 'gray',
    },
  ];

  const filteredTemplates = selectedCategory === 'all' 
    ? templates 
    : templates.filter(template => template.category === selectedCategory);

  const getColorClasses = (color: string) => {
    const colors = {
      red: 'bg-red-100 text-red-800 border-red-200',
      blue: 'bg-blue-100 text-blue-800 border-blue-200',
      green: 'bg-green-100 text-green-800 border-green-200',
      purple: 'bg-purple-100 text-purple-800 border-purple-200',
      indigo: 'bg-indigo-100 text-indigo-800 border-indigo-200',
      gray: 'bg-gray-100 text-gray-800 border-gray-200',
    };
    return colors[color] || colors.gray;
  };

  return (
    <div className="space-y-6">
      {/* Category Filter */}
      <div className="bg-white shadow rounded-lg">
        <div className="px-4 py-5 sm:px-6 border-b border-gray-200">
          <h3 className="text-lg font-medium text-gray-900">Role Templates</h3>
          <p className="mt-1 text-sm text-gray-500">
            Pre-configured role templates for quick role creation and standardization
          </p>
        </div>
        
        <div className="px-4 py-4">
          <nav className="flex space-x-8" aria-label="Tabs">
            {categories.map((category) => (
              <button
                key={category.id}
                onClick={() => setSelectedCategory(category.id)}
                className={`py-2 px-1 border-b-2 font-medium text-sm whitespace-nowrap ${
                  selectedCategory === category.id
                    ? 'border-indigo-500 text-indigo-600'
                    : 'border-transparent text-gray-500 hover:text-gray-700 hover:border-gray-300'
                }`}
              >
                {category.name}
                <span className="ml-2 py-0.5 px-2 text-xs bg-gray-100 text-gray-900 rounded-full">
                  {category.count}
                </span>
              </button>
            ))}
          </nav>
        </div>
      </div>

      {/* Templates Grid */}
      <div className="grid grid-cols-1 gap-6 sm:grid-cols-2 lg:grid-cols-3">
        {filteredTemplates.map((template) => (
          <div
            key={template.id}
            className={`border-2 rounded-lg p-6 hover:shadow-md transition-shadow ${getColorClasses(template.color)}`}
          >
            <div className="flex items-start justify-between">
              <div className="flex items-center space-x-3">
                <span className="text-2xl">{template.icon}</span>
                <div>
                  <h4 className="text-lg font-medium">{template.name}</h4>
                  <span className={`inline-flex px-2 py-1 text-xs font-semibold rounded-full bg-white`}>
                    {template.context}
                  </span>
                </div>
              </div>
              <div className="text-right">
                <div className="text-sm font-medium">{template.users}</div>
                <div className="text-xs text-gray-600">users</div>
              </div>
            </div>

            <p className="mt-4 text-sm">{template.description}</p>

            <div className="mt-4">
              <div className="text-xs font-medium text-gray-700 mb-2">Permissions:</div>
              <div className="flex flex-wrap gap-1">
                {template.permissions.slice(0, 3).map((permission, index) => (
                  <span
                    key={index}
                    className="inline-flex px-2 py-1 text-xs bg-white rounded border"
                  >
                    {permission}
                  </span>
                ))}
                {template.permissions.length > 3 && (
                  <span className="inline-flex px-2 py-1 text-xs bg-white rounded border">
                    +{template.permissions.length - 3} more
                  </span>
                )}
              </div>
            </div>

            <div className="mt-6 flex space-x-3">
              <button className="flex-1 bg-white text-gray-700 border border-gray-300 rounded-md px-3 py-2 text-sm font-medium hover:bg-gray-50">
                Preview
              </button>
              <button className="flex-1 bg-indigo-600 text-white rounded-md px-3 py-2 text-sm font-medium hover:bg-indigo-700">
                Use Template
              </button>
            </div>
          </div>
        ))}
      </div>

      {/* Usage Guide */}
      <div className="bg-white shadow rounded-lg">
        <div className="px-4 py-5 sm:p-6">
          <h3 className="text-lg font-medium text-gray-900 mb-4">How to Use Role Templates</h3>
          <div className="grid grid-cols-1 md:grid-cols-2 gap-6">
            <div>
              <h4 className="font-medium text-gray-900 mb-2">Quick Creation</h4>
              <ol className="text-sm text-gray-600 space-y-2">
                <li>1. Browse available templates by category</li>
                <li>2. Preview template permissions and context</li>
                <li>3. Click "Use Template" to create a new role</li>
                <li>4. Customize the role name and description</li>
                <li>5. Assign the role to users as needed</li>
              </ol>
            </div>
            <div>
              <h4 className="font-medium text-gray-900 mb-2">Best Practices</h4>
              <ul className="text-sm text-gray-600 space-y-2">
                <li>• Start with basic templates and customize as needed</li>
                <li>• Use system templates for administrative roles</li>
                <li>• Create custom templates for recurring role patterns</li>
                <li>• Review and audit role permissions regularly</li>
                <li>• Document custom role purposes and contexts</li>
              </ul>
            </div>
          </div>
        </div>
      </div>
    </div>
  );
};

export default UnifiedRBAC;