/**
 * AccessAuditPage Constants
 */

/**
 * Default number of items per page
 */
export const DEFAULT_ITEMS_PER_PAGE = 15;

/**
 * Default column visibility state for access audit list
 */
export const COLUMN_VISIBILITY_DEFAULTS = {
  organization: true,
  roles: true,
  scope: true,
  enabled: true,
} as const;

/**
 * Default filter values for access audit list
 */
export const FILTER_DEFAULTS = {
  assigneeType: undefined,
  assigneeId: undefined,
  roleId: undefined,
  scope: undefined,
  enabled: undefined,
  includeChildren: undefined,
} as const;

/**
 * API projection fields for memberships query
 */
export const MEMBERSHIPS_PROJECTION = 'assigneeType,assigneeId,orgId,orgPathKey,roleIds,scope,enabled,created' as const;

/**
 * Assignee type options
 */
export const ASSIGNEE_TYPE_OPTIONS = ['user', 'group'] as const;

/**
 * Scope type options
 */
export const SCOPE_OPTIONS = ['local', 'recursive'] as const;
