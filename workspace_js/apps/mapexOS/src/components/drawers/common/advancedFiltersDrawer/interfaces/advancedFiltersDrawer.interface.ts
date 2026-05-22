/**
 * Toggle field option
 */
export interface FilterToggleOption {
  label: string;
  value: any;
}

/**
 * Autocomplete option returned by fetch function
 */
export interface FilterAutocompleteOption {
  id: string;
  label: string;
  caption?: string;
}

/**
 * Base filter field configuration
 */
export interface FilterFieldBase {
  /** Unique key for the filter */
  key: string;
  /** Display label */
  label: string;
  /** Icon name */
  icon: string;
  /** Whether the field is disabled */
  disabled?: boolean;
}

/**
 * Toggle filter field (btn-toggle with options)
 */
export interface FilterFieldToggle extends FilterFieldBase {
  type: 'toggle';
  /** Options for the toggle */
  options: FilterToggleOption[];
}

/**
 * Switch filter field (simple boolean on/off)
 */
export interface FilterFieldSwitch extends FilterFieldBase {
  type: 'switch';
  /** Placeholder text / description */
  placeholder?: string;
}

/**
 * Autocomplete filter field (q-select with search)
 */
export interface FilterFieldAutocomplete extends FilterFieldBase {
  type: 'autocomplete';
  /** Placeholder text */
  placeholder?: string;
  /** Function to fetch options based on search term */
  fetchOptions: (search: string) => Promise<FilterAutocompleteOption[]>;
}

/**
 * Input filter field (q-input)
 */
export interface FilterFieldInput extends FilterFieldBase {
  type: 'input';
  /** Placeholder text */
  placeholder?: string;
  /** Input type (text, number, etc) */
  inputType?: 'text' | 'number' | 'email' | 'password' | 'tel' | 'url' | 'search' | 'date' | 'time' | 'datetime-local';
}

/**
 * Select filter field (q-select with static options)
 */
export interface FilterFieldSelect extends FilterFieldBase {
  type: 'select';
  /** Placeholder text */
  placeholder?: string;
  /** Static options */
  options: FilterToggleOption[];
  /** Loading state for async options */
  loading?: boolean;
}

/**
 * Union type for all filter field types
 */
export type FilterField =
  | FilterFieldToggle
  | FilterFieldSwitch
  | FilterFieldAutocomplete
  | FilterFieldInput
  | FilterFieldSelect;

/**
 * Filter values object
 */
export type FilterValues = Record<string, any>;

/**
 * Props for AdvancedFiltersDrawer component
 */
export interface AdvancedFiltersDrawerProps {
  /** Controls drawer visibility (v-model) */
  modelValue: boolean;
  /** Drawer title */
  title?: string;
  /** Filter fields configuration */
  fields: FilterField[];
  /** Current filter values */
  values: FilterValues;
  /** Drawer width in pixels */
  width?: number;
}

/**
 * Emits for AdvancedFiltersDrawer component
 */
export interface AdvancedFiltersDrawerEmits {
  (e: 'update:modelValue', value: boolean): void;
  (e: 'apply', values: FilterValues): void;
  (e: 'reset'): void;
  (e: 'field-change', key: string, value: any): void;
  /** Emitted when the pending changes state changes */
  (e: 'pending-change', hasPending: boolean): void;
}
