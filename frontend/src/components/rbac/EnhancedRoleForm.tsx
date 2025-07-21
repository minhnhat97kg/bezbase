import React, { useState, useEffect } from 'react';
import { rbacService, Role } from '../../services/rbacService';

interface EnhancedRoleFormProps {
  role?: Role; // For editing
  onSubmit: (data: RoleFormData) => void;
  onClose: () => void;
  availableRoles?: Role[];
}

interface RoleFormData {
  name: string;
  display_name: string;
  description: string;
  parent_role_id?: number | null;
}

const EnhancedRoleForm: React.FC<EnhancedRoleFormProps> = ({
  role,
  onSubmit,
  onClose,
  availableRoles = []
}) => {
  const [formData, setFormData] = useState<RoleFormData>({
    name: role?.name || '',
    display_name: role?.display_name || '',
    description: role?.description || '',
    parent_role_id: role?.parent_role_id || null,
  });

  const [eligibleParents, setEligibleParents] = useState<Role[]>([]);
  const [loading, setLoading] = useState(false);
  const [errors, setErrors] = useState<Record<string, string>>({});

  const isEditing = !!role;

  // Helper function to check for circular dependencies
  const isCircularDependency = (childRoleId: number, parentRoleId: number, roles: Role[]): boolean => {
    const visited = new Set<number>();
    
    const checkAncestors = (currentRoleId: number): boolean => {
      if (visited.has(currentRoleId)) {
        return false; // Prevent infinite loops
      }
      visited.add(currentRoleId);
      
      const currentRole = roles.find(r => r.id === currentRoleId);
      if (!currentRole || !currentRole.parent_role_id) {
        return false;
      }
      
      if (currentRole.parent_role_id === childRoleId) {
        return true; // Found circular dependency
      }
      
      return checkAncestors(currentRole.parent_role_id);
    };
    
    return checkAncestors(parentRoleId);
  };

  useEffect(() => {
    if (isEditing && role) {
      loadEligibleParents();
    } else {
      // For new roles, all non-system roles can be parents
      setEligibleParents(availableRoles.filter(r => !r.is_system));
    }
  }, [role, availableRoles, isEditing]);

  const loadEligibleParents = async () => {
    if (!role) return;

    try {
      setLoading(true);
      const response = await rbacService.getEligibleParentRoles(role.id);
      setEligibleParents(response.data);
    } catch (error) {
      console.error('Failed to load eligible parent roles:', error);
      // Fallback to filtering out the current role and system roles
      setEligibleParents(availableRoles.filter(r => r.id !== role.id && !r.is_system));
    } finally {
      setLoading(false);
    }
  };

  const validateForm = (): boolean => {
    const newErrors: Record<string, string> = {};

    if (!formData.name.trim()) {
      newErrors.name = 'Role name is required';
    } else if (!/^[a-z0-9_-]+$/.test(formData.name)) {
      newErrors.name = 'Role name must be lowercase letters, numbers, underscores, or hyphens only';
    }

    if (!formData.display_name.trim()) {
      newErrors.display_name = 'Display name is required';
    }

    // Check for potential circular reference and hierarchy constraints
    if (formData.parent_role_id) {
      const selectedParent = eligibleParents.find(r => r.id === formData.parent_role_id);
      if (!selectedParent) {
        newErrors.parent_role_id = 'Selected parent role is not available';
      } else {
        // Check hierarchy level constraints
        if (role && selectedParent.hierarchy_level >= role.hierarchy_level) {
          newErrors.parent_role_id = 'Parent role must have a lower hierarchy level';
        }
        
        // Check depth limit
        if (selectedParent.hierarchy_level >= 9) {
          newErrors.parent_role_id = 'Cannot set parent: would exceed maximum hierarchy depth (10 levels)';
        }
        
        // Additional frontend circular dependency check
        if (role && isCircularDependency(role.id, formData.parent_role_id, eligibleParents)) {
          newErrors.parent_role_id = 'Cannot set parent: would create circular dependency';
        }
      }
    }

    setErrors(newErrors);
    return Object.keys(newErrors).length === 0;
  };

  const handleSubmit = (e: React.FormEvent) => {
    e.preventDefault();
    
    if (!validateForm()) {
      return;
    }

    onSubmit(formData);
  };

  const handleInputChange = (field: keyof RoleFormData, value: string | number | null) => {
    setFormData(prev => ({ ...prev, [field]: value }));
    
    // Clear field error when user starts typing
    if (errors[field]) {
      setErrors(prev => ({ ...prev, [field]: '' }));
    }
  };

  const getHierarchyPreview = () => {
    if (!formData.parent_role_id) {
      return { level: 0, preview: 'Root role (Level 0)' };
    }

    const parentRole = eligibleParents.find(r => r.id === formData.parent_role_id);
    if (!parentRole) {
      return { level: 1, preview: 'Child role (Level 1)' };
    }

    const newLevel = parentRole.hierarchy_level + 1;
    return {
      level: newLevel,
      preview: `Child of "${parentRole.display_name}" (Level ${newLevel})`
    };
  };

  const hierarchyPreview = getHierarchyPreview();

  return (
    <div className="fixed inset-0 z-50 overflow-y-auto">
      <div className="flex items-center justify-center min-h-screen pt-4 px-4 pb-20 text-center sm:block sm:p-0">
        <div className="fixed inset-0 bg-gray-500 bg-opacity-75 transition-opacity" onClick={onClose}></div>
        
        <div className="inline-block align-bottom bg-white dark:bg-gray-800 rounded-lg text-left overflow-hidden shadow-xl transform transition-all sm:my-8 sm:align-middle sm:max-w-2xl sm:w-full">
          <form onSubmit={handleSubmit}>
            <div className="bg-white dark:bg-gray-800 px-4 pt-5 pb-4 sm:p-6 sm:pb-4">
              <h3 className="text-lg font-medium text-gray-900 dark:text-white mb-6">
                {isEditing ? `Edit Role: ${role.display_name}` : 'Create New Role'}
              </h3>
              
              <div className="space-y-6">
                {/* Basic Role Information */}
                <div className="grid grid-cols-1 gap-4 sm:grid-cols-2">
                  <div>
                    <label className="block text-sm font-medium text-gray-700 dark:text-gray-300">
                      Role Name (Internal) *
                    </label>
                    <input
                      type="text"
                      value={formData.name}
                      onChange={(e) => handleInputChange('name', e.target.value)}
                      className={`mt-1 block w-full border rounded-md shadow-sm focus:ring-indigo-500 focus:border-indigo-500 dark:bg-gray-700 dark:text-white ${
                        errors.name ? 'border-red-300' : 'border-gray-300 dark:border-gray-600'
                      }`}
                      placeholder="e.g., content_editor, team_manager"
                      disabled={isEditing && role?.is_system}
                    />
                    {errors.name && (
                      <p className="mt-1 text-sm text-red-600">{errors.name}</p>
                    )}
                    <p className="mt-1 text-xs text-gray-500">
                      Lowercase, no spaces, used for internal identification
                    </p>
                  </div>

                  <div>
                    <label className="block text-sm font-medium text-gray-700 dark:text-gray-300">
                      Display Name *
                    </label>
                    <input
                      type="text"
                      value={formData.display_name}
                      onChange={(e) => handleInputChange('display_name', e.target.value)}
                      className={`mt-1 block w-full border rounded-md shadow-sm focus:ring-indigo-500 focus:border-indigo-500 dark:bg-gray-700 dark:text-white ${
                        errors.display_name ? 'border-red-300' : 'border-gray-300 dark:border-gray-600'
                      }`}
                      placeholder="e.g., Content Editor, Team Manager"
                    />
                    {errors.display_name && (
                      <p className="mt-1 text-sm text-red-600">{errors.display_name}</p>
                    )}
                    <p className="mt-1 text-xs text-gray-500">
                      Human-readable name shown in the interface
                    </p>
                  </div>
                </div>

                {/* Description */}
                <div>
                  <label className="block text-sm font-medium text-gray-700 dark:text-gray-300">
                    Description
                  </label>
                  <textarea
                    value={formData.description}
                    onChange={(e) => handleInputChange('description', e.target.value)}
                    rows={3}
                    className="mt-1 block w-full border-gray-300 dark:border-gray-600 rounded-md shadow-sm focus:ring-indigo-500 focus:border-indigo-500 dark:bg-gray-700 dark:text-white"
                    placeholder="Describe the purpose and responsibilities of this role..."
                  />
                </div>

                {/* Parent Role Selection */}
                <div>
                  <label className="block text-sm font-medium text-gray-700 dark:text-gray-300 mb-3">
                    Role Hierarchy
                  </label>
                  
                  {loading ? (
                    <div className="flex justify-center py-4">
                      <div className="animate-spin rounded-full h-6 w-6 border-b-2 border-indigo-500"></div>
                    </div>
                  ) : (
                    <div className="space-y-3">
                      <div className="space-y-2">
                        <label className="flex items-center p-3 border rounded-lg hover:bg-gray-50 dark:hover:bg-gray-700 cursor-pointer">
                          <input
                            type="radio"
                            value=""
                            checked={formData.parent_role_id === null}
                            onChange={() => handleInputChange('parent_role_id', null)}
                            className="mr-3"
                            disabled={isEditing && role?.is_system}
                          />
                          <div className="flex items-center">
                            <span className="mr-3 text-lg">👑</span>
                            <div>
                              <span className="text-sm font-medium text-gray-900 dark:text-white">
                                Root Role (No Parent)
                              </span>
                              <p className="text-xs text-gray-500">
                                This role will not inherit permissions from any other role
                              </p>
                            </div>
                          </div>
                        </label>
                        
                        {eligibleParents.map(parentRole => (
                          <label 
                            key={parentRole.id} 
                            className="flex items-center p-3 border rounded-lg hover:bg-gray-50 dark:hover:bg-gray-700 cursor-pointer"
                          >
                            <input
                              type="radio"
                              value={parentRole.id}
                              checked={formData.parent_role_id === parentRole.id}
                              onChange={() => handleInputChange('parent_role_id', parentRole.id)}
                              className="mr-3"
                              disabled={isEditing && role?.is_system}
                            />
                            <div className="flex items-center flex-1">
                              <span className="mr-3 text-lg">
                                {parentRole.hierarchy_level === 0 ? '👑' : parentRole.hierarchy_level === 1 ? '🎭' : '👤'}
                              </span>
                              <div className="flex-1">
                                <div className="flex items-center">
                                  <span className="text-sm font-medium text-gray-900 dark:text-white">
                                    {parentRole.display_name}
                                  </span>
                                  <span className="ml-2 inline-flex items-center px-2 py-0.5 rounded-full text-xs font-medium bg-gray-100 text-gray-800">
                                    Level {parentRole.hierarchy_level}
                                  </span>
                                </div>
                                <p className="text-xs text-gray-500 mt-1">
                                  {parentRole.description || 'No description'}
                                </p>
                              </div>
                            </div>
                          </label>
                        ))}
                      </div>
                      
                      {errors.parent_role_id && (
                        <p className="text-sm text-red-600">{errors.parent_role_id}</p>
                      )}
                    </div>
                  )}
                </div>

                {/* Hierarchy Preview */}
                <div className="bg-indigo-50 dark:bg-indigo-900/20 p-4 rounded-lg">
                  <h4 className="text-sm font-medium text-indigo-800 dark:text-indigo-200 mb-2">
                    Hierarchy Preview
                  </h4>
                  <div className="flex items-center">
                    <span className="mr-3 text-lg">
                      {hierarchyPreview.level === 0 ? '👑' : hierarchyPreview.level === 1 ? '🎭' : '👤'}
                    </span>
                    <div>
                      <p className="text-sm text-indigo-700 dark:text-indigo-300">
                        {hierarchyPreview.preview}
                      </p>
                      <p className="text-xs text-indigo-600 dark:text-indigo-400 mt-1">
                        {formData.parent_role_id 
                          ? 'This role will inherit all permissions from its parent role'
                          : 'This role will only have directly assigned permissions'
                        }
                      </p>
                    </div>
                  </div>
                </div>

                {/* System Role Warning */}
                {isEditing && role?.is_system && (
                  <div className="bg-red-50 dark:bg-red-900/20 p-4 rounded-lg">
                    <div className="flex">
                      <div className="flex-shrink-0">
                        <span className="text-red-400 text-lg">⚠️</span>
                      </div>
                      <div className="ml-3">
                        <h3 className="text-sm font-medium text-red-800 dark:text-red-200">
                          System Role
                        </h3>
                        <p className="text-sm text-red-700 dark:text-red-300 mt-1">
                          This is a system role. Some fields are read-only to prevent system integrity issues.
                        </p>
                      </div>
                    </div>
                  </div>
                )}
              </div>
            </div>

            <div className="bg-gray-50 dark:bg-gray-700 px-4 py-3 sm:px-6 sm:flex sm:flex-row-reverse">
              <button
                type="submit"
                className="w-full inline-flex justify-center rounded-md border border-transparent shadow-sm px-4 py-2 bg-indigo-600 text-base font-medium text-white hover:bg-indigo-700 focus:outline-none focus:ring-2 focus:ring-offset-2 focus:ring-indigo-500 sm:ml-3 sm:w-auto sm:text-sm"
              >
                {isEditing ? 'Update Role' : 'Create Role'}
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

export default EnhancedRoleForm;