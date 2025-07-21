import React, { useState, useEffect } from 'react';
import { rbacService, Role } from '../../services/rbacService';
import Table from '../common/Table';
import Icon from '../common/Icons';
import RoleHierarchyView from './RoleHierarchyView';

const SimplifiedRoleManager: React.FC = () => {
  const [roles, setRoles] = useState<Role[]>([]);
  const [showCreateModal, setShowCreateModal] = useState(false);
  const [showHierarchyView, setShowHierarchyView] = useState(false);
  const [selectedRole, setSelectedRole] = useState<Role | null>(null);
  const [showSetParentModal, setShowSetParentModal] = useState(false);
  const [isLoading, setIsLoading] = useState(false);

  useEffect(() => {
    loadRoles();
  }, []);

  const loadRoles = async () => {
    setIsLoading(true);
    try {
      const response = await rbacService.getRoles();
      const rolesData = Array.isArray(response.data) ? response.data : [];
      setRoles(rolesData);
    } catch (error) {
      console.error('Failed to load roles:', error);
      setRoles([]);
    } finally {
      setIsLoading(false);
    }
  };

  const handleCreateRole = async (data: { name: string; display_name: string; description: string }) => {
    try {
      await rbacService.createRole(data);
      setShowCreateModal(false);
      await loadRoles();
    } catch (error: any) {
      console.error('Failed to create basic role:', error);
      alert(`Failed to create role: ${error.message || 'Unknown error'}`);
    }
  };

  const handleDeleteRole = async (role: Role) => {
    if (!window.confirm(`Are you sure you want to delete the role "${role.display_name}"?`)) {
      return;
    }

    try {
      await rbacService.deleteRole(role.id);
      await loadRoles();
    } catch (error: any) {
      console.error('Failed to delete role:', error);
      alert(`Failed to delete role: ${error.message || 'Unknown error'}`);
    }
  };

  const handleSetParent = (role: Role) => {
    setSelectedRole(role);
    setShowSetParentModal(true);
  };

  const handleSetParentSubmit = async (parentRoleId: number | null) => {
    if (!selectedRole) return;

    try {
      await rbacService.setRoleParent(selectedRole.id, parentRoleId);
      setShowSetParentModal(false);
      setSelectedRole(null);
      await loadRoles();
    } catch (error: any) {
      console.error('Failed to set parent role:', error);
      alert(`Failed to set parent role: ${error.message || 'Unknown error'}`);
    }
  };

  const columns = [
    {
      key: 'display_name',
      header: 'Role Name',
      sortable: true,
      render: (value: string, role: Role) => (
        <div>
          <div className="flex items-center">
            <span className="mr-2">
              {role.hierarchy_level === 0 ? '👑' : role.hierarchy_level === 1 ? '🎭' : '👤'}
            </span>
            <div>
              <div className="font-medium text-gray-900 dark:text-white">
                {value || role.name}
              </div>
              <div className="text-sm text-gray-500">
                {role.name} {role.parent_role_id && (
                  <span className="inline-flex items-center px-1.5 py-0.5 rounded-full text-xs font-medium bg-blue-100 text-blue-800 ml-1">
                    ↗ Inherits
                  </span>
                )}
              </div>
            </div>
          </div>
        </div>
      )
    },
    {
      key: 'hierarchy_level',
      header: 'Hierarchy',
      sortable: true,
      render: (value: number, role: Role) => (
        <div>
          <span className={`inline-flex items-center px-2 py-1 rounded-full text-xs font-medium ${
            value === 0 ? 'bg-yellow-100 text-yellow-800' : 
            value === 1 ? 'bg-purple-100 text-purple-800' : 
            'bg-gray-100 text-gray-800'
          }`}>
            Level {value}
          </span>
          {role.parent_role && (
            <div className="text-xs text-gray-500 mt-1">
              Parent: {role.parent_role.display_name}
            </div>
          )}
        </div>
      )
    },
    {
      key: 'description',
      header: 'Description',
      render: (value: string) => (
        <span className="text-sm text-gray-600 dark:text-gray-300">
          {value || 'No description'}
        </span>
      )
    },
    {
      key: 'is_system',
      header: 'Type',
      render: (value: boolean) => (
        <span className={`inline-flex items-center px-2 py-1 rounded-full text-xs font-medium ${
          value ? 'bg-red-100 text-red-800' : 'bg-blue-100 text-blue-800'
        }`}>
          {value ? 'System' : 'Custom'}
        </span>
      )
    },
    {
      key: 'is_active',
      header: 'Status',
      render: (value: boolean) => (
        <span className={`inline-flex items-center px-2 py-1 rounded-full text-xs font-medium ${
          value ? 'bg-green-100 text-green-800' : 'bg-gray-100 text-gray-800'
        }`}>
          {value ? 'Active' : 'Inactive'}
        </span>
      )
    },
    {
      key: 'actions',
      header: 'Actions',
      render: (value: any, role: Role) => (
        <div className="flex space-x-2">
          <button
            onClick={() => handleSetParent(role)}
            disabled={role.is_system}
            className="text-indigo-600 hover:text-indigo-900 disabled:opacity-50 disabled:cursor-not-allowed"
            title={role.is_system ? 'Cannot modify system roles' : 'Set parent role'}
          >
            <Icon name="link" className="w-4 h-4" />
          </button>
          <button
            onClick={() => handleDeleteRole(role)}
            disabled={role.is_system}
            className="text-red-600 hover:text-red-900 disabled:opacity-50 disabled:cursor-not-allowed"
            title={role.is_system ? 'Cannot delete system roles' : 'Delete role'}
          >
            <Icon name="delete" className="w-4 h-4" />
          </button>
        </div>
      )
    }
  ];

  return (
    <div className="space-y-6">
      {/* Header */}
      <div className="flex justify-between items-center">
        <div>
          <h2 className="text-xl font-bold text-gray-900 dark:text-white">Role Management</h2>
          <p className="text-sm text-gray-500 dark:text-gray-400">
            Manage system roles and permissions with hierarchy support
          </p>
        </div>
        <div className="flex space-x-2">
          <button
            onClick={() => setShowHierarchyView(!showHierarchyView)}
            className="bg-gray-100 hover:bg-gray-200 text-gray-700 px-4 py-2 rounded-md text-sm font-medium flex items-center"
          >
            <Icon name="tree" className="w-4 h-4 mr-2" />
            {showHierarchyView ? 'Hide' : 'Show'} Hierarchy
          </button>
          <button
            onClick={() => setShowCreateModal(true)}
            className="bg-indigo-600 hover:bg-indigo-700 text-white px-4 py-2 rounded-md text-sm font-medium flex items-center"
          >
            <Icon name="plus" className="w-4 h-4 mr-2" />
            Create Role
          </button>
        </div>
      </div>

      {/* Hierarchy View */}
      {showHierarchyView && (
        <RoleHierarchyView 
          roles={roles}
          onRoleSelect={(role) => setSelectedRole(role)}
          selectedRoleId={selectedRole?.id}
        />
      )}

      {/* Roles Table */}
      <div className="bg-white dark:bg-gray-800 shadow overflow-hidden rounded-lg">
        <Table
          columns={columns}
          data={roles}
          loading={isLoading}
          emptyMessage="No roles found"
        />
      </div>

      {/* Create Role Modal */}
      {showCreateModal && (
        <CreateRoleModal
          onSubmit={handleCreateRole}
          onClose={() => setShowCreateModal(false)}
        />
      )}

      {/* Set Parent Modal */}
      {showSetParentModal && selectedRole && (
        <SetParentModal
          role={selectedRole}
          availableRoles={roles.filter(r => r.id !== selectedRole.id && !r.is_system)}
          onSubmit={handleSetParentSubmit}
          onClose={() => {
            setShowSetParentModal(false);
            setSelectedRole(null);
          }}
        />
      )}
    </div>
  );
};

