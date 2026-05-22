/**
 * Props for TriggerDetailDrawer component
 */
export interface TriggerDetailDrawerProps {
  /** Controls drawer visibility */
  modelValue: boolean;
  /** ID of the trigger to display */
  triggerId: string | null;
}

/**
 * Emits for TriggerDetailDrawer component
 */
export interface TriggerDetailDrawerEmits {
  (e: 'update:modelValue', value: boolean): void;
  (e: 'edit', triggerId: string): void;
}
