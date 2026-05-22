/**
 * Props for EventFieldSelectorDrawer component
 * Drawer to select event fields from Asset Templates
 */
export interface EventFieldSelectorDrawerProps {
  /** Drawer visibility state */
  modelValue: boolean;

  /** Array of selected Asset Template IDs */
  selectedTemplates: string[];

  /** Current field value (for highlighting) */
  currentValue?: string;
}

/**
 * Emits for EventFieldSelectorDrawer component
 */
export interface EventFieldSelectorDrawerEmits {
  /** Update drawer visibility */
  (e: 'update:modelValue', value: boolean): void;

  /** Emitted when user selects a field */
  (e: 'select', field: string): void;

  /** Emitted when user wants to manage templates */
  (e: 'manage-templates'): void;
}

/**
 * Field information structure
 */
export interface FieldInfo {
  /** Field path (e.g., "data.temperature") */
  path: string;

  /** Template ID that provides this field */
  templateId: string;

  /** Template name for display */
  templateName: string;
}

/**
 * Grouped fields by template
 */
export interface GroupedFields {
  /** Template ID */
  templateId: string;

  /** Template name */
  templateName: string;

  /** Array of field paths */
  fields: string[];

  /** Visibility state in filter */
  visible: boolean;
}

/**
 * Template filter state
 * Map of templateId -> visible state
 */
export interface TemplateFilterState {
  [templateId: string]: boolean;
}
