/**
 * UserSelectorDrawer props interface
 * Controls drawer visibility and behavior
 */
export interface UserSelectorDrawerProps {
  /** Controls drawer open/close state */
  modelValue: boolean;

  /** Selected user ID (for highlighting in list) */
  selectedUserId?: string | null;
}

/**
 * UserSelectorDrawer emits interface
 * Events emitted by the drawer component
 */
export interface UserSelectorDrawerEmits {
  /** Update drawer visibility */
  'update:modelValue': [value: boolean];

  /** Emit when user is selected */
  'select': [user: any];

  /** Emit when drawer is cancelled/closed */
  'cancel': [];
}
