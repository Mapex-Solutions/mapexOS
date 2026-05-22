/** Default number of items per page */
export const DEFAULT_LIMIT = 20;

/** Maximum visible filter chips */
export const MAX_VISIBLE_CHIPS = 2;

/** Status colors for execution badges */
export const STATUS_COLORS: Record<string, string> = {
  running: 'blue-6',
  waiting: 'orange-6',
  completed: 'green-6',
  failed: 'red-6',
  cancelled: 'amber-8',
} as const;

/** Status icons */
export const STATUS_ICONS: Record<string, string> = {
  running: 'play_circle',
  waiting: 'hourglass_top',
  completed: 'check_circle',
  failed: 'error',
  cancelled: 'cancel',
} as const;

/** Duration color thresholds */
export const DURATION_THRESHOLDS = {
  fast: 1000,
  normal: 10000,
} as const;

/** Duration colors */
export const DURATION_COLORS = {
  fast: 'green-6',
  normal: 'orange-6',
  slow: 'red-6',
} as const;

/** Hot statuses (MongoDB via Workflow Service) */
export const HOT_STATUSES = ['running', 'waiting'] as const;

/** Cold statuses (ClickHouse via Events Service) */
export const COLD_STATUSES = ['completed', 'failed', 'cancelled'] as const;
