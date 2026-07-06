import type { OTARolloutConfig } from '@mapexos/schemas';

/**
 * Props for StepSchedule: the scheduling slice of the form model.
 * `startAt` null means "run immediately"; a non-null value is an ISO datetime.
 * `maxTime` is the mandatory rollout deadline (ISO datetime).
 */
export interface StepScheduleProps {
	startAt: string | null;
	maxTime: string | null;
	rolloutConfig: OTARolloutConfig;
}

/**
 * Emits for StepSchedule.
 * - update:startAt emits an ISO datetime string, or null for run-now.
 * - update:maxTime emits the deadline ISO datetime string, or null while invalid.
 * - update:rolloutConfig emits the pacing/abort config.
 * - update:valid signals whether the whole schedule slice is acceptable.
 */
export interface StepScheduleEmits {
	(e: 'update:startAt', value: string | null): void;
	(e: 'update:maxTime', value: string | null): void;
	(e: 'update:rolloutConfig', value: OTARolloutConfig): void;
	(e: 'update:valid', value: boolean): void;
}

/** Execution mode for the schedule step. */
export type ScheduleMode = 'now' | 'scheduled';
