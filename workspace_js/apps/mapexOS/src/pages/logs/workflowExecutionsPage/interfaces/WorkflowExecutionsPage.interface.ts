/**
 * Cursor state for workflow executions pagination
 */
export interface WorkflowExecutionsCursor {
  current?: string | undefined;
  next?: string | undefined;
  prev?: string | undefined;
  hasNext: boolean;
  hasPrevious: boolean;
}

/**
 * Filter state for workflow executions page
 */
export interface WorkflowExecutionsFilters {
  workflowName?: string | undefined;
  status?: string | undefined;
  instanceId?: string | undefined;
  definitionId?: string | undefined;
  startTime?: string | undefined;
  endTime?: string | undefined;
  includeChildren?: boolean | undefined;
}

/**
 * Normalized execution item for display (merges hot + cold sources)
 */
export interface WorkflowExecutionItem {
  id: string;
  workflowUUID?: string | undefined;
  instanceId: string;
  definitionId: string;
  workflowName: string;
  instanceName: string;
  definitionName: string;
  status: string;
  durationMs: number;
  errorMessage?: string | undefined;
  executionPath?: string | undefined;
  nodeOutputs?: string | undefined;
  errorInfo?: string | undefined;
  eventPayload?: string | undefined;
  triggerSource?: string | undefined;
  parentExecutionId?: string | undefined;
  depth: number;
  created: string;
  finished?: string | undefined;
  source: 'hot' | 'cold';
}
