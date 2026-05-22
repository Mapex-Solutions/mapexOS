import type { TriggerResponse } from '@mapexos/schemas';

/**
 * Props for TriggerSelectorDialog component
 * Centered modal for selecting a trigger (used in workflow context)
 */
export interface TriggerSelectorDialogProps {
  /** Dialog visibility state */
  modelValue: boolean;

  /** Pre-selected trigger ID for highlighting */
  selectedTriggerId?: string | null;
}

/**
 * Emits for TriggerSelectorDialog component
 */
export interface TriggerSelectorDialogEmits {
  /** Update dialog visibility */
  (e: 'update:modelValue', value: boolean): void;

  /** Emit selected trigger */
  (e: 'select', trigger: TriggerResponse): void;

  /** Emit cancel action */
  (e: 'cancel'): void;
}
