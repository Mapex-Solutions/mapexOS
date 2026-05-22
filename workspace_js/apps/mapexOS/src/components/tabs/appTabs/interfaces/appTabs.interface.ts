/**
 * Single tab item configuration
 */
export interface AppTabItem {
  /** Unique identifier for the tab */
  name: string;

  /** Display label for the tab */
  label: string;

  /** Material icon name */
  icon?: string | undefined;

  /** Optional badge count to display */
  badge?: number | undefined;

  /** Badge color (Quasar color name, default: 'primary') */
  badgeColor?: string | undefined;

  /** Whether the tab is disabled */
  disabled?: boolean | undefined;

  /** Optional HTML id attribute for targeting (e.g., tours) */
  id?: string | undefined;
}

/**
 * Props for AppTabs component
 */
export interface AppTabsProps {
  /** Currently active tab name (v-model) */
  modelValue: string;

  /** Array of tab configurations */
  tabs: AppTabItem[];

  /** Whether to wrap tabs in a bordered card */
  bordered?: boolean;

  /** Tab alignment: 'left' | 'center' | 'right' | 'justify' */
  align?: 'left' | 'center' | 'right' | 'justify';

  /** Whether to show separator below tabs */
  separator?: boolean;

  /** Visual variant: 'default' (underline indicator) or 'pill' (button-style with background) */
  variant?: 'default' | 'pill';
}

/**
 * Emits for AppTabs component
 */
export interface AppTabsEmits {
  (e: 'update:modelValue', value: string): void;
  (e: 'change', value: string): void;
}
