/**
 * AssetRawLogsPage Constants
 */

/**
 * Default limit for cursor pagination
 */
export const DEFAULT_LIMIT = 15;

/**
 * Default column visibility state for raw logs list
 */
export const COLUMN_VISIBILITY_DEFAULTS = {
  name: true,
  threadId: true,
  source: true,
  created: true,
} as const;

/**
 * Source color mapping (gateway types)
 */
export const SOURCE_COLORS: Record<string, string> = {
  http_gateway: 'blue-6',
  mqtt_gateway: 'purple-6',
  lorawan_gateway: 'teal-6',
  tcp_gateway: 'indigo-6',
  udp_gateway: 'orange-6',
} as const;

/**
 * Success status color mapping
 */
export const SUCCESS_COLORS = {
  success: 'green-6',
  failed: 'red-6',
} as const;

/**
 * Default fallback color
 */
export const DEFAULT_COLOR = 'grey-6';
