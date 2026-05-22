import type { CanvasToolbarState } from '../../../interfaces/CreateEditWorkflow.interface';

/**
 * Props for WorkflowCanvas component
 */
export interface WorkflowCanvasProps {
  /** Toolbar state for canvas configuration */
  toolbarState: CanvasToolbarState;
}

/**
 * Emits for WorkflowCanvas component
 */
export interface WorkflowCanvasEmits {
  (e: 'node-select', nodeId: string | null): void;
  (e: 'canvas-click'): void;
}
