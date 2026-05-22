import type { HandleDefinition } from './workflowPlugin.interface';

/**
 * Props for BaseWorkflowNode component
 */
export interface BaseWorkflowNodeProps {
  /** Node ID */
  id: string;

  /** Display label (may be undefined during Vue Flow internal re-renders) */
  label?: string;

  /** Material icon name */
  icon: string;

  /** Color theme for the node */
  color: string;

  /** Whether the node is selected */
  selected?: boolean;

  /** Input handles */
  inputs?: HandleDefinition[];

  /** Output handles */
  outputs?: HandleDefinition[];

  /** Whether the node can be deleted (default: true) */
  deletable?: boolean;

  /** Visual shape — 'square' (default) or 'circle' */
  shape?: 'square' | 'circle';

  /** Optional hex color override for --node-color (used by nodes with dynamic colors) */
  colorHex?: string;

  /** Whether the node has validation errors */
  hasErrors?: boolean;
}

/**
 * Shared props interface for all Vue Flow custom node components.
 * Vue Flow passes these props automatically to custom node components.
 */
export interface WorkflowNodeComponentProps {
  /** Node ID assigned by Vue Flow */
  id: string;

  /** Node data containing config, label, and validation state */
  data: { config: Record<string, unknown>; label?: string; __nodeType?: string; hasErrors?: boolean };

  /** Whether the node is currently selected */
  selected?: boolean;
}

/**
 * Shared props interface for all node config panel components.
 */
export interface NodeConfigComponentProps {
  /** Current node configuration */
  config: Record<string, unknown>;
}

/**
 * Shared emits interface for all node config panel components.
 */
export interface NodeConfigComponentEmits {
  (e: 'update:config', config: Record<string, unknown>): void;
}
