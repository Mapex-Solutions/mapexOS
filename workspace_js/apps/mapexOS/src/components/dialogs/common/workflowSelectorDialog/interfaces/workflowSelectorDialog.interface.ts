/**
 * Minimal workflow item for the selector dialog.
 * Lightweight representation — no full definition needed.
 */
export interface WorkflowSelectorItem {
  /** Unique workflow ID */
  id: string;

  /** Display name */
  name: string;

  /** Optional description */
  description?: string;

  /** Whether the workflow is currently enabled */
  enabled: boolean;
}

/**
 * Props for WorkflowSelectorDialog component
 */
export interface WorkflowSelectorDialogProps {
  /** Dialog visibility state */
  modelValue: boolean;

  /** Pre-selected workflow ID for highlighting */
  selectedWorkflowId?: string | null;

  /** Workflow ID to exclude from the list (prevents self-reference) */
  excludeWorkflowId?: string | null;
}

/**
 * Emits for WorkflowSelectorDialog component
 */
export interface WorkflowSelectorDialogEmits {
  /** Update dialog visibility */
  (e: 'update:modelValue', value: boolean): void;

  /** Emit selected workflow */
  (e: 'select', workflow: WorkflowSelectorItem): void;

  /** Emit cancel action */
  (e: 'cancel'): void;
}
