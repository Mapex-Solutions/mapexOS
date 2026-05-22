/**
 * JsExecLogsPage Constants
 */

/**
 * Default limit for cursor pagination
 */
export const DEFAULT_LIMIT = 15;

/**
 * Default column visibility state for JS executor logs list
 */
export const COLUMN_VISIBILITY_DEFAULTS = {
  threadId: true,
  name: true,
  success: true,
  executionTime: true,
  created: true,
} as const;

/**
 * Success status color mapping
 */
export const SUCCESS_COLORS = {
  success: 'green-6',
  failed: 'red-6',
} as const;

/**
 * Failed step color mapping
 */
export const FAILED_AT_COLORS = {
  script_validation: 'orange-6',
  script_execution: 'red-6',
  payload_transformation: 'purple-6',
  timeout: 'grey-7',
} as const;

/**
 * Default fallback color
 */
export const DEFAULT_COLOR = 'grey-6';
