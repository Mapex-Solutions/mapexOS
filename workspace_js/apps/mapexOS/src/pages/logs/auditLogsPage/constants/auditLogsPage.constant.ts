/**
 * AuditLogsPage Constants
 */

import type { AuditLogsPageFilters } from '../interfaces';

/**
 * Default number of items per page
 */
export const DEFAULT_ITEMS_PER_PAGE = 15;

/**
 * Initial page number
 */
export const INITIAL_PAGE = 1;

/**
 * Maximum number of visible filter chips
 */
export const MAX_VISIBLE_CHIPS = 2;

/**
 * Default column visibility state for audit logs list
 */
export const COLUMN_VISIBILITY_DEFAULTS = {
  action: true,
  resource: true,
  created: true,
} as const;

/**
 * Default filter values
 */
export const FILTER_DEFAULTS = {} as AuditLogsPageFilters;

/**
 * Action color mapping
 */
export const ACTION_COLORS: Record<string, string> = {
  Create: 'green-6',
  Edit: 'blue-6',
  Delete: 'red-6',
} as const;

/**
 * Type icon mapping
 */
export const TYPE_ICONS: Record<string, string> = {
  userLog: 'person',
  dataSource: 'storage',
  assets: 'devices',
  payloadHandler: 'code',
  triggers: 'flash_on',
  users: 'group',
  customers: 'business',
} as const;

/**
 * Default fallback icon
 */
export const DEFAULT_TYPE_ICON = 'folder';

/**
 * Default fallback color
 */
export const DEFAULT_ACTION_COLOR = 'grey-6';
