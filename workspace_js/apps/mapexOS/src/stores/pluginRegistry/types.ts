import type { WorkflowPlugin, PluginNodeType, PluginCategory } from '@src/pages/automations/workflows/createEditWorkflowPage/interfaces';

/**
 * Plugin registry store state
 */
export interface PluginRegistryState {
  /** All registered plugins keyed by plugin ID */
  plugins: Map<string, WorkflowPlugin>;

  /** All node types keyed by type string (e.g., 'core/delay') */
  nodeTypeMap: Map<string, PluginNodeType>;
}

export type { WorkflowPlugin, PluginNodeType, PluginCategory };
