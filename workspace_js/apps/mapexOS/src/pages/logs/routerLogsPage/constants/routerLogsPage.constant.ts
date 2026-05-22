/**
 * Router Logs Page Constants
 */

/** Default number of items per page */
export const DEFAULT_LIMIT = 20;

/** Default column visibility settings */
export const COLUMN_VISIBILITY_DEFAULTS = {
  threadId: true,
  name: true,
  success: true,
  routersCount: true,
  publishedCount: true,
  created: true,
} as const;

/** Colors for success status */
export const SUCCESS_COLORS = {
  success: 'green-6',
  failed: 'red-6',
} as const;

/** Colors for router count badges */
export const COUNT_COLORS = {
  zero: 'grey-6',
  partial: 'orange-6',
  full: 'green-6',
} as const;
