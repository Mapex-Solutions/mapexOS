/**
 * RoleSelectorDrawer props interface
 * Controls drawer visibility and behavior
 */
export interface RoleSelectorDrawerProps {
  /** Controls drawer open/close state */
  modelValue: boolean;

  /** Selected role ID (for highlighting in list) */
  selectedRoleId?: string | null;
}

/**
 * RoleSelectorDrawer emits interface
 * Events emitted by the drawer component
 */
export interface RoleSelectorDrawerEmits {
  /** Update drawer visibility */
  'update:modelValue': [value: boolean];

  /** Emit when role is selected */
  'select': [role: any];

  /** Emit when drawer is cancelled/closed */
  'cancel': [];
}
