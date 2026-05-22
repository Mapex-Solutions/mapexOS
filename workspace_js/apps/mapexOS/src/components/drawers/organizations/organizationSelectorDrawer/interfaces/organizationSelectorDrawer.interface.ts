/**
 * OrganizationSelectorDrawer props interface
 * Controls drawer visibility and behavior
 */
export interface OrganizationSelectorDrawerProps {
  /** Controls drawer open/close state */
  modelValue: boolean;

  /** Selected organization ID (for highlighting in list) */
  selectedOrganizationId?: string | null;
}

/**
 * OrganizationSelectorDrawer emits interface
 * Events emitted by the drawer component
 */
export interface OrganizationSelectorDrawerEmits {
  /** Update drawer visibility */
  'update:modelValue': [value: boolean];

  /** Emit when organization is selected */
  'select': [organization: any];

  /** Emit when drawer is cancelled/closed */
  'cancel': [];
}
