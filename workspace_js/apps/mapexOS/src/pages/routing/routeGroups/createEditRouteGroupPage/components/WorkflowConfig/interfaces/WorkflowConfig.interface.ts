/** TYPE IMPORTS */
import type { WorkflowData } from '@interfaces/routing/routeGroups.interface';

/**
 * Props for WorkflowConfig component
 */
export interface WorkflowConfigProps {
  /** Workflow data from the router */
  workflow: WorkflowData;

  /** Translation object from useRouteGroupsTranslations */
  t: Record<string, any>;
}

/**
 * Emits for WorkflowConfig component
 */
export interface WorkflowConfigEmits {
  (e: 'update:workflow', value: WorkflowData): void;
}
