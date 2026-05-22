import type { PermissionState } from './types';

/**
 * Check if a user permission matches a required permission.
 * Replicates the exact backend wildcard matching logic.
 *
 * Match Rules:
 * 1. Root/Admin wildcards match EVERYTHING
 * 2. Exact match
 * 3. Resource wildcard: "assets.*" matches "assets.list"
 * 4. Nested wildcard: "events.*" matches "events.raw.list" (prefix match)
 *
 * @param {string} userPerm - Permission the user has
 * @param {string} requiredPerm - Permission being checked
 * @returns {boolean} True if userPerm satisfies requiredPerm
 */
function matchesPermission(userPerm: string, requiredPerm: string): boolean {
  // 1. Root permission — matches EVERYTHING
  if (userPerm === 'mapex.*') return true;

  // 2. Admin wildcards — match everything within their scope
  if (userPerm === 'admin_vendor.*') return true;
  if (userPerm === 'admin_customer.*') return true;
  if (userPerm === 'admin.*') return true;

  // 3. Exact match
  if (userPerm === requiredPerm) return true;

  // 4. Wildcard match: "resource.*" matches "resource.action" and "resource.sub.action"
  if (userPerm.endsWith('.*')) {
    const resource = userPerm.slice(0, -2);
    if (requiredPerm.startsWith(resource + '.')) return true;
  }

  return false;
}

export const getters = {
  /**
   * Check if user has a specific permission (supports wildcards)
   */
  hasPermission: (state: PermissionState) => (required: string): boolean => {
    return state.permissions.some(userPerm => matchesPermission(userPerm, required));
  },

  /**
   * Check if user has ANY of the required permissions
   */
  hasAnyPermission: (state: PermissionState) => (required: string[]): boolean => {
    if (!required.length) return true;
    return required.some(req =>
      state.permissions.some(userPerm => matchesPermission(userPerm, req)),
    );
  },

  /**
   * Check if user has ALL of the required permissions
   */
  hasAllPermissions: (state: PermissionState) => (required: string[]): boolean => {
    if (!required.length) return true;
    return required.every(req =>
      state.permissions.some(userPerm => matchesPermission(userPerm, req)),
    );
  },

  /**
   * Whether permissions have been loaded and are not currently loading
   */
  isLoaded: (state: PermissionState): boolean => {
    return state.permissions.length > 0 && !state.loading;
  },

  /**
   * Whether cached permissions are stale (> 5 minutes old)
   */
  isStale: (state: PermissionState): boolean => {
    if (!state.lastFetched) return true;
    const fiveMinutesAgo = Date.now() - 5 * 60 * 1000;
    return new Date(state.lastFetched).getTime() < fiveMinutesAgo;
  },
};
