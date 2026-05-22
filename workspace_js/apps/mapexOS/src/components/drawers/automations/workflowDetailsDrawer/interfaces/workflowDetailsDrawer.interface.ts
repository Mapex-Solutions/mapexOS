/**
 * Props for WorkflowDetailsDrawer component
 */
export interface WorkflowDetailsDrawerProps {
  /** Whether the drawer is open */
  modelValue: boolean;

  /** ID of the workflow to display details for */
  workflowId: string | null;
}

/**
 * Emits for WorkflowDetailsDrawer component
 */
export interface WorkflowDetailsDrawerEmits {
  (e: 'update:modelValue', value: boolean): void;
  (e: 'edit', workflowId: string): void;
}
