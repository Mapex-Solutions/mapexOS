// ────────────────────────────────────────────────────────────────────────────
// Re-exports from @src/components/workflow/interfaces
// These types are re-exported here for backward compatibility so existing
// imports throughout the app continue to work.
// ────────────────────────────────────────────────────────────────────────────

import type {
  WorkflowNode,
  WorkflowEdge,
  WorkflowVariable,
  CaptureField,
  ExternalSignal,
  ExternalVariable,
} from '@src/components/workflow/interfaces';

export type {
  // Plugin system
  NodePropertyType,
  NodePropertyDefinition,
  DynamicNodeFormProps,
  DynamicNodeFormEmits,
  PluginCategory,
  HandleDefinition,
  HandleResolver,
  HandleOverrides,
  ValidationResult,
  PluginNodeType,
  WorkflowPlugin,
  CatalogGroup,
  ResolvedHandles,
  Disposable,
  PluginActivationContext,

  // Data models
  WorkflowNode,
  WorkflowEdge,
  WorkflowVariable,
  CaptureField,
  ExternalSignal,
  ExternalVariable,

  // Component props (SDK)
  BaseWorkflowNodeProps,
  WorkflowNodeComponentProps,
  NodeConfigComponentProps,
  NodeConfigComponentEmits,

  // Context
  IWorkflowEditorContext,
} from '@src/components/workflow/interfaces';

// ────────────────────────────────────────────────────────────────────────────
// App-specific types (NOT in SDK — specific to the workflow editor host)
// ────────────────────────────────────────────────────────────────────────────

import type { TimezoneConfig, RetryPolicy } from '../components/GeneralTab/interfaces/GeneralTab.interface';

export type { TimezoneConfig, RetryPolicy };

/**
 * General settings for the workflow (Tab 1)
 */
export interface WorkflowGeneralSettings {
  /** Workflow name (required) */
  name: string;

  /** Description */
  description: string;

  /** Whether workflow is enabled */
  enabled: boolean;

  /** Whether this is a reusable template */
  isTemplate: boolean;

  /** Share with child organizations */
  sharedWithChildren: boolean;

  /** Timezone configuration */
  timezone: TimezoneConfig;

  /** Global retry policy */
  retryPolicy: RetryPolicy;
}

/**
 * Canvas viewport metadata
 */
export interface CanvasViewport {
  /** X offset */
  x: number;

  /** Y offset */
  y: number;

  /** Zoom level */
  zoom: number;
}

/**
 * Snapshot of canvas state for undo/redo history
 */
export interface HistorySnapshot {
  /** Deep-cloned nodes at time of snapshot */
  nodes: WorkflowNode[];

  /** Deep-cloned edges at time of snapshot */
  edges: WorkflowEdge[];

  /** Human-readable action label (e.g., 'Add node', 'Delete') */
  label: string;
}

/**
 * Complete workflow definition (JSON DSL)
 */
export interface WorkflowDefinition {
  /** Workflow ID (set by backend) */
  workflowId?: string;

  /** Workflow name */
  name: string;

  /** Description */
  description: string;

  /** Whether workflow is enabled */
  enabled: boolean;

  /** Whether this is a reusable template */
  isTemplate: boolean;

  /** Definition schema version */
  definitionVersion: number;

  /** Timezone configuration */
  timezone: TimezoneConfig;

  /** Global retry policy */
  retryPolicy: RetryPolicy;

  /** State variables (persisted during execution) */
  states: WorkflowVariable[];

  /** Capture fields (stored in ClickHouse) */
  captureFields: CaptureField[];

  /** External inputs (input contract for callers) */
  externalInputs: ExternalVariable[];

  /** External signals (signal contract for wait_signal nodes) */
  externalSignals: ExternalSignal[];

  /** DAG nodes */
  nodes: WorkflowNode[];

  /** DAG edges (connections) */
  edges: WorkflowEdge[];

  /** Installed marketplace plugin IDs (e.g., ['telegram', 'slack']) */
  installedPlugins: string[];

  /** Missing marketplace plugin IDs (computed by backend) */
  missingPlugins?: string[];

  /** Definition status computed by backend: 'valid' | 'plugin_missing' | 'invalid' */
  status?: 'valid' | 'plugin_missing' | 'invalid';

  /** UI metadata (ignored by backend engine) */
  metadata: {
    canvasViewport: CanvasViewport;
  };
}

// ────────────────────────────────────────────────────────────────────────────
// Canvas & Toolbar
// ────────────────────────────────────────────────────────────────────────────

/**
 * Canvas toolbar state
 */
export interface CanvasToolbarState {
  /** Show minimap overlay */
  showMinimap: boolean;

  /** Show background grid */
  showGrid: boolean;

  /** Lock node movement */
  locked: boolean;

  /** Whether the canvas is maximized (fullscreen dialog) */
  maximized: boolean;
}

/**
 * Tab configuration for the page
 */
export interface WorkflowTab {
  /** Tab HTML id (for tour targeting) */
  id: string;

  /** Tab name (v-model value) */
  name: string;

  /** Display label */
  label: string;

  /** Material icon name */
  icon: string;

  /** Optional badge count */
  badge?: number;

  /** Badge color */
  badgeColor?: string;
}