// Simple Create Role Modal
interface CreateRoleModalProps {
  onSubmit: (data: { name: string; display_name: string; description: string }) => void;
  onClose: () => void;
}

const CreateRoleModal: React.FC<CreateRoleModalProps> = ({ onSubmit, onClose }) => {
  const [formData, setFormData] = useState({
    name: '',
    display_name: '',
    description: '',
  });

  const handleSubmit = (e: React.FormEvent) => {
    e.preventDefault();
    if (!formData.name.trim() || !formData.display_name.trim()) {
      alert('Please fill in required fields');
      return;
    }
    onSubmit(formData);
  };

  return (
    <div className="fixed inset-0 z-50 overflow-y-auto">
      <div className="flex items-center justify-center min-h-screen pt-4 px-4 pb-20 text-center sm:block sm:p-0">
        <div className="fixed inset-0 bg-gray-500 bg-opacity-75 transition-opacity" onClick={onClose}></div>
        
        <div className="inline-block align-bottom bg-white dark:bg-gray-800 rounded-lg text-left overflow-hidden shadow-xl transform transition-all sm:my-8 sm:align-middle sm:max-w-lg sm:w-full">
          <form onSubmit={handleSubmit}>
            <div className="bg-white dark:bg-gray-800 px-4 pt-5 pb-4 sm:p-6 sm:pb-4">
              <h3 className="text-lg font-medium text-gray-900 dark:text-white mb-4">Create New Role</h3>
              
              <div className="space-y-4">
                <div>
                  <label className="block text-sm font-medium text-gray-700 dark:text-gray-300">
                    Role Name (Internal) *
                  </label>
                  <input
                    type="text"
                    value={formData.name}
                    onChange={(e) => setFormData({ ...formData, name: e.target.value })}
                    className="mt-1 block w-full border-gray-300 dark:border-gray-600 rounded-md shadow-sm focus:ring-indigo-500 focus:border-indigo-500 dark:bg-gray-700 dark:text-white"
                    placeholder="e.g., editor, manager"
                    required
                  />
                  <p className="mt-1 text-xs text-gray-500">Lowercase, no spaces, used internally</p>
                </div>

                <div>
                  <label className="block text-sm font-medium text-gray-700 dark:text-gray-300">
                    Display Name *
                  </label>
                  <input
                    type="text"
                    value={formData.display_name}
                    onChange={(e) => setFormData({ ...formData, display_name: e.target.value })}
                    className="mt-1 block w-full border-gray-300 dark:border-gray-600 rounded-md shadow-sm focus:ring-indigo-500 focus:border-indigo-500 dark:bg-gray-700 dark:text-white"
                    placeholder="e.g., Content Editor, Team Manager"
                    required
                  />
                  <p className="mt-1 text-xs text-gray-500">Human-readable name shown in UI</p>
                </div>

                <div>
                  <label className="block text-sm font-medium text-gray-700 dark:text-gray-300">
                    Description
                  </label>
                  <textarea
                    value={formData.description}
                    onChange={(e) => setFormData({ ...formData, description: e.target.value })}
                    rows={3}
                    className="mt-1 block w-full border-gray-300 dark:border-gray-600 rounded-md shadow-sm focus:ring-indigo-500 focus:border-indigo-500 dark:bg-gray-700 dark:text-white"
                    placeholder="Describe what this role is for..."
                  />
                </div>
              </div>
            </div>

            <div className="bg-gray-50 dark:bg-gray-700 px-4 py-3 sm:px-6 sm:flex sm:flex-row-reverse">
              <button
                type="submit"
                className="w-full inline-flex justify-center rounded-md border border-transparent shadow-sm px-4 py-2 bg-indigo-600 text-base font-medium text-white hover:bg-indigo-700 focus:outline-none focus:ring-2 focus:ring-offset-2 focus:ring-indigo-500 sm:ml-3 sm:w-auto sm:text-sm"
              >
                Create Role
              </button>
              <button
                type="button"
                onClick={onClose}
                className="mt-3 w-full inline-flex justify-center rounded-md border border-gray-300 dark:border-gray-600 shadow-sm px-4 py-2 bg-white dark:bg-gray-600 text-base font-medium text-gray-700 dark:text-gray-300 hover:bg-gray-50 dark:hover:bg-gray-500 focus:outline-none focus:ring-2 focus:ring-offset-2 focus:ring-indigo-500 sm:mt-0 sm:ml-3 sm:w-auto sm:text-sm"
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

// Set Parent Role Modal
interface SetParentModalProps {
  role: Role;
  availableRoles: Role[];
  onSubmit: (parentRoleId: number | null) => void;
  onClose: () => void;
}

const SetParentModal: React.FC<SetParentModalProps> = ({ role, availableRoles, onSubmit, onClose }) => {
  const [selectedParentId, setSelectedParentId] = useState<number | null>(role.parent_role_id || null);

  const handleSubmit = (e: React.FormEvent) => {
    e.preventDefault();
    onSubmit(selectedParentId);
  };

  return (
    <div className="fixed inset-0 z-50 overflow-y-auto">
      <div className="flex items-center justify-center min-h-screen pt-4 px-4 pb-20 text-center sm:block sm:p-0">
        <div className="fixed inset-0 bg-gray-500 bg-opacity-75 transition-opacity" onClick={onClose}></div>
        
        <div className="inline-block align-bottom bg-white dark:bg-gray-800 rounded-lg text-left overflow-hidden shadow-xl transform transition-all sm:my-8 sm:align-middle sm:max-w-lg sm:w-full">
          <form onSubmit={handleSubmit}>
            <div className="bg-white dark:bg-gray-800 px-4 pt-5 pb-4 sm:p-6 sm:pb-4">
              <h3 className="text-lg font-medium text-gray-900 dark:text-white mb-4">
                Set Parent Role for "{role.display_name}"
              </h3>
              
              <div className="space-y-4">
                <div>
                  <label className="block text-sm font-medium text-gray-700 dark:text-gray-300 mb-2">
                    Parent Role
                  </label>
                  
                  <div className="space-y-2">
                    <label className="flex items-center">
                      <input
                        type="radio"
                        value=""
                        checked={selectedParentId === null}
                        onChange={() => setSelectedParentId(null)}
                        className="mr-2"
                      />
                      <span className="text-sm text-gray-700 dark:text-gray-300">
                        No Parent (Root Role)
                      </span>
                    </label>
                    
                    {availableRoles.map(availableRole => (
                      <label key={availableRole.id} className="flex items-center">
                        <input
                          type="radio"
                          value={availableRole.id}
                          checked={selectedParentId === availableRole.id}
                          onChange={() => setSelectedParentId(availableRole.id)}
                          className="mr-2"
                        />
                        <div className="flex items-center">
                          <span className="mr-2">
                            {availableRole.hierarchy_level === 0 ? '👑' : availableRole.hierarchy_level === 1 ? '🎭' : '👤'}
                          </span>
                          <div>
                            <span className="text-sm font-medium text-gray-900 dark:text-white">
                              {availableRole.display_name}
                            </span>
                            <span className="ml-2 inline-flex items-center px-2 py-0.5 rounded-full text-xs font-medium bg-gray-100 text-gray-800">
                              Level {availableRole.hierarchy_level}
                            </span>
                          </div>
                        </div>
                      </label>
                    ))}
                  </div>
                  
                  <p className="mt-2 text-xs text-gray-500">
                    Select a parent role to inherit permissions from, or choose no parent to make this a root role.
                  </p>
                </div>

                {/* Current Status */}
                <div className="bg-gray-50 dark:bg-gray-700 p-3 rounded-md">
                  <h4 className="text-sm font-medium text-gray-700 dark:text-gray-300 mb-1">Current Status</h4>
                  <p className="text-sm text-gray-600 dark:text-gray-400">
                    {role.parent_role_id ? 
                      `Currently inherits from: ${role.parent_role?.display_name || 'Unknown'}` :
                      'Currently a root role (no parent)'
                    }
                  </p>
                  <p className="text-xs text-gray-500 dark:text-gray-400 mt-1">
                    Hierarchy Level: {role.hierarchy_level}
                  </p>
                </div>
              </div>
            </div>

            <div className="bg-gray-50 dark:bg-gray-700 px-4 py-3 sm:px-6 sm:flex sm:flex-row-reverse">
              <button
                type="submit"
                className="w-full inline-flex justify-center rounded-md border border-transparent shadow-sm px-4 py-2 bg-indigo-600 text-base font-medium text-white hover:bg-indigo-700 focus:outline-none focus:ring-2 focus:ring-offset-2 focus:ring-indigo-500 sm:ml-3 sm:w-auto sm:text-sm"
              >
                Update Parent
              </button>
              <button
                type="button"
                onClick={onClose}
                className="mt-3 w-full inline-flex justify-center rounded-md border border-gray-300 dark:border-gray-600 shadow-sm px-4 py-2 bg-white dark:bg-gray-600 text-base font-medium text-gray-700 dark:text-gray-300 hover:bg-gray-50 dark:hover:bg-gray-500 focus:outline-none focus:ring-2 focus:ring-offset-2 focus:ring-indigo-500 sm:mt-0 sm:ml-3 sm:w-auto sm:text-sm"
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

export default SimplifiedRoleManager;