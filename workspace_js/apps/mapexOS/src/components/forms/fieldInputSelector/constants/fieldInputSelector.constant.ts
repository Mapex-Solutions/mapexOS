/**
 * Field type options for Field Input Selector
 * Defines available field types with their visual properties
 */
export const FIELD_TYPE_SELECTOR_OPTIONS = [
  { label: 'Event', value: 'event', icon: 'event', color: 'blue-6' },
  { label: 'State', value: 'state', icon: 'storage', color: 'purple-6' },
  { label: 'Variable', value: 'variable', icon: 'code', color: 'orange-6' },
  { label: 'Literal', value: 'literal', icon: 'format_quote', color: 'green-6' }
] as const;

/**
 * Default field type
 */
export const DEFAULT_FIELD_TYPE = 'literal' as const;
