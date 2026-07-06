/**
 * Props for StepSchedule: the scheduleAt slice of the form model.
 * `null` means "run now"; a non-null value is an ISO datetime string.
 */
export interface StepScheduleProps {
  modelValue: string | null;
}

/**
 * Emits for StepSchedule.
 * - update:modelValue emits an ISO datetime string, or null for run-now.
 * - update:valid signals whether the current choice is acceptable.
 */
export interface StepScheduleEmits {
  (e: 'update:modelValue', value: string | null): void;
  (e: 'update:valid', value: boolean): void;
}

/** Execution mode for the schedule step. */
export type ScheduleMode = 'now' | 'scheduled';
