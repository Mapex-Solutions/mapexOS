/**
 * CustomersListPage Constants
 */

/**
 * Default number of items per page
 */
export const DEFAULT_ITEMS_PER_PAGE = 15;

/**
 * Default column visibility state for customer list
 */
export const COLUMN_VISIBILITY_DEFAULTS = {
  organization: true,
  address: true,
  created: true,
} as const;

/**
 * Default filter values for customer list
 */
export const FILTER_DEFAULTS = {
  name: undefined,
  enabled: undefined,
  includeChildren: undefined,
  type: undefined,
} as const;

/**
 * API projection fields for customer list query
 */
export const CUSTOMERS_PROJECTION = 'name,type,enabled,address,phone,logo,created,customerId' as const;
