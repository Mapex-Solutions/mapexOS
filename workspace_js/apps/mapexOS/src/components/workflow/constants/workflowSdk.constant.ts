/**
 * Start node type identifier
 */
export const START_NODE_TYPE = 'core/start';

/**
 * Default start node ID (auto-created on every workflow)
 */
export const START_NODE_ID = '__start__';

/**
 * @deprecated Use outputHints on PluginNodeType instead — only nodes with outputHints
 * appear in the nodeOutput selector automatically.
 */
export const NODE_OUTPUT_EXCLUDED_TYPES: readonly string[] = [
  'core/start',
  'core/end',
  'core/text_note',
  'core/group_frame',
] as const;

/**
 * Handle position options for the context menu in BaseWorkflowNode
 */
export const POSITION_OPTIONS = [
  { label: 'Top', value: 'top', icon: 'arrow_upward' },
  { label: 'Bottom', value: 'bottom', icon: 'arrow_downward' },
  { label: 'Left', value: 'left', icon: 'arrow_back' },
  { label: 'Right', value: 'right', icon: 'arrow_forward' },
] as const;
