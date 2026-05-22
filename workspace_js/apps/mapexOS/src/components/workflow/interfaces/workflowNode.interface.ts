/**
 * Timeout configuration for async nodes.
 * Lives at node level (not inside config) — separate concern from node functionality.
 */
export interface NodeTimeoutConfig {
  /** Duration value */
  duration: number;

  /** Time unit: 'seconds' | 'minutes' | 'hours' | 'days' */
  unit: string;

  /** When true, timeout routes to "timeout" output handle instead of failing execution */
  enableOutput: boolean;
}

/**
 * Error handler configuration for nodes with "error" output handle.
 * Lives at node level (not inside config) — separate concern from node functionality.
 * When enabled, retries the node with exponential backoff before following the "error" handle.
 * Independent from timeout — timeout and error handler never interact.
 */
export interface NodeErrorHandlerConfig {
  /** Whether retry is enabled for this node */
  enabled: boolean;

  /** Maximum number of retry attempts before following "error" handle */
  maxAttempts: number;

  /** Delay before first retry */
  initialInterval: number;

  /** Time unit for initialInterval: 'seconds' | 'minutes' | 'hours' */
  intervalUnit: string;

  /** Multiplier applied to interval after each failed attempt (exponential backoff) */
  backoffMultiplier: number;
}

/**
 * Workflow node instance on canvas
 */
export interface WorkflowNode {
  /** Unique node ID */
  id: string;

  /** Plugin node type (e.g., 'core/delay') */
  type: string;

  /** Position on canvas (relative to parent if parentNodeId is set) */
  position: { x: number; y: number };

  /** Node configuration */
  config: Record<string, unknown>;

  /** Async timeout configuration (nil = use plugin/platform default) */
  timeout?: NodeTimeoutConfig;

  /** Error handler with retry policy (nil = no retry, error goes to handle directly) */
  errorHandler?: NodeErrorHandlerConfig;

  /** Display label (user-editable) */
  label?: string;

  /** Parent node ID — child moves with parent, position becomes relative */
  parentNodeId?: string;
}

/**
 * Workflow edge (connection between nodes)
 */
export interface WorkflowEdge {
  /** Unique edge ID */
  id: string;

  /** Source node ID */
  source: string;

  /** Source handle ID */
  sourceHandle?: string;

  /** Target node ID */
  target: string;

  /** Target handle ID */
  targetHandle?: string;

  /** Edge label */
  label?: string;

  /** Horizontal offset (px) from natural edge center — for untangling overlapping edges */
  pathOffsetX?: number;

  /** Vertical offset (px) from natural edge center — for untangling overlapping edges */
  pathOffsetY?: number;
}

/**
 * Workflow variable definition
 */
export interface WorkflowVariable {
  /** Variable name (used as state.{field}) */
  field: string;

  /** Variable type */
  type: 'string' | 'number' | 'boolean' | 'json';

  /** Default value */
  defaultValue: string | number | boolean | Record<string, unknown>;

  /** Optional description */
  description?: string;

  /** Whether variable persists across workflow runs (false = ephemeral, true = durable) */
  durable: boolean;
}

/**
 * Capture field definition (for ClickHouse search/reports)
 */
export interface CaptureField {
  /** Field name */
  field: string;

  /** Field type */
  type: 'string' | 'number' | 'boolean' | 'json';

  /** Description of what this field captures */
  description: string;
}

/**
 * External signal definition (signal contract).
 * Defines named signals that this workflow can wait for via wait_signal nodes.
 */
export interface ExternalSignal {
  /** Signal name (unique identifier used in wait_signal nodes) */
  name: string;

  /** Description of what this signal represents */
  description: string;
}

/**
 * Allowed types for external input variables.
 * - Primitive types: string, number, boolean, json
 * - literal: fixed value embedded in the definition — not provided by callers
 * - assetFromTemplate: UI hint — renders an asset picker filtered by template
 */
export type ExternalVariableType = 'string' | 'number' | 'boolean' | 'json' | 'literal' | 'assetFromTemplate';

/**
 * External variable definition (input contract).
 * Defines variables that callers must provide when invoking this workflow
 * (e.g., from a subworkflow node, API trigger, or manual run).
 */
export interface ExternalVariable {
  /** Variable key (unique, used as input.{field} in expressions) */
  field: string;

  /** UI display label */
  label: string;

  /** Material icon name for UI presentation */
  icon: string;

  /** Variable type */
  type: ExternalVariableType;

  /** UI help text / description */
  description?: string;

  /** Default value when caller does not provide it */
  defaultValue: string | number | boolean | Record<string, unknown>;

  /** Whether caller must provide this variable */
  required: boolean;

  /** Asset template ID — only when type = 'assetFromTemplate' */
  assetTemplateId?: string;

  /** Asset field path to extract as value — only when type = 'assetFromTemplate' */
  fieldPath?: string;
}
