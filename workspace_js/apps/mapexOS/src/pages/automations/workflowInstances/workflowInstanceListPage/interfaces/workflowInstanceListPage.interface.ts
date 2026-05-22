/**
 * Filter state for workflow instance list page
 */
export interface WorkflowInstanceListPageFilters {
  /** Search by name */
  name: string | undefined;
  /** Filter by enabled status */
  status: boolean | undefined;
}

/**
 * Column visibility state for workflow instance list page
 */
export interface WorkflowInstanceListPageColumnVisibility {
  /** Show definition name column */
  definitionName: boolean;
  /** Show inputs count column */
  inputsCount: boolean;
  /** Show unique execution column */
  uniqueExecution: boolean;
}

/**
 * Workflow instance item for list display
 */
export interface WorkflowInstanceListItem {
  /** Unique instance ID */
  id: string;
  /** Instance name */
  name: string;
  /** Description */
  description: string;
  /** Whether the instance is enabled */
  enabled: boolean;
  /** Definition name (denormalized) */
  definitionName: string;
  /** Number of external inputs */
  inputsCount: number;
  /** Whether this uses unique execution */
  uniqueExecution: boolean;
  /** Fixed UUID (only when uniqueExecution is true) */
  workflowUUID: string;
}
