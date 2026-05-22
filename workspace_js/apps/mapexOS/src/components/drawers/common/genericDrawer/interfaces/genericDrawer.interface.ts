/**
 * Props for GenericDrawer component
 * Base drawer with standardized header, scrollable content, and optional footer
 */
export interface GenericDrawerProps {
  /** Controls drawer visibility (v-model) */
  modelValue: boolean;

  /** Drawer title displayed in header */
  title: string;

  /** Icon name displayed before title */
  icon?: string;

  /** Icon color (Quasar color name) */
  iconColor?: string;

  /** Drawer width in pixels */
  width?: number;

  /** Tooltip text for close button */
  closeTooltip?: string;
}

/**
 * Emits for GenericDrawer component
 */
export interface GenericDrawerEmits {
  (e: 'update:modelValue', value: boolean): void;
  (e: 'close'): void;
}
