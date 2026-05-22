/**
 * RolesListPage Constants
 */

/**
 * Default number of items per page
 */
export const DEFAULT_ITEMS_PER_PAGE = 15;

/**
 * Default column visibility state for roles list
 */
export const COLUMN_VISIBILITY_DEFAULTS = {
  organization: true,
  description: true,
  permissions: true,
  scope: true,
  isTemplate: true,
  created: true,
} as const;

/**
 * Default filter values for roles list
 */
export const FILTER_DEFAULTS = {
  name: undefined,
  isSystem: undefined,
  scope: undefined,
  permission: undefined,
  includeChildren: undefined,
  isTemplate: undefined,
} as const;

/**
 * API projection fields for roles query
 */
export const ROLES_PROJECTION = 'name,description,permissions,isSystem,scope,isTemplate,orgId,created' as const;
