import React, { useState } from 'react';
import { rbacService, Role } from '../../services/rbacService';

interface SimpleRoleHierarchyProps {
  roles: Role[];
  onRoleUpdate?: () => void;
}

const SimpleRoleHierarchy: React.FC<SimpleRoleHierarchyProps> = ({
  roles,
  onRoleUpdate
}) => {
  const [selectedRole, setSelectedRole] = useState<Role | null>(null);
  const [showUpdateModal, setShowUpdateModal] = useState(false);
  const [availableParents, setAvailableParents] = useState<Role[]>([]);
  const [loading, setLoading] = useState(false);

  const loadAvailableParents = async (roleId: number) => {
    try {
      setLoading(true);
      const response = await rbacService.getEligibleParentRoles(roleId);
      setAvailableParents(Array.isArray(response.data) ? response.data : []);
    } catch (error) {
      console.error('Failed to load available parents:', error);
      setAvailableParents(roles?.filter(r => r.id !== roleId && !r.is_system) || []);
    } finally {
      setLoading(false);
    }
  };

  const handleRoleClick = (role: Role) => {
    setSelectedRole(role);
  };

  const handleUpdateParent = async (role: Role) => {
    setSelectedRole(role);
    setShowUpdateModal(true);
    await loadAvailableParents(role.id);
  };

  const handleParentUpdate = async (newParentId: number | null) => {
    if (!selectedRole) return;

    try {
      setLoading(true);
      await rbacService.setRoleParent(selectedRole.id, newParentId);
      setShowUpdateModal(false);
      setSelectedRole(null);
      onRoleUpdate?.();
    } catch (error: any) {
      console.error('Failed to update parent:', error);
      alert(`Failed to update parent: ${error.response?.data?.message || error.message || 'Unknown error'}`);
    } finally {
      setLoading(false);
    }
  };

  const getParentName = (role: Role) => {
    if (!role.parent_role_id) return '-';
    const parent = roles.find(r => r.id === role.parent_role_id);
    return parent?.display_name || 'Unknown';
  };

  const getLevelIcon = (level: number) => {
    if (level === 0) return '👑';
    if (level === 1) return '🎭';
    return '👤';
  };

  // Sort roles by hierarchy level, then by name
  const sortedRoles = [...roles].sort((a, b) => {
    if (a.hierarchy_level !== b.hierarchy_level) {
      return a.hierarchy_level - b.hierarchy_level;
    }
    return a.display_name.localeCompare(b.display_name);
  });

  return (
    <div className="space-y-6">
      <div className="bg-white shadow rounded-lg">
        <div className="px-4 py-5 sm:px-6 border-b border-gray-200">
          <h3 className="text-lg font-medium text-gray-900">Role Hierarchy</h3>
          <p className="mt-1 text-sm text-gray-500">
            Click on a role to view details. Use "Update Parent" to change hierarchy.
          </p>
        </div>

        <div className="overflow-hidden">
          <table className="min-w-full divide-y divide-gray-200">
            <thead className="bg-gray-50">
              <tr>
                <th className="px-6 py-3 text-left text-xs font-medium text-gray-500 uppercase tracking-wider">
                  Role
                </th>
                <th className="px-6 py-3 text-left text-xs font-medium text-gray-500 uppercase tracking-wider">
                  Level
                </th>
                <th className="px-6 py-3 text-left text-xs font-medium text-gray-500 uppercase tracking-wider">
                  Parent Role
                </th>
                <th className="px-6 py-3 text-left text-xs font-medium text-gray-500 uppercase tracking-wider">
                  Type
                </th>
                <th className="px-6 py-3 text-left text-xs font-medium text-gray-500 uppercase tracking-wider">
                  Actions
                </th>
              </tr>
            </thead>
            <tbody className="bg-white divide-y divide-gray-200">
              {sortedRoles.map((role) => (
                <tr 
                  key={role.id}
                  className={`hover:bg-gray-50 cursor-pointer ${
                    selectedRole?.id === role.id ? 'bg-indigo-50' : ''
                  }`}
                  onClick={() => handleRoleClick(role)}
                >
                  <td className="px-6 py-4 whitespace-nowrap">
                    <div className="flex items-center">
                      <span className="mr-3 text-lg">
                        {getLevelIcon(role.hierarchy_level)}
                      </span>
                      <div>
                        <div className="text-sm font-medium text-gray-900">
                          {role.display_name}
                        </div>
                        <div className="text-sm text-gray-500">
                          {role.name}
                        </div>
                      </div>
                    </div>
                  </td>
                  <td className="px-6 py-4 whitespace-nowrap">
                    <span className={`inline-flex items-center px-2.5 py-0.5 rounded-full text-xs font-medium ${
                      role.hierarchy_level === 0 ? 'bg-yellow-100 text-yellow-800' :
                      role.hierarchy_level === 1 ? 'bg-purple-100 text-purple-800' :
                      'bg-gray-100 text-gray-800'
                    }`}>
                      Level {role.hierarchy_level}
                    </span>
                  </td>
                  <td className="px-6 py-4 whitespace-nowrap text-sm text-gray-900">
                    {getParentName(role)}
                  </td>
                  <td className="px-6 py-4 whitespace-nowrap">
                    <span className={`inline-flex items-center px-2.5 py-0.5 rounded-full text-xs font-medium ${
                      role.is_system ? 'bg-red-100 text-red-800' : 'bg-blue-100 text-blue-800'
                    }`}>
                      {role.is_system ? 'System' : 'Custom'}
                    </span>
                  </td>
                  <td className="px-6 py-4 whitespace-nowrap text-sm font-medium">
                    <button
                      onClick={(e) => {
                        e.stopPropagation();
                        handleUpdateParent(role);
                      }}
                      disabled={role.is_system}
                      className={`${
                        role.is_system 
                          ? 'text-gray-400 cursor-not-allowed' 
                          : 'text-indigo-600 hover:text-indigo-900'
                      }`}
                    >
                      {role.is_system ? 'Protected' : 'Update Parent'}
                    </button>
                  </td>
                </tr>
              ))}
            </tbody>
          </table>
        </div>
      </div>

      {/* Role Details Panel */}
      {selectedRole && (
        <div className="bg-white shadow rounded-lg p-6">
          <h4 className="text-lg font-medium text-gray-900 mb-4">
            Role Details: {selectedRole.display_name}
          </h4>
          <div className="grid grid-cols-2 gap-4">
            <div>
              <label className="block text-sm font-medium text-gray-700">Internal Name</label>
              <p className="text-sm text-gray-900">{selectedRole.name}</p>
            </div>
            <div>
              <label className="block text-sm font-medium text-gray-700">Hierarchy Level</label>
              <p className="text-sm text-gray-900">
                {getLevelIcon(selectedRole.hierarchy_level)} Level {selectedRole.hierarchy_level}
              </p>
            </div>
            <div>
              <label className="block text-sm font-medium text-gray-700">Parent Role</label>
              <p className="text-sm text-gray-900">
                {getParentName(selectedRole)}
              </p>
            </div>
            <div>
              <label className="block text-sm font-medium text-gray-700">Type</label>
              <p className="text-sm text-gray-900">
                {selectedRole.is_system ? 'System Role' : 'Custom Role'}
              </p>
            </div>
            <div>
              <label className="block text-sm font-medium text-gray-700">Status</label>
              <p className="text-sm text-gray-900">
                {selectedRole.is_active ? 'Active' : 'Inactive'}
              </p>
            </div>
            <div>
              <label className="block text-sm font-medium text-gray-700">Child Roles</label>
              <p className="text-sm text-gray-900">
                {roles.filter(r => r.parent_role_id === selectedRole.id).length} children
              </p>
            </div>
            <div className="col-span-2">
              <label className="block text-sm font-medium text-gray-700">Description</label>
              <p className="text-sm text-gray-900">
                {selectedRole.description || 'No description provided'}
              </p>
            </div>
          </div>
        </div>
      )}

      {/* Update Parent Modal */}
      {showUpdateModal && selectedRole && (
        <div className="fixed inset-0 z-50 overflow-y-auto">
          <div className="flex items-center justify-center min-h-screen pt-4 px-4 pb-20 text-center sm:block sm:p-0">
            <div className="fixed inset-0 bg-gray-500 bg-opacity-75 transition-opacity" onClick={() => setShowUpdateModal(false)}></div>
            
            <div className="inline-block align-bottom bg-white rounded-lg text-left overflow-hidden shadow-xl transform transition-all sm:my-8 sm:align-middle sm:max-w-lg sm:w-full">
              <div className="bg-white px-4 pt-5 pb-4 sm:p-6 sm:pb-4">
                <h3 className="text-lg font-medium text-gray-900 mb-4">
                  Update Parent for "{selectedRole.display_name}"
                </h3>
                
                {loading ? (
                  <div className="flex justify-center py-4">
                    <div className="animate-spin rounded-full h-6 w-6 border-b-2 border-indigo-500"></div>
                  </div>
                ) : (
                  <div className="space-y-3">
                    <label className="flex items-center p-3 border rounded-lg hover:bg-gray-50 cursor-pointer">
                      <input
                        type="radio"
                        name="parentRole"
                        value=""
                        onChange={() => handleParentUpdate(null)}
                        className="mr-3"
                      />
                      <div className="flex items-center">
                        <span className="mr-3 text-lg">👑</span>
                        <div>
                          <span className="text-sm font-medium text-gray-900">
                            Root Level (No Parent)
                          </span>
                          <p className="text-xs text-gray-500">Make this a top-level role</p>
                        </div>
                      </div>
                    </label>
                    
                    {availableParents?.map(parent => (
                      <label key={parent.id} className="flex items-center p-3 border rounded-lg hover:bg-gray-50 cursor-pointer">
                        <input
                          type="radio"
                          name="parentRole"
                          value={parent.id}
                          onChange={() => handleParentUpdate(parent.id)}
                          className="mr-3"
                        />
                        <div className="flex items-center">
                          <span className="mr-3 text-lg">
                            {getLevelIcon(parent.hierarchy_level)}
                          </span>
                          <div>
                            <span className="text-sm font-medium text-gray-900">
                              {parent.display_name}
                            </span>
                            <p className="text-xs text-gray-500">
                              Level {parent.hierarchy_level} • {parent.description || 'No description'}
                            </p>
                          </div>
                        </div>
                      </label>
                    ))}
                  </div>
                )}
              </div>

              <div className="bg-gray-50 px-4 py-3 sm:px-6 sm:flex sm:flex-row-reverse">
                <button
                  type="button"
                  onClick={() => setShowUpdateModal(false)}
                  className="mt-3 w-full inline-flex justify-center rounded-md border border-gray-300 shadow-sm px-4 py-2 bg-white text-base font-medium text-gray-700 hover:bg-gray-50 focus:outline-none focus:ring-2 focus:ring-offset-2 focus:ring-indigo-500 sm:mt-0 sm:w-auto sm:text-sm"
                >
                  Cancel
                </button>
              </div>
            </div>
          </div>
        </div>
      )}
    </div>
  );
};

export default SimpleRoleHierarchy;