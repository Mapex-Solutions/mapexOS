/**
 * Field type for Field Input Selector
 * Represents different types of field values that can be selected
 */
export type FieldType = 'event' | 'state' | 'variable' | 'literal';

/**
 * Field Input Selector component props
 * Simple selector for choosing field input type
 */
export interface FieldInputSelectorProps {
  /** Current selected field type */
  modelValue?: FieldType | undefined;
  /** Label for the selector */
  label?: string | undefined;
  /** Whether the selector is disabled */
  disabled?: boolean | undefined;
  /** Whether to show icons in options */
  showIcons?: boolean | undefined;
}

/**
 * Field Input Selector component emits
 */
export interface FieldInputSelectorEmits {
  /** Emit when model value changes */
  (e: 'update:modelValue', value: FieldType): void;
}
