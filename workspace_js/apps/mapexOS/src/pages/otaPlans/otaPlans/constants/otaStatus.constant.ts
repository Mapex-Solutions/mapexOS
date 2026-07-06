import type { DetailChipColor } from '@components/chips';

/**
 * Shared OTA status visuals (colors + icons), keyed by the exact wire values
 * the platform returns. Used by the plans list (DataRow, Quasar palette) and
 * the plan detail page (DetailChip, design-token color). One visual language
 * across both screens.
 */

/**
 * Plan status -> Quasar palette color (DataRow chips / avatars). Mirrors the
 * template-migrations color language (scheduled=blue, running=primary,
 * complete=positive, closed≈completed_with_errors=orange, canceled/aborted=red).
 */
export const OTA_PLAN_STATUS_COLORS: Record<string, string> = {
	SCHEDULED: 'blue-6',
	IN_PROGRESS: 'primary',
	COMPLETED: 'positive',
	CLOSED: 'orange-6',
	CANCELED: 'negative',
	ABORTED: 'negative',
};

/** Plan status -> DetailChip color token (same language as above) */
export const OTA_PLAN_STATUS_CHIP_COLORS: Record<string, DetailChipColor> = {
	SCHEDULED: 'blue',
	IN_PROGRESS: 'primary',
	COMPLETED: 'positive',
	CLOSED: 'orange',
	CANCELED: 'negative',
	ABORTED: 'negative',
};

/** Plan status -> Material icon (mirrors the migration icon set) */
export const OTA_PLAN_STATUS_ICONS: Record<string, string> = {
	SCHEDULED: 'event',
	IN_PROGRESS: 'sync',
	COMPLETED: 'check_circle',
	CLOSED: 'flag',
	CANCELED: 'cancel',
	ABORTED: 'report',
};

/** Execution state -> DetailChip color token */
export const OTA_EXECUTION_STATE_COLORS: Record<string, DetailChipColor> = {
	QUEUED: 'grey',
	INITIATED: 'blue',
	DOWNLOADING: 'orange',
	DOWNLOADED: 'amber',
	VERIFIED: 'cyan',
	UPDATING: 'purple',
	UPDATED: 'positive',
	FAILED: 'negative',
	TIMED_OUT: 'blue-grey',
};

/** Execution state -> Material icon */
export const OTA_EXECUTION_STATE_ICONS: Record<string, string> = {
	QUEUED: 'hourglass_empty',
	INITIATED: 'send',
	DOWNLOADING: 'downloading',
	DOWNLOADED: 'download_done',
	VERIFIED: 'verified',
	UPDATING: 'system_update',
	UPDATED: 'check_circle',
	FAILED: 'error',
	TIMED_OUT: 'timer_off',
};

/** Plan statuses that still allow cancel (non-terminal) */
export const OTA_PLAN_CANCELABLE_STATUSES = ['SCHEDULED', 'IN_PROGRESS'];
