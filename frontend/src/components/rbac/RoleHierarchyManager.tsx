import React, { useState, useEffect } from 'react';
import { rbacService, Role } from '../../services/rbacService';

interface RoleHierarchyManagerProps {
  roles: Role[];
  onRoleUpdate?: () => void;
}

interface HierarchyNode {
  role: Role;
  children: HierarchyNode[];
  level: number;
}

interface MoveOperation {
  roleId: number;
  newParentId: number | null;
  roleName: string;
  newParentName: string;
}

const RoleHierarchyManager: React.FC<RoleHierarchyManagerProps> = ({
  roles,
  onRoleUpdate
}) => {
  const [hierarchyTree, setHierarchyTree] = useState<HierarchyNode[]>([]);
  const [selectedRole, setSelectedRole] = useState<Role | null>(null);
  const [availableParents, setAvailableParents] = useState<Role[]>([]);
  const [pendingMoves, setPendingMoves] = useState<MoveOperation[]>([]);
  const [loading, setLoading] = useState(false);
  const [expandedNodes, setExpandedNodes] = useState<Set<number>>(new Set());
  const [draggedRole, setDraggedRole] = useState<Role | null>(null);
  const [dropTarget, setDropTarget] = useState<number | null>(null);

  useEffect(() => {
    buildHierarchyTree();
    // Auto-expand all nodes initially
    const allRoleIds = new Set(roles.map(r => r.id));
    setExpandedNodes(allRoleIds);
  }, [roles]);

  const buildHierarchyTree = () => {
    const rootRoles = roles.filter(role => !role.parent_role_id);
    
    const buildNode = (role: Role, level: number = 0): HierarchyNode => {
      const children = roles
        .filter(r => r.parent_role_id === role.id)
        .map(child => buildNode(child, level + 1));
      
      return {
        role,
        children,
        level
      };
    };

    const tree = rootRoles.map(role => buildNode(role));
    setHierarchyTree(tree);
  };

  const loadAvailableParents = async (roleId: number) => {
    try {
      setLoading(true);
      const response = await rbacService.getEligibleParentRoles(roleId);
      setAvailableParents(response.data);
    } catch (error) {
      console.error('Failed to load available parents:', error);
      // Fallback: filter out the role itself and system roles
      const filtered = roles.filter(r => r.id !== roleId && !r.is_system);
      setAvailableParents(filtered);
    } finally {
      setLoading(false);
    }
  };

  const handleRoleSelect = async (role: Role) => {
    setSelectedRole(role);
    await loadAvailableParents(role.id);
  };

  const handleMoveRole = (role: Role, newParentId: number | null) => {
    const newParent = newParentId ? roles.find(r => r.id === newParentId) : null;
    
    const moveOp: MoveOperation = {
      roleId: role.id,
      newParentId,
      roleName: role.display_name,
      newParentName: newParent?.display_name || 'Root Level'
    };

    setPendingMoves(prev => {
      // Remove any existing move for this role
      const filtered = prev.filter(m => m.roleId !== role.id);
      return [...filtered, moveOp];
    });
  };

  const removePendingMove = (roleId: number) => {
    setPendingMoves(prev => prev.filter(m => m.roleId !== roleId));
  };

  const applyPendingMoves = async () => {
    setLoading(true);
    const errors: string[] = [];

    for (const move of pendingMoves) {
      try {
        await rbacService.setRoleParent(move.roleId, move.newParentId);
      } catch (error: any) {
        errors.push(`Failed to move ${move.roleName}: ${error.message || 'Unknown error'}`);
      }
    }

    if (errors.length > 0) {
      alert(`Some moves failed:\n${errors.join('\n')}`);
    } else {
      alert(`Successfully applied ${pendingMoves.length} hierarchy changes!`);
    }

    setPendingMoves([]);
    setLoading(false);
    onRoleUpdate?.();
  };

  const toggleExpanded = (roleId: number) => {
    const newExpanded = new Set(expandedNodes);
    if (newExpanded.has(roleId)) {
      newExpanded.delete(roleId);
    } else {
      newExpanded.add(roleId);
    }
    setExpandedNodes(newExpanded);
  };

  const expandAll = () => {
    const allRoleIds = new Set(roles.map(r => r.id));
    setExpandedNodes(allRoleIds);
  };

  const collapseAll = () => {
    setExpandedNodes(new Set());
  };

  // Drag and Drop handlers
  const handleDragStart = (e: React.DragEvent, role: Role) => {
    if (role.is_system) {
      e.preventDefault();
      return;
    }
    setDraggedRole(role);
    e.dataTransfer.effectAllowed = 'move';
    e.dataTransfer.setData('text/plain', role.id.toString());
  };

  const handleDragOver = (e: React.DragEvent, targetRole: Role) => {
    e.preventDefault();
    if (draggedRole && !targetRole.is_system && draggedRole.id !== targetRole.id) {
      e.dataTransfer.dropEffect = 'move';
      setDropTarget(targetRole.id);
    }
  };

  const handleDragLeave = (e: React.DragEvent) => {
    e.preventDefault();
    setDropTarget(null);
  };

  const handleDrop = (e: React.DragEvent, targetRole: Role) => {
    e.preventDefault();
    setDropTarget(null);
    
    if (draggedRole && !targetRole.is_system && draggedRole.id !== targetRole.id) {
      // Check if this would be a valid move (prevent circular dependencies)
      const wouldCreateCircular = isCircularDependency(targetRole.id, draggedRole.id);
      
      if (wouldCreateCircular) {
        alert(`Cannot move "${draggedRole.display_name}" under "${targetRole.display_name}" - this would create a circular dependency.`);
        return;
      }
      
      handleMoveRole(draggedRole, targetRole.id);
    }
    setDraggedRole(null);
  };

  const handleDropToRoot = (e: React.DragEvent) => {
    e.preventDefault();
    setDropTarget(null);
    
    if (draggedRole) {
      handleMoveRole(draggedRole, null);
    }
    setDraggedRole(null);
  };

  // Helper function to check circular dependencies
  const isCircularDependency = (parentId: number, childId: number): boolean => {
    const visited = new Set<number>();
    
    const checkAncestors = (currentRoleId: number): boolean => {
      if (visited.has(currentRoleId)) {
        return false;
      }
      visited.add(currentRoleId);
      
      const currentRole = roles.find(r => r.id === currentRoleId);
      if (!currentRole || !currentRole.parent_role_id) {
        return false;
      }
      
      if (currentRole.parent_role_id === childId) {
        return true;
      }
      
      return checkAncestors(currentRole.parent_role_id);
    };
    
    return checkAncestors(parentId);
  };

  const renderHierarchyNode = (node: HierarchyNode, isRoot: boolean = false) => {
    const isExpanded = expandedNodes.has(node.role.id);
    const hasChildren = node.children.length > 0;
    const isSelected = selectedRole?.id === node.role.id;
    const hasPendingMove = pendingMoves.some(m => m.roleId === node.role.id);
    const isDraggedOver = dropTarget === node.role.id;
    const isDragging = draggedRole?.id === node.role.id;

    return (
      <div key={node.role.id} className={`role-node ${isRoot ? 'border-l-4 border-indigo-500' : ''}`}>
        <div 
          className={`flex items-center p-3 rounded-lg cursor-pointer transition-all ${
            isSelected 
              ? 'bg-indigo-100 border-indigo-300 border-2' 
              : hasPendingMove
                ? 'bg-yellow-50 border-yellow-300 border'
                : isDraggedOver
                  ? 'bg-green-100 border-green-300 border-2 border-dashed'
                  : isDragging
                    ? 'opacity-50 bg-gray-100'
                    : 'hover:bg-gray-50 border border-transparent'
          }`}
          draggable={!node.role.is_system}
          onDragStart={(e) => handleDragStart(e, node.role)}
          onDragOver={(e) => handleDragOver(e, node.role)}
          onDragLeave={handleDragLeave}
          onDrop={(e) => handleDrop(e, node.role)}
          onClick={() => handleRoleSelect(node.role)}
          title={node.role.is_system ? 'System roles cannot be moved' : 'Click to select, drag to move'}
        >
          {/* Expand/Collapse Button */}
          {hasChildren && (
            <button
              onClick={(e) => {
                e.stopPropagation();
                toggleExpanded(node.role.id);
              }}
              className="mr-2 p-1 hover:bg-gray-200 rounded transition-colors"
            >
              {isExpanded ? (
                <svg className="w-4 h-4" fill="currentColor" viewBox="0 0 20 20">
                  <path fillRule="evenodd" d="M5.293 7.293a1 1 0 011.414 0L10 10.586l3.293-3.293a1 1 0 111.414 1.414l-4 4a1 1 0 01-1.414 0l-4-4a1 1 0 010-1.414z" clipRule="evenodd" />
                </svg>
              ) : (
                <svg className="w-4 h-4" fill="currentColor" viewBox="0 0 20 20">
                  <path fillRule="evenodd" d="M7.293 14.707a1 1 0 010-1.414L10.586 10 7.293 6.707a1 1 0 011.414-1.414l4 4a1 1 0 010 1.414l-4 4a1 1 0 01-1.414 0z" clipRule="evenodd" />
                </svg>
              )}
            </button>
          )}

          {/* Role Icon */}
          <span className="mr-3 text-lg">
            {node.level === 0 ? '👑' : node.level === 1 ? '🎭' : '👤'}
          </span>

          {/* Role Info */}
          <div className="flex-1">
            <div className="flex items-center space-x-2">
              <span className="font-medium text-gray-900">
                {node.role.display_name}
              </span>
              
              <span className={`inline-flex items-center px-2 py-0.5 rounded-full text-xs font-medium ${
                node.level === 0 ? 'bg-yellow-100 text-yellow-800' :
                node.level === 1 ? 'bg-purple-100 text-purple-800' :
                'bg-gray-100 text-gray-800'
              }`}>
                Level {node.level}
              </span>

              {node.role.is_system && (
                <span className="inline-flex items-center px-2 py-0.5 rounded-full text-xs font-medium bg-red-100 text-red-800">
                  System
                </span>
              )}

              {hasPendingMove && (
                <span className="inline-flex items-center px-2 py-0.5 rounded-full text-xs font-medium bg-yellow-100 text-yellow-800">
                  📝 Pending Move
                </span>
              )}

              {!node.role.is_system && (
                <span className="inline-flex items-center px-2 py-0.5 rounded-full text-xs font-medium bg-blue-100 text-blue-800">
                  🔄 Draggable
                </span>
              )}
            </div>
            
            <div className="text-sm text-gray-500 mt-1">
              {node.role.description || 'No description'}
            </div>
          </div>

          {/* Child Count */}
          {hasChildren && (
            <span className="ml-2 text-sm text-gray-500 bg-gray-100 px-2 py-1 rounded-full">
              {node.children.length} child{node.children.length !== 1 ? 'ren' : ''}
            </span>
          )}
        </div>

        {/* Render Children */}
        {isExpanded && hasChildren && (
          <div className="ml-8 mt-2 border-l-2 border-gray-200 pl-4">
            {node.children.map(child => renderHierarchyNode(child))}
          </div>
        )}
      </div>
    );
  };

  return (
    <div className="role-hierarchy-manager">
      <div className="bg-white shadow rounded-lg">
        <div className="px-4 py-5 sm:px-6 border-b border-gray-200">
          <div className="flex justify-between items-center">
            <div>
              <h3 className="text-lg font-medium text-gray-900">Role Hierarchy Management</h3>
              <p className="mt-1 text-sm text-gray-500">
                Select roles to edit, or drag & drop to reorganize the hierarchy. Changes are applied in batch.
              </p>
            </div>
            
            <div className="flex space-x-2">
              <button
                onClick={expandAll}
                className="bg-gray-100 hover:bg-gray-200 text-gray-700 px-3 py-1 rounded-md text-sm"
              >
                Expand All
              </button>
              <button
                onClick={collapseAll}
                className="bg-gray-100 hover:bg-gray-200 text-gray-700 px-3 py-1 rounded-md text-sm"
              >
                Collapse All
              </button>
            </div>
          </div>
        </div>

        <div className="p-6">
          <div className="grid grid-cols-1 lg:grid-cols-3 gap-6">
            {/* Hierarchy Tree */}
            <div className="lg:col-span-2">
              <div className="flex justify-between items-center mb-4">
                <h4 className="text-sm font-medium text-gray-700">Role Hierarchy Tree</h4>
                {draggedRole && (
                  <div 
                    className="px-3 py-2 bg-yellow-100 border-2 border-dashed border-yellow-300 rounded-lg cursor-pointer hover:bg-yellow-200 transition-colors"
                    onDragOver={(e) => { e.preventDefault(); e.dataTransfer.dropEffect = 'move'; }}
                    onDrop={handleDropToRoot}
                  >
                    <span className="text-sm text-yellow-800">
                      👑 Drop here to make "{draggedRole.display_name}" a root role
                    </span>
                  </div>
                )}
              </div>
              
              {hierarchyTree.length === 0 ? (
                <div className="text-center py-8">
                  <span className="text-gray-400 text-4xl mb-4 block">🎭</span>
                  <p className="text-gray-500">No roles found</p>
                </div>
              ) : (
                <div className="space-y-2 max-h-96 overflow-y-auto">
                  {hierarchyTree.map(node => renderHierarchyNode(node, true))}
                </div>
              )}
            </div>

            {/* Role Editor Panel */}
            <div className="lg:col-span-1">
              <h4 className="text-sm font-medium text-gray-700 mb-4">Role Editor</h4>
              
              {!selectedRole ? (
                <div className="bg-gray-50 rounded-lg p-6 text-center">
                  <span className="text-gray-400 text-2xl mb-2 block">👆</span>
                  <p className="text-gray-500 text-sm">Select a role from the hierarchy to edit its position</p>
                </div>
              ) : (
                <div className="space-y-4">
                  {/* Selected Role Info */}
                  <div className="bg-indigo-50 rounded-lg p-4">
                    <h5 className="font-medium text-indigo-900 mb-2">Selected Role</h5>
                    <p className="text-sm text-indigo-800">{selectedRole.display_name}</p>
                    <p className="text-xs text-indigo-600 mt-1">
                      Current Level: {selectedRole.hierarchy_level}
                    </p>
                    {selectedRole.parent_role_id && (
                      <p className="text-xs text-indigo-600">
                        Current Parent: {roles.find(r => r.id === selectedRole.parent_role_id)?.display_name}
                      </p>
                    )}
                  </div>

                  {/* Move Role Section */}
                  {selectedRole.is_system ? (
                    <div className="bg-red-50 rounded-lg p-4">
                      <p className="text-sm text-red-800">
                        ⚠️ System roles cannot be moved in the hierarchy
                      </p>
                    </div>
                  ) : (
                    <div>
                      <label className="block text-sm font-medium text-gray-700 mb-2">
                        Move to Parent Role
                      </label>
                      
                      {loading ? (
                        <div className="flex justify-center py-4">
                          <div className="animate-spin rounded-full h-5 w-5 border-b-2 border-indigo-500"></div>
                        </div>
                      ) : (
                        <div className="space-y-2 max-h-40 overflow-y-auto">
                          <label className="flex items-center p-2 border rounded hover:bg-gray-50 cursor-pointer">
                            <input
                              type="radio"
                              name="newParent"
                              value=""
                              onChange={() => handleMoveRole(selectedRole, null)}
                              className="mr-2"
                            />
                            <span className="text-sm">
                              <span className="mr-2">👑</span>
                              Root Level (No Parent)
                            </span>
                          </label>
                          
                          {availableParents.map(parent => (
                            <label key={parent.id} className="flex items-center p-2 border rounded hover:bg-gray-50 cursor-pointer">
                              <input
                                type="radio"
                                name="newParent"
                                value={parent.id}
                                onChange={() => handleMoveRole(selectedRole, parent.id)}
                                className="mr-2"
                              />
                              <div className="flex items-center text-sm">
                                <span className="mr-2">
                                  {parent.hierarchy_level === 0 ? '👑' : parent.hierarchy_level === 1 ? '🎭' : '👤'}
                                </span>
                                <div>
                                  <span className="font-medium">{parent.display_name}</span>
                                  <span className="text-xs text-gray-500 ml-1">(Level {parent.hierarchy_level})</span>
                                </div>
                              </div>
                            </label>
                          ))}
                        </div>
                      )}
                    </div>
                  )}
                </div>
              )}
            </div>
          </div>
        </div>

        {/* Pending Changes */}
        {pendingMoves.length > 0 && (
          <div className="border-t border-gray-200 bg-yellow-50 p-4">
            <div className="flex items-center justify-between">
              <div>
                <h4 className="text-sm font-medium text-yellow-800">
                  📝 Pending Changes ({pendingMoves.length})
                </h4>
                <div className="mt-2 space-y-1">
                  {pendingMoves.map((move, index) => (
                    <div key={index} className="flex items-center justify-between text-sm text-yellow-700">
                      <span>
                        Move "{move.roleName}" to "{move.newParentName}"
                      </span>
                      <button
                        onClick={() => removePendingMove(move.roleId)}
                        className="text-yellow-600 hover:text-yellow-800 ml-2"
                      >
                        ✕
                      </button>
                    </div>
                  ))}
                </div>
              </div>
              
              <div className="flex space-x-2">
                <button
                  onClick={() => setPendingMoves([])}
                  className="bg-gray-200 hover:bg-gray-300 text-gray-700 px-3 py-1 rounded-md text-sm"
                  disabled={loading}
                >
                  Clear All
                </button>
                <button
                  onClick={applyPendingMoves}
                  className="bg-yellow-600 hover:bg-yellow-700 text-white px-4 py-1 rounded-md text-sm"
                  disabled={loading}
                >
                  {loading ? 'Applying...' : 'Apply Changes'}
                </button>
              </div>
            </div>
          </div>
        )}
      </div>

      {/* Quick Stats */}
      <div className="mt-6 grid grid-cols-1 md:grid-cols-4 gap-4">
        <div className="bg-blue-50 rounded-lg p-4">
          <div className="flex items-center">
            <span className="text-blue-600 text-lg mr-3">👑</span>
            <div>
              <p className="text-sm font-medium text-blue-800">Root Roles</p>
              <p className="text-lg font-semibold text-blue-900">
                {roles.filter(r => !r.parent_role_id).length}
              </p>
            </div>
          </div>
        </div>

        <div className="bg-purple-50 rounded-lg p-4">
          <div className="flex items-center">
            <span className="text-purple-600 text-lg mr-3">🎭</span>
            <div>
              <p className="text-sm font-medium text-purple-800">Level 1 Roles</p>
              <p className="text-lg font-semibold text-purple-900">
                {roles.filter(r => r.hierarchy_level === 1).length}
              </p>
            </div>
          </div>
        </div>

        <div className="bg-green-50 rounded-lg p-4">
          <div className="flex items-center">
            <span className="text-green-600 text-lg mr-3">👤</span>
            <div>
              <p className="text-sm font-medium text-green-800">Child Roles</p>
              <p className="text-lg font-semibold text-green-900">
                {roles.filter(r => r.hierarchy_level >= 2).length}
              </p>
            </div>
          </div>
        </div>

        <div className="bg-orange-50 rounded-lg p-4">
          <div className="flex items-center">
            <span className="text-orange-600 text-lg mr-3">📝</span>
            <div>
              <p className="text-sm font-medium text-orange-800">Pending Changes</p>
              <p className="text-lg font-semibold text-orange-900">{pendingMoves.length}</p>
            </div>
          </div>
        </div>
      </div>
    </div>
  );
};

export default RoleHierarchyManager;