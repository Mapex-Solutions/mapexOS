/**
 * Props for WorkflowInstanceDetailsDrawer component.
 * Displays full instance details in a right-side drawer.
 */
export interface WorkflowInstanceDetailsDrawerProps {
  /** Drawer visibility state */
  modelValue: boolean;

  /** Instance ID to display (null when closed) */
  instanceId: string | null;
}

/**
 * Emits for WorkflowInstanceDetailsDrawer component.
 */
export interface WorkflowInstanceDetailsDrawerEmits {
  /** Update drawer visibility */
  (e: 'update:modelValue', value: boolean): void;

  /** Edit action — emit instance ID */
  (e: 'edit', instanceId: string): void;
}
