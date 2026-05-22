/**
 * Plugin catalog category labels and icons
 */
export const CATALOG_CATEGORIES = [
  { category: 'triggers', label: 'Triggers', icon: 'bolt' },
  { category: 'logic', label: 'Logic & Conditions', icon: 'rule' },
  { category: 'state', label: 'Data', icon: 'edit_note' },
  { category: 'flow_control', label: 'Flow Control', icon: 'call_split' },
  { category: 'timers', label: 'Timers & Wait', icon: 'hourglass_empty' },
  { category: 'integrations', label: 'Integrations', icon: 'http' },
  { category: 'observability', label: 'Observability', icon: 'article' },
  { category: 'annotations', label: 'Annotations', icon: 'sticky_note_2' },
  { category: 'custom', label: 'Custom', icon: 'extension' },
] as const;

/**
 * Plugin catalog sidebar width (expanded)
 */
export const CATALOG_WIDTH = 260;

/**
 * Plugin catalog sidebar width (collapsed)
 */
export const CATALOG_WIDTH_COLLAPSED = 48;
