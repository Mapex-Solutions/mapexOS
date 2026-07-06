/**
 * MigrationPlanDetailPage Constants
 */

import type { MigrationExecutionStatus, MigrationPlanStatus } from '@mapexos/schemas';
import type { DetailChipColor } from '@components/chips';

/**
 * Default number of executions per page
 */
export const DEFAULT_EXECUTIONS_PER_PAGE = 15;

/**
 * Plan statuses during which edit and cancel actions are allowed.
 * Mirrors the backend guard (409 on any other status).
 */
export const MIGRATION_EDITABLE_STATUSES: MigrationPlanStatus[] = ['pending', 'scheduled'];

/**
 * Ordered list of execution statuses used to build the executions status filter
 */
export const EXECUTION_STATUS_VALUES: MigrationExecutionStatus[] = ['pending', 'migrated', 'failed', 'cancelled'];

/**
 * DetailChip color per execution status (restricted DetailChip palette)
 */
export const EXECUTION_STATUS_COLORS: Record<MigrationExecutionStatus, DetailChipColor> = {
  pending: 'grey',
  migrated: 'positive',
  failed: 'negative',
  cancelled: 'negative',
};

/**
 * DetailChip color per plan lifecycle status (restricted DetailChip palette).
 * Owned here so the detail page does not depend on the list page's Quasar-shade map.
 */
export const PLAN_STATUS_DETAIL_COLORS: Record<MigrationPlanStatus, DetailChipColor> = {
  pending: 'grey',
  scheduled: 'blue',
  running: 'primary',
  complete: 'positive',
  completed_with_errors: 'orange',
  cancelled: 'negative',
};

/**
 * Chip icon per execution status
 */
export const EXECUTION_STATUS_ICONS: Record<MigrationExecutionStatus, string> = {
  pending: 'schedule',
  migrated: 'check_circle',
  failed: 'error',
  cancelled: 'cancel',
};

