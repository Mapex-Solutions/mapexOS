/**
 * MigrationPlansListPage Constants
 */

import type { MigrationPlanStatus } from '@mapexos/schemas';

/**
 * Default number of items per page
 */
export const DEFAULT_ITEMS_PER_PAGE = 15;

/**
 * Default filter values for the migration plans list
 */
export const MIGRATION_FILTER_DEFAULTS = {
  name: undefined,
  status: undefined,
} as const;

/**
 * Plan statuses that still allow removal/cancellation (mirrors the backend 409 guard).
 */
export const MIGRATION_EDITABLE_STATUSES: MigrationPlanStatus[] = ['pending', 'scheduled'];

/**
 * Ordered list of migration lifecycle statuses used to build the status filter
 */
export const MIGRATION_STATUS_VALUES: MigrationPlanStatus[] = [
  'pending',
  'scheduled',
  'running',
  'complete',
  'completed_with_errors',
  'cancelled',
];

/**
 * Chip color per lifecycle status (Quasar theme-aware color names)
 */
export const MIGRATION_STATUS_COLORS: Record<MigrationPlanStatus, string> = {
  pending: 'grey-6',
  scheduled: 'blue-6',
  running: 'primary',
  complete: 'positive',
  completed_with_errors: 'orange-6',
  cancelled: 'negative',
};

/**
 * Chip icon per lifecycle status
 */
export const MIGRATION_STATUS_ICONS: Record<MigrationPlanStatus, string> = {
  pending: 'schedule',
  scheduled: 'event',
  running: 'sync',
  complete: 'check_circle',
  completed_with_errors: 'error',
  cancelled: 'cancel',
};

/**
 * Date/time mask for the schedule column
 */
export const MIGRATION_SCHEDULE_MASK = 'MMM DD, YYYY HH:mm';

/**
 * Base route for migration plan navigation
 */
export const MIGRATION_ROUTE_BASE = '/assets_template/migrations';
