import React, { useState, useEffect, useCallback } from 'react';
import { userService } from '../../services/api';
import { rbacService } from '../../services/rbacService';
import Icon from '../common/Icons';

const UserRoleAssignment = ({ roles, onRefresh }) => {
  const [selectedUser, setSelectedUser] = useState(null);
  const [userRoles, setUserRoles] = useState([]);
  const [availableRoles, setAvailableRoles] = useState([]);
  const [loading, setLoading] = useState(false);
  const [error, setError] = useState(null);
  const [showAssignForm, setShowAssignForm] = useState(false);
  const [roleToAssign, setRoleToAssign] = useState('');

  // User search state
  const [userSearchTerm, setUserSearchTerm] = useState('');
  const [searchResults, setSearchResults] = useState([]);
  const [showUserDropdown, setShowUserDropdown] = useState(false);
  const [searchLoading, setSearchLoading] = useState(false);

  // Debounced search for users
  const searchUsers = useCallback(async (searchTerm) => {
    if (!searchTerm.trim()) {
      setSearchResults([]);
      return;
    }

    try {
      setSearchLoading(true);
      const response = await userService.getUsers(searchTerm.trim());
      const userData = response.data || [];
      setSearchResults(userData);
    } catch (err) {
      console.error('Failed to search users:', err);
      setSearchResults([]);
    } finally {
      setSearchLoading(false);
    }
  }, []);

  // Debounce search
  useEffect(() => {
    const timeoutId = setTimeout(() => {
      if (userSearchTerm) {
        searchUsers(userSearchTerm);
      } else {
        setSearchResults([]);
      }
    }, 300);

    return () => clearTimeout(timeoutId);
  }, [userSearchTerm, searchUsers]);

  const fetchUserRoles = useCallback(async (userId) => {
    if (!userId) return;
    
    try {
      setLoading(true);
      setError(null);
      const response = await rbacService.getUserRoles(userId);
      setUserRoles(response.data || []);
      
      // Filter out roles that user already has
      const assignedRoleNames = response.data || [];
      const available = roles.filter(role => !assignedRoleNames.includes(role.name));
      setAvailableRoles(available);
    } catch (err) {
      setError(err.response?.data?.message || 'Failed to fetch user roles');
      setUserRoles([]);
      setAvailableRoles(roles);
    } finally {
      setLoading(false);
    }
  }, [roles]);

  useEffect(() => {
    if (selectedUser) {
      fetchUserRoles(selectedUser.id);
    }
  }, [selectedUser, fetchUserRoles]);

  const handleUserSelect = (user) => {
    setSelectedUser(user);
    setUserSearchTerm(`${user.first_name} ${user.last_name} (${user.email})`);
    setShowUserDropdown(false);
    setUserRoles([]);
    setError(null);
  };

  const clearUserSelection = () => {
    setSelectedUser(null);
    setUserSearchTerm('');
    setUserRoles([]);
    setAvailableRoles(roles);
    setShowUserDropdown(false);
    setError(null);
  };

  const handleSearchInputChange = (e) => {
    const value = e.target.value;
    setUserSearchTerm(value);
    setShowUserDropdown(true);
    
    // If user clears the input, reset selection
    if (!value.trim()) {
      clearUserSelection();
    }
  };

  const handleAssignRole = async (e) => {
    e.preventDefault();
    if (!selectedUser || !roleToAssign) {
      return;
    }

    try {
      setLoading(true);
      await rbacService.assignRole({
        user_id: selectedUser.id,
        role: roleToAssign,
      });
      setRoleToAssign('');
      setShowAssignForm(false);
      setShowUserDropdown(false); // Ensure dropdown is closed
      setError(null); // Clear any previous errors
      // Do NOT clear user selection or search input
      // Don't call onRefresh() to prevent parent re-render
      await fetchUserRoles(selectedUser.id);
    } catch (err) {
      setError(err.response?.data?.message || 'Failed to assign role');
    } finally {
      setLoading(false);
    }

  };

  const handleRemoveRole = async (roleName) => {
    if (!window.confirm(`Remove role "${roleName}" from user?`)) {
      return;
    }

    try {
      setLoading(true);
      await rbacService.removeRole({
        user_id: selectedUser.id,
        role: roleName,
      });
      setShowUserDropdown(false); // Ensure dropdown is closed
      setError(null); // Clear any previous errors
      // Do NOT clear user selection or search input
      // Don't call onRefresh() to prevent parent re-render
      await fetchUserRoles(selectedUser.id);
    } catch (err) {
      setError(err.response?.data?.message || 'Failed to remove role');
    } finally {
      setLoading(false);
    }
  };

  const getRoleByName = (roleName) => {
    return roles.find(role => role.name === roleName);
  };

  const getRoleBadge = (roleName) => {
    const role = getRoleByName(roleName);
    if (!role) return null;
    
    return (
      <span className={`inline-flex items-center px-2.5 py-0.5 rounded-full text-xs font-medium ${
        role.is_system 
          ? 'bg-gray-100 text-gray-800 dark:bg-gray-800 dark:text-gray-100' 
          : 'bg-blue-100 text-blue-800 dark:bg-blue-800 dark:text-blue-100'
      }`}>
        {role.display_name}
        {role.is_system && (
          <span className="ml-1 text-xs">
            (System)
          </span>
        )}
      </span>
    );
  };

  const formatUserName = (user) => {
    const fullName = `${user.first_name} ${user.last_name}`.trim();
    return fullName || 'Unknown User';
  };

  return (
    <div className="space-y-6">
      {error && (
        <div className="bg-red-50 border border-red-200 text-red-700 px-4 py-3 rounded-md">
          {error}
        </div>
      )}

      {/* User Search */}
      <div className="relative">
        <label htmlFor="user-search" className="block text-sm font-medium text-gray-700 dark:text-gray-300 mb-2">
          Search for User
        </label>
        <div className="relative">
          <input
            id="user-search"
            type="text"
            value={userSearchTerm}
            onChange={handleSearchInputChange}
            onFocus={() => setShowUserDropdown(true)}
            placeholder="Type to search users by name or email..."
            className="block w-full px-3 py-2 pr-10 border border-gray-300 dark:border-gray-600 rounded-md shadow-sm focus:outline-none focus:ring-blue-500 focus:border-blue-500 dark:bg-gray-700 dark:text-white"
          />
          
          {/* Search icon or clear button */}
          <div className="absolute inset-y-0 right-0 flex items-center pr-3">
            {selectedUser ? (
              <button
                onClick={clearUserSelection}
                className="text-gray-400 hover:text-gray-600 dark:hover:text-gray-300"
              >
                <Icon name="close" />
              </button>
            ) : (
              <Icon name="search" className="text-gray-400" />
            )}
          </div>
        </div>

        {/* Search Results Dropdown */}
        {showUserDropdown && (userSearchTerm || searchResults.length > 0) && (
          <div className="absolute z-10 w-full mt-1 bg-white dark:bg-gray-800 border border-gray-300 dark:border-gray-600 rounded-md shadow-lg max-h-60 overflow-auto">
            {searchLoading ? (
              <div className="flex justify-center py-3">
                <div className="animate-spin rounded-full h-4 w-4 border-b-2 border-blue-600"></div>
              </div>
            ) : searchResults.length > 0 ? (
              searchResults.map((user) => (
                <button
                  key={user.id}
                  onClick={() => handleUserSelect(user)}
                  className="w-full text-left px-4 py-2 hover:bg-gray-100 dark:hover:bg-gray-700 focus:outline-none focus:bg-gray-100 dark:focus:bg-gray-700"
                >
                  <div className="flex justify-between items-center">
                    <div>
                      <div className="text-sm font-medium text-gray-900 dark:text-white">
                        {formatUserName(user)}
                      </div>
                      <div className="text-sm text-gray-500 dark:text-gray-400">
                        {user.email}
                      </div>
                    </div>
                    <div className="text-xs text-gray-400 dark:text-gray-500">
                      ID: {user.id}
                    </div>
                  </div>
                </button>
              ))
            ) : userSearchTerm && !searchLoading ? (
              <div className="px-4 py-3 text-sm text-gray-500 dark:text-gray-400">
                No users found matching "{userSearchTerm}"
              </div>
            ) : null}
          </div>
        )}
      </div>

      {selectedUser && (
        <div>
          {/* Selected User Info */}
          <div className="bg-blue-50 dark:bg-blue-900 border border-blue-200 dark:border-blue-700 rounded-md p-4 mb-6">
            <div className="flex items-center justify-between">
              <div>
                <h3 className="text-lg font-medium text-blue-900 dark:text-blue-100">
                  {formatUserName(selectedUser)}
                </h3>
                <p className="text-sm text-blue-700 dark:text-blue-300">
                  {selectedUser.email} • ID: {selectedUser.id}
                </p>
                <p className="text-xs text-blue-600 dark:text-blue-400 mt-1">
                  Status: {selectedUser.status} • Email Verified: {selectedUser.email_verified ? 'Yes' : 'No'}
                </p>
              </div>
              {availableRoles.length > 0 && (
                <button
                  onClick={() => setShowAssignForm(true)}
                  className="bg-blue-600 hover:bg-blue-700 text-white px-4 py-2 rounded-md text-sm font-medium"
                >
                  Assign Role
                </button>
              )}
            </div>
          </div>

          {/* Current User Roles */}
          <div className="mb-6">
            <h4 className="text-md font-medium text-gray-900 dark:text-white mb-3">
              Current Roles ({userRoles.length})
            </h4>
            
            {loading && !showAssignForm ? (
              <div className="flex justify-center py-4">
                <div className="animate-spin rounded-full h-6 w-6 border-b-2 border-blue-600"></div>
              </div>
            ) : userRoles.length > 0 ? (
              <div className="grid grid-cols-1 sm:grid-cols-2 lg:grid-cols-3 gap-3">
                {userRoles.map((roleName, index) => {
                  const role = getRoleByName(roleName);
                  return (
                    <div key={index} className="flex items-center justify-between p-3 bg-gray-50 dark:bg-gray-700 rounded-md border">
                      <div className="flex flex-col">
                        {getRoleBadge(roleName)}
                        <span className="text-xs text-gray-600 dark:text-gray-400 mt-1">
                          {role?.description || 'No description'}
                        </span>
                      </div>
                      <button
                        onClick={() => handleRemoveRole(roleName)}
                        className="text-red-600 hover:text-red-800 dark:text-red-400 dark:hover:text-red-300 ml-2"
                        disabled={loading}
                        title="Remove role"
                      >
                        <Icon name="delete" />
                      </button>
                    </div>
                  );
                })}
              </div>
            ) : (
              <div className="text-center py-6 text-gray-500 dark:text-gray-400 bg-gray-50 dark:bg-gray-700 rounded-md">
                No roles assigned to this user
              </div>
            )}
          </div>

          {/* Assign Role Form */}
          {showAssignForm && (
            <div className="mb-6 p-4 border border-gray-200 dark:border-gray-600 rounded-md bg-gray-50 dark:bg-gray-700">
              <h4 className="text-sm font-medium text-gray-900 dark:text-white mb-3">Assign New Role</h4>
              <form onSubmit={handleAssignRole} className="space-y-4">
                <div>
                  <label className="block text-sm font-medium text-gray-700 dark:text-gray-300 mb-1">
                    Select Role
                  </label>
                  <select
                    value={roleToAssign}
                    onChange={(e) => setRoleToAssign(e.target.value)}
                    className="block w-full px-3 py-2 border border-gray-300 dark:border-gray-600 rounded-md shadow-sm focus:outline-none focus:ring-blue-500 focus:border-blue-500 dark:bg-gray-600 dark:text-white"
                    required
                  >
                    <option value="">Choose a role...</option>
                    {availableRoles.map((role) => (
                      <option key={role.id} value={role.name}>
                        {role.display_name} ({role.name})
                      </option>
                    ))}
                  </select>
                </div>
                <div className="flex justify-end space-x-3">
                  <button
                    type="button"
                    onClick={() => setShowAssignForm(false)}
                    className="px-4 py-2 text-sm font-medium text-gray-700 bg-white border border-gray-300 rounded-md hover:bg-gray-50 focus:outline-none focus:ring-2 focus:ring-offset-2 focus:ring-blue-500 dark:bg-gray-600 dark:text-gray-300 dark:border-gray-500 dark:hover:bg-gray-500"
                  >
                    Cancel
                  </button>
                  <button
                    type="submit"
                    disabled={loading}
                    className="px-4 py-2 text-sm font-medium text-white bg-blue-600 border border-transparent rounded-md hover:bg-blue-700 focus:outline-none focus:ring-2 focus:ring-offset-2 focus:ring-blue-500 disabled:opacity-50"
                  >
                    {loading ? 'Assigning...' : 'Assign Role'}
                  </button>
                </div>
              </form>
            </div>
          )}

          {/* Permission Check Tool */}
          <div className="mt-8 p-4 border-t border-gray-200 dark:border-gray-600">
            <h4 className="text-sm font-medium text-gray-900 dark:text-white mb-3">
              Check User Permissions
            </h4>
            <PermissionChecker userId={selectedUser.id} />
          </div>
        </div>
      )}
    </div>
  );
};

