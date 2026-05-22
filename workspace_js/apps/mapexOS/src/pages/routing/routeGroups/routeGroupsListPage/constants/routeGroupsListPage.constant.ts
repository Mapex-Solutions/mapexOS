/**
 * RouteGroupsListPage Constants
 */

/**
 * Default number of items per page
 */
export const DEFAULT_ITEMS_PER_PAGE = 15;

/**
 * Default column visibility state for route groups list
 */
export const COLUMN_VISIBILITY_DEFAULTS = {
  organization: true,
  routers: true,
  isTemplate: true,
} as const;

/**
 * Default filter values for route groups list
 */
export const FILTER_DEFAULTS = {
  name: undefined,
  enabled: undefined,
  isTemplate: undefined,
  includeChildren: undefined,
} as const;

/**
 * API projection fields for route groups query
 */
export const ROUTE_GROUPS_PROJECTION = 'name,description,enabled,routers,isTemplate,orgId' as const;
