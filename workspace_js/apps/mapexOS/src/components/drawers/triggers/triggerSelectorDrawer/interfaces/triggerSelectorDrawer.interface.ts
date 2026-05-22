import type { TriggerResponse } from '@mapexos/schemas';

/**
 * Props for TriggerSelectorDrawer component
 */
export interface TriggerSelectorDrawerProps {
  /** Whether the dialog is open */
  modelValue: boolean;

  /** Pre-selected trigger ID for highlighting */
  selectedTriggerId?: string | null;
}

/**
 * Emits for TriggerSelectorDrawer component
 */
export interface TriggerSelectorDrawerEmits {
  'update:modelValue': [value: boolean];
  'select': [trigger: TriggerResponse];
  'cancel': [];
}
