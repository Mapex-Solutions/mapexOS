/**
 * DynamicFiltersDrawer Interfaces
 * Drawer for selecting a source (Asset or Asset Template) and filtering by dynamic EVA fields
 */

/** Source type for dynamic filter resolution */
export type DynamicSourceType = 'asset' | 'assetTemplate';

/** Operator option for EVA filter select */
export interface EvaOperatorOption {
  /** Operator value sent to backend (eq, neq, gt, etc.) */
  value: string;
  /** i18n key for the operator label */
  labelKey: string;
  /** Material icon name for the operator */
  icon: string;
}

/**
 * Field definition loaded from the template (available pool)
 * Represents a field that can be added as an active filter
 */
export interface DynamicFieldDefinition {
  /** Field key/name from EVA schema */
  key: string;
  /** Field label for display */
  label: string;
  /** Field data type */
  type: 'string' | 'number' | 'boolean' | 'date';
  /** Original EVA type from template (e.g. 'bool', 'geo') */
  originalType?: string;
  /** uint16 fieldId for backend EVA MAP column access */
  fieldId: number;
}

/**
 * Dynamic field for EVA filters with operator support
 */
export interface DynamicFilterField {
  /** Field key/name from EVA schema */
  key: string;
  /** Field label for display */
  label: string;
  /** Field data type */
  type: 'string' | 'number' | 'boolean' | 'date';
  /** Selected comparison operator (eq, neq, gt, gte, lt, lte, between, like) */
  operator: string;
  /** Current filter value (typed per field.type at runtime) */
  value: any;
  /** End value for "between" operator (range end) */
  endValue?: any;
  /** Original EVA type from template (e.g. 'bool', 'geo') */
  originalType?: string;
  /** uint16 fieldId for backend EVA MAP column access */
  fieldId: number;
}

/**
 * Result emitted when user applies dynamic filters
 */
export interface DynamicFiltersResult {
  /** Source type selected by the user */
  sourceType: DynamicSourceType;
  /** ID of the selected source entity (asset or business rule) */
  sourceId: string;
  /** Name of the selected source entity */
  sourceName: string;
  /** Resolved asset template ID */
  assetTemplateId: string;
  /** Resolved template name */
  templateName: string;
  /** Dynamic fields with user-entered values and operators */
  fields: DynamicFilterField[];
}

/**
 * Props for DynamicFiltersDrawer component
 */
export interface DynamicFiltersDrawerProps {
  /** Controls drawer visibility (v-model) */
  modelValue: boolean;
  /** When true, populates the drawer with demo data for tour/onboarding */
  demo?: boolean;
}

/**
 * Emits for DynamicFiltersDrawer component
 */
export interface DynamicFiltersDrawerEmits {
  (e: 'update:modelValue', value: boolean): void;
  (e: 'apply', result: DynamicFiltersResult): void;
  (e: 'reset'): void;
  (e: 'pending-change', hasPending: boolean): void;
}
