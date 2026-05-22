/**
 * Re-export canonical types — FieldSourceValue, SourceType, NodeOutputOption
 * are defined in @src/components/workflow/interfaces for cross-plugin reusability.
 */
export type { SourceType, FieldSourceValue, NodeOutputOption } from '@src/components/workflow/interfaces';

/** TYPE IMPORTS */
import type { FieldSourceValue, SourceType, NodeOutputOption } from '@src/components/workflow/interfaces';

/**
 * Props for FieldSourceSelector component
 */
export interface FieldSourceSelectorProps {
  /** Current value (v-model) */
  modelValue: FieldSourceValue;
  /** Allowed source types — component adapts showing only these */
  allowedTypes: SourceType[];
  /** Label for the field (optional) */
  label?: string;
  /** Placeholder text (optional) */
  placeholder?: string;
  /** Whether the component is disabled */
  disabled?: boolean;

  /** Whether templates are available for event browsing (event type) */
  hasTemplates?: boolean;
  /** Count of available templates (event type) */
  templateCount?: number;

  /** Available state field options for autocomplete (state type) */
  stateFields?: Array<{ name: string; type: string }>;

  /** Nodes available on the canvas for nodeOutput type */
  nodeOutputOptions?: NodeOutputOption[];

  /** Credential ID for fetchOptions API calls (required when allowedTypes includes 'fetchOptions') */
  credentialId?: string;

  /** Resource key for fetchOptions endpoint (e.g., 'getChats') — from prop.fetchOptions.rules[].key */
  fetchOptionsKey?: string;

  /** Dynamic label for fetchOptions source type (e.g., 'Search Chats') — from prop.fetchOptions.rules[].label */
  fetchOptionsLabel?: string;
}

/**
 * Emits for FieldSourceSelector component
 */
export interface FieldSourceSelectorEmits {
  /** Emit when model value changes */
  (e: 'update:modelValue', value: FieldSourceValue): void;
  /** Parent should open the event field selector (modal or drawer) */
  (e: 'openEventSelector'): void;
  /** Parent should open the template selector (modal or drawer) */
  (e: 'openTemplateSelector'): void;
}
