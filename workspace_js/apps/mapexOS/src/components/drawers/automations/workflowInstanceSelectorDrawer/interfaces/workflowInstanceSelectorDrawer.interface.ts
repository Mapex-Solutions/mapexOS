import type { InstanceResponse } from '@mapexos/schemas';

/**
 * Props for WorkflowInstanceSelectorDrawer component.
 * Provides a drawer for selecting a single workflow instance.
 */
export interface WorkflowInstanceSelectorDrawerProps {
  /** Drawer visibility state */
  modelValue: boolean;

  /** Pre-selected instance ID (optional) */
  selectedInstanceId?: string;
}

/**
 * Emits for WorkflowInstanceSelectorDrawer component.
 */
export interface WorkflowInstanceSelectorDrawerEmits {
  /** Update drawer visibility */
  (e: 'update:modelValue', value: boolean): void;

  /** Emit selected instance (full response for caller to extract fields) */
  (e: 'select', instance: InstanceResponse): void;

  /** Emit cancel action */
  (e: 'cancel'): void;
}
