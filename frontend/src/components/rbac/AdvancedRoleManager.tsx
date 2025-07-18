import { normalize } from 'node:path/posix';
import React, { useState, useEffect } from 'react';
import { useOrganization } from '../../context/OrganizationContext';
import { advancedRbacService, Role, RoleTemplate, CreateRoleFromTemplateRequest } from '../../services/advancedRbacService';
import Table from '../common/Table';

const AdvancedRoleManager: React.FC = () => {
  const { currentOrganization, getCurrentOrgId } = useOrganization();
  const [roles, setRoles] = useState<Role[]>([]);
  const [roleTemplates, setRoleTemplates] = useState<RoleTemplate[]>([]);
  const [selectedRole, setSelectedRole] = useState<Role | null>(null);
  const [showCreateModal, setShowCreateModal] = useState(false);
  const [showCreateBasicModal, setShowCreateBasicModal] = useState(false);
  const [showHierarchyModal, setShowHierarchyModal] = useState(false);
  const [isLoading, setIsLoading] = useState(false);

  useEffect(() => {
    loadRoles();
    loadRoleTemplates();
  }, [currentOrganization]);

  const loadRoles = async () => {
    setIsLoading(true);
    try {
      console.log('Loading roles...');
      // Use the existing /rbac/roles endpoint
      const response = await fetch('/api/v1/rbac/roles', {
        method: 'GET',
        headers: {
          'Authorization': `Bearer ${localStorage.getItem('token')}`,
          'X-Organization-ID': getCurrentOrgId()?.toString() || '',
        },
      });

      if (!response.ok) {
        throw new Error(`HTTP ${response.status}`);
      }

      const result = await response.json();
      console.log('Raw API response:', result);
      const rolesData = Array.isArray(result.data) ? result.data : Array.isArray(result) ? result : [];
      console.log('Processed roles data:', rolesData);
      console.log('First role (if any):', rolesData[0]);

      // Simple sort by display_name for now (no hierarchy until backend supports it)
      const sortedRoles = rolesData.sort((a, b) => a.display_name.localeCompare(b.display_name));
      setRoles(sortedRoles);
    } catch (error) {
      console.error('Failed to load roles:', error);
      setRoles([]);
    } finally {
      setIsLoading(false);
    }
  };

  const loadRoleTemplates = async () => {
    try {
      console.log('Loading role templates...');
      const response = await advancedRbacService.getRoleTemplates();
      const templates = Array.isArray(response.data) ? response.data : [];
      console.log('Role templates loaded:', templates);
      setRoleTemplates(templates);
    } catch (error) {
      console.error('Failed to load role templates:', error);
      // Provide fallback mock templates if backend isn't available
      setRoleTemplates([
        {
          id: 1,
          name: 'org_admin',
          display_name: 'Organization Administrator',
          description: 'Full administrative access within organization',
          category: 'system',
          config: '{"permissions":["*:*"],"hierarchy_level":1}',
          is_active: true,
          created_at: new Date().toISOString(),
          updated_at: new Date().toISOString(),
        },
        {
          id: 2,
          name: 'team_lead',
          display_name: 'Team Lead',
          description: 'Team management and project oversight',
          category: 'business',
          config: '{"permissions":["users:read","users:update","projects:*"],"hierarchy_level":2}',
          is_active: true,
          created_at: new Date().toISOString(),
          updated_at: new Date().toISOString(),
        },
        {
          id: 3,
          name: 'viewer',
          display_name: 'Viewer',
          description: 'Read-only access to assigned resources',
          category: 'basic',
          config: '{"permissions":["*:read"],"hierarchy_level":3}',
          is_active: true,
          created_at: new Date().toISOString(),
          updated_at: new Date().toISOString(),
        }
      ]);
    }
  };

  const handleCreateFromTemplate = async (data: CreateRoleFromTemplateRequest) => {
    try {
      console.log('Creating role from template:', data);
      const response = await advancedRbacService.createRoleFromTemplate(data);
      console.log('Role created successfully:', response.data);
      setShowCreateModal(false);
      await loadRoles();
    } catch (error: any) {
      console.error('Failed to create role from template:', error);

      // If template creation fails, try to create a basic role instead
      const template = roleTemplates.find(t => t.id === data.template_id);
      if (template) {
        console.log('Falling back to basic role creation...');
        try {
          const basicRoleData = {
            name: data.custom_name || template.name,
            display_name: data.custom_name || template.display_name,
            description: template.description,
          };
          await handleCreateBasicRole(basicRoleData);
        } catch (fallbackError: any) {
          console.error('Fallback role creation also failed:', fallbackError);
          alert(`Failed to create role: ${error.response?.data?.message || error.message || 'Template and fallback creation both failed'}`);
        }
      } else {
        alert(`Failed to create role: ${error.response?.data?.message || error.message || 'Unknown error'}`);
      }
    }
  };

  const handleCreateBasicRole = async (data: { name: string; display_name: string; description: string }) => {
    try {
      console.log('Creating basic role:', data);
      // Use the existing /rbac/roles endpoint
      const response = await fetch('/api/v1/rbac/roles', {
        method: 'POST',
        headers: {
          'Content-Type': 'application/json',
          'Authorization': `Bearer ${localStorage.getItem('token')}`,
          'X-Organization-ID': getCurrentOrgId()?.toString() || '',
        },
        body: JSON.stringify(data),
      });

      if (!response.ok) {
        const errorData = await response.json();
        throw new Error(errorData.message || `HTTP ${response.status}`);
      }

      const result = await response.json();
      console.log('Basic role created successfully:', result);
      setShowCreateBasicModal(false);
      await loadRoles();
    } catch (error: any) {
      console.error('Failed to create basic role:', error);
      alert(`Failed to create role: ${error.message || 'Unknown error'}`);
    }
  };

  const handleSetParent = async (roleId: number, parentRoleId?: number) => {
    try {
      await advancedRbacService.setRoleParent(roleId, { parent_role_id: parentRoleId });
      await loadRoles();
    } catch (error) {
      console.error('Failed to set role parent:', error);
    }
  };

  const roleColumns = [
    {
      key: 'display_name',
      header: 'Role',
      render: (_, role: Role & { indentLevel?: number }) => (
        <div className="flex items-center">
          <div style={{ paddingLeft: `${(role.indentLevel || 0) * 20}px` }}>
            <div className="flex items-center space-x-2">
              {/* Show hierarchy indicators */}
              {role.indentLevel > 0 && (
                <span className="text-gray-400">
                  {'└' + '─'.repeat(Math.min(role.indentLevel, 3))}
                </span>
              )}
              <span className="text-lg">{advancedRbacService.getRoleIcon(role)}</span>
              <div>
                <div className="font-medium text-gray-900">{role?.display_name || role?.name || 'Unknown Role'}</div>
                <div className="text-sm text-gray-500">{role?.name || 'unknown'}</div>
              </div>
            </div>
          </div>
        </div>
      ),
    },
    {
      key: 'hierarchy_level',
      header: 'Level',
      render: (_, role: Role) => (
        <span className={`inline-flex px-2 py-1 text-xs font-semibold rounded-full ${advancedRbacService.getRoleColor(role)}`}>
          Level {role?.hierarchy_level ?? 0}
        </span>
      ),
    },
    {
      key: 'type',
      header: 'Type',
      render: (_, role: Role) => (
        <span className={`inline-flex px-2 py-1 text-xs font-semibold rounded-full ${role?.is_system ? 'bg-red-100 text-red-800' : 'bg-blue-100 text-blue-800'}`}>
          {role?.is_system ? 'System' : 'Custom'}
        </span>
      ),
    },
    {
      key: 'parent_role',
      header: 'Parent Role',
      render: (_, role: Role) => (
        role?.parent_role ? (
          <span className="text-sm text-gray-600">{role.parent_role?.display_name || role.parent_role?.name || 'Unknown'}</span>
        ) : (
          <span className="text-sm text-gray-400">None</span>
        )
      ),
    },
    {
      key: 'actions',
      header: 'Actions',
      render: (_, role: Role) => (
        <div className="flex space-x-2">
          <button
            onClick={() => {
              setSelectedRole(role);
              setShowHierarchyModal(true);
            }}
            className="text-indigo-600 hover:text-indigo-900 text-sm"
            disabled={role?.is_system}
            title={role?.is_system ? "Cannot edit system roles" : "Edit Hierarchy"}
          >
            Edit Hierarchy
          </button>
          <button
            onClick={() => {
              setSelectedRole(role);
              setShowHierarchyModal(false); // <-- ensure details panel shows
            }}
            className="text-blue-600 hover:text-blue-900 text-sm"
          >
            View Details
          </button>
        </div>
      ),
    },
  ];

  // Build proper role hierarchy display
  const displayRoles = React.useMemo(() => {
    if (!Array.isArray(roles) || roles.length === 0) {
      return [];
    }

    // First, normalize the roles
    console.log(roles);
    const normalizedRoles = roles.map(role => ({
      ...role,
      hierarchy_level: role?.hierarchy_level ?? 0,
      display_name: role?.display_name || role?.name || 'Unknown Role',
      name: role?.name || 'unknown',
      is_system: role?.is_system ?? false,
      parent_role_id: role?.parent_role_id || null,
    }));

    // Try to build tree structure if parent_role_id is available
    const hasHierarchy = normalizedRoles.some(role => role.parent_role_id);

    if (hasHierarchy) {
      // Use the advancedRbacService utility to build and flatten the tree
      const roleTree = advancedRbacService.buildRoleTree(normalizedRoles);
      return advancedRbacService.flattenRoleTree(roleTree);
    } else {
      // No hierarchy, just sort by hierarchy level and name
      return normalizedRoles
        .sort((a, b) => {
          if (a.hierarchy_level !== b.hierarchy_level) {
            return a.hierarchy_level - b.hierarchy_level;
          }
          return a.display_name.localeCompare(b.display_name);
        })
        .map(role => ({ ...role, indentLevel: 0 }));
    }
  }, [roles]);

  console.log(displayRoles)
  return (
    <div className="space-y-6">
      {/* Header */}
      <div className="bg-white shadow rounded-lg">
        <div className="px-4 py-5 sm:px-6 border-b border-gray-200">
          <div className="flex items-center justify-between">
            <div>
              <h3 className="text-lg font-medium text-gray-900">Advanced Role Management</h3>
              <p className="mt-1 text-sm text-gray-500">
                Manage hierarchical roles and permissions for {currentOrganization?.name || 'Global'}
              </p>
            </div>
            <div className="flex space-x-2">
              <button
                onClick={() => setShowCreateModal(true)}
                disabled={roleTemplates.length === 0}
                className="inline-flex items-center px-4 py-2 border border-transparent text-sm font-medium rounded-md shadow-sm text-white bg-indigo-600 hover:bg-indigo-700 disabled:opacity-50 disabled:cursor-not-allowed"
                title={roleTemplates.length === 0 ? "No role templates available" : "Create Role from Template"}
              >
                Create from Template
              </button>
              <button
                onClick={() => setShowCreateBasicModal(true)}
                className="inline-flex items-center px-4 py-2 border border-gray-300 text-sm font-medium rounded-md shadow-sm text-gray-700 bg-white hover:bg-gray-50"
              >
                Create Basic Role
              </button>
            </div>
          </div>
        </div>

        {/* Role Statistics */}
        <div className="px-4 py-5 sm:p-6">
          <div className="grid grid-cols-1 gap-5 sm:grid-cols-3">
            <div className="bg-gray-50 px-4 py-5 rounded-lg">
              <dt className="text-sm font-medium text-gray-500 truncate">Total Roles</dt>
              <dd className="mt-1 text-3xl font-semibold text-gray-900">{roles.length}</dd>
            </div>
            <div className="bg-gray-50 px-4 py-5 rounded-lg">
              <dt className="text-sm font-medium text-gray-500 truncate">Organization Roles</dt>
              <dd className="mt-1 text-3xl font-semibold text-gray-900">
                {roles.filter(role => role.org_id).length}
              </dd>
            </div>
            <div className="bg-gray-50 px-4 py-5 rounded-lg">
              <dt className="text-sm font-medium text-gray-500 truncate">Available Templates</dt>
              <dd className="mt-1 text-3xl font-semibold text-gray-900">{roleTemplates.length}</dd>
            </div>
          </div>
        </div>
      </div>

      {/* Roles Table */}
      <div className="bg-white shadow rounded-lg">
        <div className="px-4 py-5 sm:px-6 border-b border-gray-200">
          <h3 className="text-lg font-medium text-gray-900">Role Hierarchy</h3>
        </div>
        <Table
          data={displayRoles}
          columns={roleColumns}
          loading={isLoading}
          emptyMessage="No roles found"
        />
      </div>

      {/* Create Role Modal */}
      {showCreateModal && (
        <CreateRoleFromTemplateModal
          templates={roleTemplates}
          onSubmit={handleCreateFromTemplate}
          onClose={() => setShowCreateModal(false)}
        />
      )}

      {/* Create Basic Role Modal */}
      {showCreateBasicModal && (
        <CreateBasicRoleModal
          onSubmit={handleCreateBasicRole}
          onClose={() => setShowCreateBasicModal(false)}
        />
      )}

      {/* Hierarchy Management Modal */}
      {showHierarchyModal && selectedRole && (
        <RoleHierarchyModal
          role={selectedRole}
          availableParents={advancedRbacService.getAvailableParentRoles(roles, selectedRole.id)}
          onSubmit={handleSetParent}
          onClose={() => {
            setShowHierarchyModal(false);
            setSelectedRole(null);
          }}
        />
      )}

      {/* Role Details Panel */}
      {selectedRole && !showHierarchyModal && (
        <RoleDetailsPanel
          role={selectedRole}
          onClose={() => setSelectedRole(null)}
        />
      )}
    </div>
  );
};

