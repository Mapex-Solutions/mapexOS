/**
 * Props for AvailableFieldsList component
 */
export interface AvailableFieldsListProps {
  /** Array of field paths to display */
  fields: string[];

  /** Optional maximum height for scroll area (in pixels) */
  maxHeight?: number;

  /** Show loading state */
  loading?: boolean;
}

/**
 * Emits for AvailableFieldsList component
 */
export interface AvailableFieldsListEmits {
  /** Emitted when a field is clicked (for copy/select functionality) */
  'field-click': [field: string];
}
