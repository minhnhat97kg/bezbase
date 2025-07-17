// PermissionService: check if user has permission for resource/action using AuthContext

/**
 * Checks if the user has a specific permission.
 * @param {Array} permissions - Array of permission objects from AuthContext
 * @param {string} resource - Resource name to check
 * @param {string} action - Action name to check
 * @returns {boolean}
 */
export function hasPermission(permissions, resource, action) {
  if (!Array.isArray(permissions)) return false;
  return permissions.some(p => {
    // permissions format like role:resource:action then have to split it
    const [role, res, act] = p.split(':');
    p = { role, resource: res, action: act };

    // Support wildcard '*' for resource or action
    const resourceMatch = p.resource === resource || p.resource === '*';
    const actionMatch = p.action === action || p.action === '*';
    return resourceMatch && actionMatch;
  });
}

/**
 * HOC to hide children if user lacks permission
 * @param {ReactNode} children
 * @param {Array} permissions
 * @param {string} resource
 * @param {string} action
 */
export function PermissionGuard({ permissions, resource, action, children }) {
  if (!hasPermission(permissions, resource, action)) return null;
  return children;
}