// Create Role from Template Modal
interface CreateRoleFromTemplateModalProps {
  templates: RoleTemplate[];
  onSubmit: (data: CreateRoleFromTemplateRequest) => void;
  onClose: () => void;
}

const CreateRoleFromTemplateModal: React.FC<CreateRoleFromTemplateModalProps> = ({
  templates,
  onSubmit,
  onClose,
}) => {
  const [formData, setFormData] = useState<CreateRoleFromTemplateRequest>({
    template_id: 0,
    custom_name: '',
  });

  const handleSubmit = (e: React.FormEvent) => {
    e.preventDefault();

    if (formData.template_id === 0) {
      alert('Please select a template');
      return;
    }

    onSubmit(formData);
  };

  const selectedTemplate = templates.find(t => t.id === formData.template_id);

  return (
    <div className="fixed inset-0 z-50 overflow-y-auto">
      <div className="flex items-center justify-center min-h-screen pt-4 px-4 pb-20 text-center sm:block sm:p-0">
        <div className="fixed inset-0 bg-gray-500 bg-opacity-75 transition-opacity" onClick={onClose}></div>

        <div className="inline-block align-bottom bg-white rounded-lg text-left overflow-hidden shadow-xl transform transition-all sm:my-8 sm:align-middle sm:max-w-lg sm:w-full">
          <form onSubmit={handleSubmit}>
            <div className="bg-white px-4 pt-5 pb-4 sm:p-6 sm:pb-4">
              <h3 className="text-lg font-medium text-gray-900 mb-4">Create Role from Template</h3>

              <div className="space-y-4">
                <div>
                  <label className="block text-sm font-medium text-gray-700">Role Template</label>
                  <select
                    value={formData.template_id}
                    onChange={(e) => setFormData({ ...formData, template_id: parseInt(e.target.value) })}
                    className="mt-1 block w-full border-gray-300 rounded-md shadow-sm focus:ring-indigo-500 focus:border-indigo-500"
                    required
                  >
                    <option value={0}>Select a template</option>
                    {templates.map((template) => (
                      <option key={template.id} value={template.id}>
                        {template.display_name} ({template.category})
                      </option>
                    ))}
                  </select>
                  {selectedTemplate && (
                    <p className="mt-1 text-sm text-gray-500">{selectedTemplate.description}</p>
                  )}
                </div>

                <div>
                  <label className="block text-sm font-medium text-gray-700">Custom Name (Optional)</label>
                  <input
                    type="text"
                    value={formData.custom_name}
                    onChange={(e) => setFormData({ ...formData, custom_name: e.target.value })}
                    className="mt-1 block w-full border-gray-300 rounded-md shadow-sm focus:ring-indigo-500 focus:border-indigo-500"
                    placeholder="Leave empty to use template name"
                  />
                </div>
              </div>
            </div>

            <div className="bg-gray-50 px-4 py-3 sm:px-6 sm:flex sm:flex-row-reverse">
              <button
                type="submit"
                disabled={formData.template_id === 0}
                className="w-full inline-flex justify-center rounded-md border border-transparent shadow-sm px-4 py-2 bg-indigo-600 text-base font-medium text-white hover:bg-indigo-700 focus:outline-none focus:ring-2 focus:ring-offset-2 focus:ring-indigo-500 sm:ml-3 sm:w-auto sm:text-sm disabled:opacity-50"
              >
                Create Role
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

// Role Hierarchy Modal
interface RoleHierarchyModalProps {
  role: Role;
  availableParents: Role[];
  onSubmit: (roleId: number, parentRoleId?: number) => void;
  onClose: () => void;
}

const RoleHierarchyModal: React.FC<RoleHierarchyModalProps> = ({
  role,
  availableParents,
  onSubmit,
  onClose,
}) => {
  const [parentRoleId, setParentRoleId] = useState<number | undefined>(role.parent_role_id);

  const handleSubmit = (e: React.FormEvent) => {
    e.preventDefault();
    onSubmit(role.id, parentRoleId);
    onClose();
  };

  return (
    <div className="fixed inset-0 z-50 overflow-y-auto">
      <div className="flex items-center justify-center min-h-screen pt-4 px-4 pb-20 text-center sm:block sm:p-0">
        <div className="fixed inset-0 bg-gray-500 bg-opacity-75 transition-opacity" onClick={onClose}></div>

        <div className="inline-block align-bottom bg-white rounded-lg text-left overflow-hidden shadow-xl transform transition-all sm:my-8 sm:align-middle sm:max-w-lg sm:w-full">
          <form onSubmit={handleSubmit}>
            <div className="bg-white px-4 pt-5 pb-4 sm:p-6 sm:pb-4">
              <h3 className="text-lg font-medium text-gray-900 mb-4">
                Edit Role Hierarchy: {role.display_name}
              </h3>

              <div className="space-y-4">
                <div>
                  <label className="block text-sm font-medium text-gray-700">Parent Role</label>
                  <select
                    value={parentRoleId || ''}
                    onChange={(e) => setParentRoleId(e.target.value ? parseInt(e.target.value) : undefined)}
                    className="mt-1 block w-full border-gray-300 rounded-md shadow-sm focus:ring-indigo-500 focus:border-indigo-500"
                  >
                    <option value="">No Parent (Root Level)</option>
                    {availableParents.map((parent) => (
                      <option key={parent.id} value={parent.id}>
                        {parent.display_name} (Level {parent.hierarchy_level})
                      </option>
                    ))}
                  </select>
                  <p className="mt-1 text-sm text-gray-500">
                    Selecting a parent will inherit its permissions and place this role at level {parentRoleId ? (availableParents.find(p => p.id === parentRoleId)?.hierarchy_level || 0) + 1 : 0}
                  </p>
                </div>
              </div>
            </div>

            <div className="bg-gray-50 px-4 py-3 sm:px-6 sm:flex sm:flex-row-reverse">
              <button
                type="submit"
                className="w-full inline-flex justify-center rounded-md border border-transparent shadow-sm px-4 py-2 bg-indigo-600 text-base font-medium text-white hover:bg-indigo-700 focus:outline-none focus:ring-2 focus:ring-offset-2 focus:ring-indigo-500 sm:ml-3 sm:w-auto sm:text-sm"
              >
                Update Hierarchy
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

// Role Details Panel
interface RoleDetailsPanelProps {
  role: Role;
  onClose: () => void;
}

const RoleDetailsPanel: React.FC<RoleDetailsPanelProps> = ({ role, onClose }) => {
  return (
    <div className="bg-white shadow rounded-lg">
      <div className="px-4 py-5 sm:px-6 border-b border-gray-200">
        <div className="flex items-center justify-between">
          <div className="flex items-center space-x-3">
            <span className="text-2xl">{role?.is_system ? '🔒' : '👤'}</span>
            <div>
              <h3 className="text-lg font-medium text-gray-900">{role?.display_name || role?.name || 'Unknown Role'}</h3>
              <p className="text-sm text-gray-500">{role?.name || 'unknown'}</p>
            </div>
          </div>
          <button
            onClick={onClose}
            className="text-gray-400 hover:text-gray-600"
          >
            <svg className="w-6 h-6" fill="none" stroke="currentColor" viewBox="0 0 24 24">
              <path strokeLinecap="round" strokeLinejoin="round" strokeWidth={2} d="M6 18L18 6M6 6l12 12" />
            </svg>
          </button>
        </div>
      </div>

      <div className="px-4 py-5 sm:p-6">
        <dl className="grid grid-cols-1 gap-x-4 gap-y-6 sm:grid-cols-2">
          <div>
            <dt className="text-sm font-medium text-gray-500">Description</dt>
            <dd className="mt-1 text-sm text-gray-900">{role?.description || 'No description'}</dd>
          </div>
          <div>
            <dt className="text-sm font-medium text-gray-500">Hierarchy Level</dt>
            <dd className="mt-1 text-sm text-gray-900">Level {role?.hierarchy_level ?? 0}</dd>
          </div>
          <div>
            <dt className="text-sm font-medium text-gray-500">Type</dt>
            <dd className="mt-1">
              <span className={`inline-flex px-2 py-1 text-xs font-semibold rounded-full ${role?.is_system ? 'bg-red-100 text-red-800' : 'bg-blue-100 text-blue-800'}`}>
                {role?.is_system ? 'System Role' : 'Custom Role'}
              </span>
            </dd>
          </div>
          <div>
            <dt className="text-sm font-medium text-gray-500">Status</dt>
            <dd className="mt-1">
              <span className={`inline-flex px-2 py-1 text-xs font-semibold rounded-full ${role?.is_active !== false ? 'bg-green-100 text-green-800' : 'bg-red-100 text-red-800'}`}>
                {role?.is_active !== false ? 'Active' : 'Inactive'}
              </span>
            </dd>
          </div>
          {role?.parent_role && (
            <div>
              <dt className="text-sm font-medium text-gray-500">Parent Role</dt>
              <dd className="mt-1 text-sm text-gray-900">{role.parent_role?.display_name || role.parent_role?.name || 'Unknown'}</dd>
            </div>
          )}
          {role?.organization && (
            <div>
              <dt className="text-sm font-medium text-gray-500">Organization</dt>
              <dd className="mt-1 text-sm text-gray-900">{role.organization?.name || 'Unknown'}</dd>
            </div>
          )}
        </dl>
      </div>
    </div>
  );
};

// Create Basic Role Modal
interface CreateBasicRoleModalProps {
  onSubmit: (data: { name: string; display_name: string; description: string }) => void;
  onClose: () => void;
}

const CreateBasicRoleModal: React.FC<CreateBasicRoleModalProps> = ({
  onSubmit,
  onClose,
}) => {
  const [formData, setFormData] = useState({
    name: '',
    display_name: '',
    description: '',
  });

  const handleSubmit = (e: React.FormEvent) => {
    e.preventDefault();

    if (!formData.name.trim() || !formData.display_name.trim()) {
      alert('Please fill in all required fields');
      return;
    }

    onSubmit(formData);
  };

  return (
    <div className="fixed inset-0 z-50 overflow-y-auto">
      <div className="flex items-center justify-center min-h-screen pt-4 px-4 pb-20 text-center sm:block sm:p-0">
        <div className="fixed inset-0 bg-gray-500 bg-opacity-75 transition-opacity" onClick={onClose}></div>

        <div className="inline-block align-bottom bg-white rounded-lg text-left overflow-hidden shadow-xl transform transition-all sm:my-8 sm:align-middle sm:max-w-lg sm:w-full">
          <form onSubmit={handleSubmit}>
            <div className="bg-white px-4 pt-5 pb-4 sm:p-6 sm:pb-4">
              <h3 className="text-lg font-medium text-gray-900 mb-4">Create Basic Role</h3>

              <div className="space-y-4">
                <div>
                  <label className="block text-sm font-medium text-gray-700">Role Name *</label>
                  <input
                    type="text"
                    value={formData.name}
                    onChange={(e) => setFormData({ ...formData, name: e.target.value })}
                    className="mt-1 block w-full border-gray-300 rounded-md shadow-sm focus:ring-indigo-500 focus:border-indigo-500"
                    placeholder="e.g., editor"
                    required
                  />
                  <p className="mt-1 text-xs text-gray-500">
                    Unique identifier for the role (lowercase, no spaces)
                  </p>
                </div>

                <div>
                  <label className="block text-sm font-medium text-gray-700">Display Name *</label>
                  <input
                    type="text"
                    value={formData.display_name}
                    onChange={(e) => setFormData({ ...formData, display_name: e.target.value })}
                    className="mt-1 block w-full border-gray-300 rounded-md shadow-sm focus:ring-indigo-500 focus:border-indigo-500"
                    placeholder="e.g., Content Editor"
                    required
                  />
                  <p className="mt-1 text-xs text-gray-500">
                    Human-readable name for the role
                  </p>
                </div>

                <div>
                  <label className="block text-sm font-medium text-gray-700">Description</label>
                  <textarea
                    value={formData.description}
                    onChange={(e) => setFormData({ ...formData, description: e.target.value })}
                    className="mt-1 block w-full border-gray-300 rounded-md shadow-sm focus:ring-indigo-500 focus:border-indigo-500"
                    rows={3}
                    placeholder="Brief description of what this role can do..."
                  />
                </div>
              </div>
            </div>

            <div className="bg-gray-50 px-4 py-3 sm:px-6 sm:flex sm:flex-row-reverse">
              <button
                type="submit"
                disabled={!formData.name.trim() || !formData.display_name.trim()}
                className="w-full inline-flex justify-center rounded-md border border-transparent shadow-sm px-4 py-2 bg-indigo-600 text-base font-medium text-white hover:bg-indigo-700 focus:outline-none focus:ring-2 focus:ring-offset-2 focus:ring-indigo-500 sm:ml-3 sm:w-auto sm:text-sm disabled:opacity-50"
              >
                Create Role
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

export default AdvancedRoleManager;
