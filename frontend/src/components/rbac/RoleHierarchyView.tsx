import React, { useState } from 'react';
import { rbacService, Role } from '../../services/rbacService';

interface RoleHierarchyViewProps {
  roles: Role[];
  onRoleSelect?: (role: Role) => void;
  selectedRoleId?: number;
}

interface HierarchyData {
  role_id: number;
  parent_roles: Role[];
  child_roles: Role[];
}

const RoleHierarchyView: React.FC<RoleHierarchyViewProps> = ({
  roles,
  onRoleSelect,
  selectedRoleId
}) => {
  const [hierarchyData, setHierarchyData] = useState<Record<number, HierarchyData>>({});
  const [expandedNodes, setExpandedNodes] = useState<Set<number>>(new Set());
  const [loading, setLoading] = useState(false);

  // Build hierarchy tree from roles
  const buildHierarchyTree = () => {
    const rootRoles = roles.filter(role => !role.parent_role_id);
    
    const buildNode = (role: Role): RoleNode => {
      const children = roles.filter(r => r.parent_role_id === role.id);
      return {
        role,
        children: children.map(child => buildNode(child)),
        level: role.hierarchy_level || 0
      };
    };

    return rootRoles.map(role => buildNode(role));
  };

  interface RoleNode {
    role: Role;
    children: RoleNode[];
    level: number;
  }

  const loadHierarchyData = async (roleId: number) => {
    if (hierarchyData[roleId]) return;
    
    try {
      setLoading(true);
      const response = await rbacService.getRoleHierarchy(roleId);
      setHierarchyData(prev => ({
        ...prev,
        [roleId]: response.data
      }));
    } catch (error) {
      console.error('Failed to load hierarchy data:', error);
    } finally {
      setLoading(false);
    }
  };

  const toggleExpanded = (roleId: number) => {
    const newExpanded = new Set(expandedNodes);
    if (newExpanded.has(roleId)) {
      newExpanded.delete(roleId);
    } else {
      newExpanded.add(roleId);
      loadHierarchyData(roleId);
    }
    setExpandedNodes(newExpanded);
  };

  const renderRoleNode = (node: RoleNode, depth: number = 0) => {
    const isExpanded = expandedNodes.has(node.role.id);
    const hasChildren = node.children.length > 0;
    const isSelected = selectedRoleId === node.role.id;

    return (
      <div key={node.role.id} className="role-hierarchy-node">
        <div 
          className={`flex items-center p-2 rounded-lg cursor-pointer transition-colors ${
            isSelected 
              ? 'bg-indigo-100 border-indigo-300 border' 
              : 'hover:bg-gray-50'
          }`}
          style={{ marginLeft: `${depth * 24}px` }}
          onClick={() => onRoleSelect?.(node.role)}
        >
          {/* Expand/Collapse Button */}
          {hasChildren && (
            <button
              onClick={(e) => {
                e.stopPropagation();
                toggleExpanded(node.role.id);
              }}
              className="mr-2 p-1 hover:bg-gray-200 rounded"
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

          {/* Role Icon based on hierarchy level */}
          <span className="mr-3 text-lg">
            {node.level === 0 ? '👑' : node.level === 1 ? '🎭' : '👤'}
          </span>

          {/* Role Info */}
          <div className="flex-1">
            <div className="flex items-center space-x-2">
              <span className="font-medium text-gray-900">
                {node.role.display_name || node.role.name}
              </span>
              
              {/* Hierarchy Level Badge */}
              <span className="inline-flex items-center px-2 py-0.5 rounded-full text-xs font-medium bg-gray-100 text-gray-800">
                Level {node.level}
              </span>

              {/* System Role Badge */}
              {node.role.is_system && (
                <span className="inline-flex items-center px-2 py-0.5 rounded-full text-xs font-medium bg-red-100 text-red-800">
                  System
                </span>
              )}

              {/* Parent Role Indicator */}
              {node.role.parent_role_id && (
                <span className="inline-flex items-center px-2 py-0.5 rounded-full text-xs font-medium bg-blue-100 text-blue-800">
                  ↗ Inherits
                </span>
              )}
            </div>
            
            {node.role.description && (
              <p className="text-sm text-gray-500 mt-1">{node.role.description}</p>
            )}
          </div>

          {/* Child Count */}
          {hasChildren && (
            <span className="ml-2 text-sm text-gray-500">
              {node.children.length} child{node.children.length !== 1 ? 'ren' : ''}
            </span>
          )}
        </div>

        {/* Render Children */}
        {isExpanded && hasChildren && (
          <div className="ml-6 border-l border-gray-200 pl-2">
            {node.children.map(child => renderRoleNode(child, depth + 1))}
          </div>
        )}
      </div>
    );
  };

  const hierarchyTree = buildHierarchyTree();

  return (
    <div className="role-hierarchy-view">
      <div className="bg-white shadow rounded-lg">
        <div className="px-4 py-5 sm:px-6 border-b border-gray-200">
          <h3 className="text-lg font-medium text-gray-900">Role Hierarchy</h3>
          <p className="mt-1 text-sm text-gray-500">
            Visual representation of role inheritance relationships
          </p>
        </div>

        <div className="p-6">
          {hierarchyTree.length === 0 ? (
            <div className="text-center py-8">
              <span className="text-gray-400 text-4xl mb-4 block">🎭</span>
              <p className="text-gray-500">No roles found</p>
            </div>
          ) : (
            <div className="space-y-2">
              {hierarchyTree.map(node => renderRoleNode(node))}
            </div>
          )}

          {loading && (
            <div className="flex justify-center py-4">
              <div className="animate-spin rounded-full h-6 w-6 border-b-2 border-indigo-500"></div>
            </div>
          )}
        </div>
      </div>

      {/* Legend */}
      <div className="mt-4 bg-gray-50 rounded-lg p-4">
        <h4 className="text-sm font-medium text-gray-700 mb-2">Legend</h4>
        <div className="grid grid-cols-2 gap-2 text-xs text-gray-600">
          <div className="flex items-center">
            <span className="mr-2">👑</span>
            <span>Root Role (Level 0)</span>
          </div>
          <div className="flex items-center">
            <span className="mr-2">🎭</span>
            <span>Parent Role (Level 1)</span>
          </div>
          <div className="flex items-center">
            <span className="mr-2">👤</span>
            <span>Child Role (Level 2+)</span>
          </div>
          <div className="flex items-center">
            <span className="mr-2 inline-flex items-center px-2 py-0.5 rounded-full text-xs font-medium bg-blue-100 text-blue-800">↗ Inherits</span>
            <span>Has Parent Role</span>
          </div>
        </div>
      </div>
    </div>
  );
};

export default RoleHierarchyView;