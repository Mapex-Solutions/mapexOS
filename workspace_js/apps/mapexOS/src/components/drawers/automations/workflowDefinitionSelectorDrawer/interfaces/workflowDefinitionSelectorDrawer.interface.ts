import type { DefinitionResponse } from '@mapexos/schemas';

/**
 * Props for WorkflowDefinitionSelectorDrawer component.
 * Provides a drawer for selecting a single workflow definition.
 */
export interface WorkflowDefinitionSelectorDrawerProps {
  /** Drawer visibility state */
  modelValue: boolean;

  /** Pre-selected definition ID (optional) */
  selectedDefinitionId?: string;
}

/**
 * Emits for WorkflowDefinitionSelectorDrawer component.
 */
export interface WorkflowDefinitionSelectorDrawerEmits {
  /** Update drawer visibility */
  (e: 'update:modelValue', value: boolean): void;

  /** Emit selected definition */
  (e: 'select', definition: DefinitionResponse): void;

  /** Emit cancel action */
  (e: 'cancel'): void;
}
