/**
 * AssetConnectivityLogsPage Constants
 */

/**
 * Default limit for cursor pagination
 */
export const DEFAULT_LIMIT = 15;

/**
 * Default column visibility state for connectivity logs list
 */
export const COLUMN_VISIBILITY_DEFAULTS = {
  asset: true,
  assetUUID: true,
  eventType: true,
  lastSeenAt: true,
  missCount: true,
  thresholdMinutes: false,
  created: true,
} as const;

/**
 * Event type color mapping (offline = negative, online = positive).
 */
export const EVENT_TYPE_COLORS: Record<string, string> = {
  offline: 'red-6',
  online: 'green-6',
} as const;

/**
 * Default fallback color
 */
export const DEFAULT_COLOR = 'grey-6';
