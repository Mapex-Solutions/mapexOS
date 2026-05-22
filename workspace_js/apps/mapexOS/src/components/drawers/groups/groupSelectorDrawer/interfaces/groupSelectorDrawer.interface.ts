/**
 * GroupSelectorDrawer props interface
 * Controls drawer visibility and behavior
 */
export interface GroupSelectorDrawerProps {
  /** Controls drawer open/close state */
  modelValue: boolean;

  /** Selected group ID (for highlighting in list) */
  selectedGroupId?: string | null;
}

/**
 * GroupSelectorDrawer emits interface
 * Events emitted by the drawer component
 */
export interface GroupSelectorDrawerEmits {
  /** Update drawer visibility */
  'update:modelValue': [value: boolean];

  /** Emit when group is selected */
  'select': [group: any];

  /** Emit when drawer is cancelled/closed */
  'cancel': [];
}
