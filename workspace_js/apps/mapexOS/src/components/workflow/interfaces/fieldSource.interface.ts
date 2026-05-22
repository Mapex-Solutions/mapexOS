/**
 * Source types available for field value selection.
 * Determines where the value comes from in the workflow context.
 */
export type SourceType = 'event' | 'state' | 'input' | 'variable' | 'literal' | 'nodeOutput' | 'fetchOptions' | 'assetStatus';

/**
 * Value model representing a field source selection.
 * Used across all workflow nodes that need dynamic value inputs
 * (conditions, set_state, loop, trigger variables, etc.)
 */
export interface FieldSourceValue {
  /** Source type — determines where the value comes from */
  type: SourceType;

  /** The actual value (field path, variable name, literal string, etc.) */
  value: string;

  /** Node ID — only required when type is 'nodeOutput' */
  nodeId?: string;

  /** Input mode — only for event type (dynamic: browse, manual: type path, assetStatus: predefined health fields) */
  mode?: 'dynamic' | 'manual' | 'assetStatus';
}

/**
 * Node output option for the nodeOutput source type selector.
 * Represents a node on the canvas that produces outputs.
 */
export interface NodeOutputOption {
  /** Node ID */
  id: string;

  /** Display label (e.g. "Loop_1 (Loop)") */
  label: string;

  /** Node type (e.g. "core/loop") */
  type: string;

  /** Hint of available outputs (e.g. "item, index") */
  outputHint?: string;
}

/**
 * Option for the asset status field dropdown.
 * Represents a predefined health monitoring event field (sensor.offline / sensor.online).
 */
export interface AssetStatusFieldOption {
  label: string;
  value: string;
  icon: string;
  type: string;
  availability: 'all' | 'offline';
}

/**
 * UI option for displaying a source type in selectors/menus.
 */
export interface SourceTypeOption {
  /** Display label */
  label: string;

  /** Value matching SourceType */
  value: SourceType;

  /** Material icon name */
  icon: string;

  /** Quasar color class */
  color: string;
}
