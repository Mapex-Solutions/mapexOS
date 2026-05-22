/** TYPE IMPORTS */
import type { FieldSourceValue } from '@src/components/workflow/interfaces';

/**
 * @deprecated Use `FieldSourceValue` from `@src/components/workflow/interfaces` instead.
 * Kept as alias for backward compatibility with rules pages.
 */
export type FieldValue = FieldSourceValue;

/**
 * Event Field Input component props
 * Manages input for event fields with type selection and browsing capabilities
 */
export interface EventFieldInputProps {
  /** Current field value */
  modelValue?: FieldValue | undefined;
  /** Label for the input */
  label?: string | undefined;
  /** Placeholder text */
  placeholder?: string | undefined;
  /** Whether the input is disabled */
  disabled?: boolean | undefined;
  /** Whether to show the type selector */
  showTypeSelector?: boolean | undefined;
  /** Whether templates are available for event browsing */
  hasTemplates?: boolean | undefined;
  /** Count of available templates */
  templateCount?: number | undefined;
  /** Whether state fields are available */
  hasStateFields?: boolean | undefined;
  /** Count of available state fields */
  stateFieldCount?: number | undefined;
  /** Available state field options */
  stateFields?: Array<{ name: string; type: string }> | undefined;
}

/**
 * Event Field Input component emits
 */
export interface EventFieldInputEmits {
  /** Emit when model value changes */
  (e: 'update:modelValue', value: FieldValue): void;
  /** Emit when event field selector should open */
  (e: 'openEventSelector'): void;
  /** Emit when template selector should open */
  (e: 'openTemplateSelector'): void;
  /** Emit when state field is selected */
  (e: 'selectStateField', fieldName: string): void;
}
