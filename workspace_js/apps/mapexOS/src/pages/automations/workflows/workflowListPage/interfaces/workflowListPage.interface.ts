/**
 * Filter state for workflow list page
 */
export interface WorkflowListPageFilters {
  /** Search by name */
  name: string | undefined;
  /** Filter by enabled status */
  status: boolean | undefined;
}

/**
 * Column visibility state for workflow list page
 */
export interface WorkflowListPageColumnVisibility {
  /** Show version column */
  version: boolean;
  /** Show nodes count column */
  nodesCount: boolean;
  /** Show plugins count column */
  pluginsCount: boolean;
}

/**
 * Workflow item for list display
 */
export interface WorkflowListItem {
  /** Unique workflow ID */
  id: string;
  /** Workflow name */
  name: string;
  /** Description */
  description: string;
  /** Whether the workflow is enabled */
  enabled: boolean;
  /** Whether this is a reusable template */
  isTemplate: boolean;
  /** Definition schema version */
  definitionVersion: number;
  /** Number of nodes in the workflow */
  nodesCount: number;
  /** Number of edges in the workflow */
  edgesCount: number;
  /** Timezone identifier */
  timezone: string;
  /** Created timestamp */
  created: string;
  /** Updated timestamp */
  updated: string;
  /** Definition status computed by backend */
  status: 'valid' | 'plugin_missing' | 'invalid';
  /** Missing plugin IDs */
  missingPlugins: string[];
  /** Number of installed plugins */
  pluginsCount: number;
}