// Helper component for checking permissions
const PermissionChecker = ({ userId }) => {
  const [resource, setResource] = useState('');
  const [action, setAction] = useState('');
  const [result, setResult] = useState(null);
  const [loading, setLoading] = useState(false);

  const checkPermission = async (e) => {
    e.preventDefault();
    if (!userId || !resource || !action) return;

    try {
      setLoading(true);
      const response = await rbacService.checkUserPermission(userId, resource, action);
      setResult(response.data.allowed);
    } catch (err) {
      setResult(false);
    } finally {
      setLoading(false);
    }
  };

  return (
    <form onSubmit={checkPermission} className="space-y-3">
      <div className="grid grid-cols-3 gap-3">
        <input
          type="text"
          value={resource}
          onChange={(e) => setResource(e.target.value)}
          placeholder="Resource (e.g., users)"
          className="px-3 py-2 border border-gray-300 dark:border-gray-600 rounded-md shadow-sm focus:outline-none focus:ring-blue-500 focus:border-blue-500 dark:bg-gray-600 dark:text-white"
        />
        <input
          type="text"
          value={action}
          onChange={(e) => setAction(e.target.value)}
          placeholder="Action (e.g., read)"
          className="px-3 py-2 border border-gray-300 dark:border-gray-600 rounded-md shadow-sm focus:outline-none focus:ring-blue-500 focus:border-blue-500 dark:bg-gray-600 dark:text-white"
        />
        <button
          type="submit"
          disabled={loading}
          className="px-4 py-2 text-sm font-medium text-white bg-blue-600 border border-transparent rounded-md hover:bg-blue-700 focus:outline-none focus:ring-2 focus:ring-offset-2 focus:ring-blue-500 disabled:opacity-50"
        >
          {loading ? 'Checking...' : 'Check'}
        </button>
      </div>
      {result !== null && (
        <div className={`text-sm px-3 py-2 rounded-md ${
          result 
            ? 'bg-green-100 text-green-800 dark:bg-green-800 dark:text-green-100' 
            : 'bg-red-100 text-red-800 dark:bg-red-800 dark:text-red-100'
        }`}>
          Permission {result ? 'ALLOWED' : 'DENIED'}
        </div>
      )}
    </form>
  );
};

export default UserRoleAssignment;