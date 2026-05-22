/**
 * Field type options for Event Field Input
 * Defines available input modes with their visual properties
 */
export const FIELD_TYPE_OPTIONS = [
  { label: 'Event', value: 'event', icon: 'event', color: 'blue-6' },
  { label: 'State', value: 'state', icon: 'storage', color: 'purple-6' },
  { label: 'Variable', value: 'variable', icon: 'code', color: 'orange-6' },
  { label: 'Literal', value: 'literal', icon: 'format_quote', color: 'green-6' }
] as const;

/**
 * Default field value for Event Field Input
 */
export const DEFAULT_FIELD_VALUE = {
  type: 'literal' as const,
  value: '',
  mode: 'dynamic' as const
};
