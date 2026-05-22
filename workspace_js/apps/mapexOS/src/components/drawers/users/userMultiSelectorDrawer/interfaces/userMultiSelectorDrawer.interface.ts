/**
 * UserMultiSelectorDrawer Interfaces
 *
 * Interfaces for the user multi-selector drawer component
 * that supports paginated user selection with search toggle.
 */

/**
 * User item for display in selector
 */
export interface UserSelectorItem {
  /** User ID */
  id: string;

  /** First name */
  firstName: string;

  /** Last name */
  lastName: string;

  /** Email address */
  email: string;

  /** Optional avatar URL */
  avatar?: string;
}

/**
 * Filter mode for user search
 */
export type UserFilterMode = 'name' | 'email';

/**
 * UserMultiSelectorDrawer props interface
 */
export interface UserMultiSelectorDrawerProps {
  /** Controls drawer open/close state */
  modelValue: boolean;

  /** IDs of users to exclude from the list (already selected elsewhere) */
  excludeUserIds?: string[];

  /** Pre-selected user IDs (for highlighting in list) */
  selectedUserIds?: string[];
}

/**
 * UserMultiSelectorDrawer emits interface
 */
export interface UserMultiSelectorDrawerEmits {
  /** Update drawer visibility */
  'update:modelValue': [value: boolean];

  /** Emit when user selection is confirmed */
  'confirm': [users: UserSelectorItem[]];

  /** Emit when drawer is cancelled/closed */
  'cancel': [];
}
