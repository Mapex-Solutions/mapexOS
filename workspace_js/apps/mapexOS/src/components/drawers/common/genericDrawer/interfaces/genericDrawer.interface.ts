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
 * Vertical scroll metrics forwarded from the drawer's internal scroll area
 * (a subset of Quasar's QScrollArea @scroll payload). Consumers use these to
 * drive infinite scroll — e.g. load the next page when nearing the bottom.
 */
export interface DrawerScrollInfo {
  verticalPosition: number;
  verticalPercentage: number;
  verticalSize: number;
  verticalContainerSize: number;
}

/**
 * Emits for GenericDrawer component
 */
export interface GenericDrawerEmits {
  (e: 'update:modelValue', value: boolean): void;
  (e: 'close'): void;
  (e: 'scroll', info: DrawerScrollInfo): void;
}
