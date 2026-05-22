/**
 * Single item in the icon section navigation menu
 */
export interface IconSectionNavItem {
  /** Unique section key used for v-model binding */
  name: string;

  /** Material icon name */
  icon: string;

  /** Tooltip text shown on hover */
  tooltip: string;

  /** Show a dot badge indicator on this item */
  badge?: boolean;

  /** Badge dot color (Quasar color name, default: 'primary') */
  badgeColor?: string;
}

/**
 * Props for IconSectionNav component
 */
export interface IconSectionNavProps {
  /** Active section name (v-model) */
  modelValue: string;

  /** Menu items to display */
  items: IconSectionNavItem[];

  /** Width of the navigation rail in pixels (default: 40) */
  width?: number;
}

/**
 * Emits for IconSectionNav component
 */
export interface IconSectionNavEmits {
  /**
   * Emitted when active section changes
   * @param e - Event name
   * @param value - New section name
   */
  (e: 'update:modelValue', value: string): void;
}
