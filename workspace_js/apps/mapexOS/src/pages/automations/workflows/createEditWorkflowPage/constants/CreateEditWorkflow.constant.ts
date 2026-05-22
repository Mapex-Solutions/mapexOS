import type { WorkflowGeneralSettings } from '../interfaces';

/**
 * Default general settings for a new workflow
 */
export const DEFAULT_GENERAL_SETTINGS: WorkflowGeneralSettings = {
  name: '',
  description: '',
  enabled: true,
  isTemplate: false,
  sharedWithChildren: false,
  timezone: {
    type: 'literal',
    value: 'UTC',
  },
  retryPolicy: {
    enabled: false,
    maxAttempts: 3,
    initialInterval: '1s',
    backoffMultiplier: 2.0,
    maxInterval: '5m',
    nonRetryableErrors: [],
  },
};

/**
 * Snap grid size in pixels
 */
export const SNAP_GRID_SIZE: [number, number] = [15, 15];

/**
 * Variable type options for workflow variables and capture fields
 */
export const VARIABLE_TYPE_OPTIONS = [
  { label: 'String', value: 'string' },
  { label: 'Number', value: 'number' },
  { label: 'Boolean', value: 'boolean' },
  { label: 'JSON', value: 'json' },
] as const;

/**
 * Extended type options for external inputs (includes asset-based types)
 */
export const EXTERNAL_INPUT_TYPE_OPTIONS = [
  { label: 'String', value: 'string' },
  { label: 'Number', value: 'number' },
  { label: 'Boolean', value: 'boolean' },
  { label: 'JSON', value: 'json' },
  { label: 'Literal', value: 'literal' },
  { label: 'Asset from Template', value: 'assetFromTemplate' },
] as const;

/**
 * Fixed asset field path options — metadata fields from the Asset entity
 */
export const ASSET_FIELD_PATH_OPTIONS = [
  { label: 'Asset UUID', value: 'assetUUID', icon: 'fingerprint' },
  { label: 'Asset Name', value: 'name', icon: 'label' },
  { label: 'Asset ID', value: 'id', icon: 'tag' },
] as const;

/**
 * Default value by variable type
 */
export const DEFAULT_VALUE_BY_TYPE: Record<string, string | number | boolean> = {
  string: '',
  number: 0,
  boolean: false,
  json: '{}',
};

/**
 * Base node width estimate for dagre layout (px).
 * Nodes with multiple outputs grow wider automatically.
 */
export const LAYOUT_NODE_WIDTH = 120;

/**
 * Default node height estimate for dagre layout (px)
 */
export const LAYOUT_NODE_HEIGHT = 70;

/**
 * Maximum number of undo history entries
 */
export const MAX_HISTORY_SIZE = 50;

