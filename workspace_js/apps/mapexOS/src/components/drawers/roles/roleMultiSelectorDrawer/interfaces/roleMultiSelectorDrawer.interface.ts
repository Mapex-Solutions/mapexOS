/**
 * RoleMultiSelectorDrawer props interface
 * Controls drawer visibility and behavior for multi-select roles
 */
export interface RoleMultiSelectorDrawerProps {
  /** Controls drawer open/close state */
  modelValue: boolean;

  /** Selected role IDs (for highlighting in list) */
  selectedRoleIds?: string[];
}

/**
 * RoleMultiSelectorDrawer emits interface
 * Events emitted by the drawer component
 */
export interface RoleMultiSelectorDrawerEmits {
  /** Update drawer visibility */
  'update:modelValue': [value: boolean];

  /** Emit when roles selection is confirmed */
  'confirm': [roles: any[]];

  /** Emit when drawer is cancelled/closed */
  'cancel': [];
}
