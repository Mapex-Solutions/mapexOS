/**
 * CreateEditGroupPage Interfaces
 */

import type { UserSelectorItem } from '@components/drawers';

/**
 * Role selection item for display
 */
export interface RoleSelectionItem {
  /** Role ID */
  id: string;

  /** Role name */
  name: string;
}

/**
 * Group form data for create/edit operations
 */
export interface GroupFormData {
  /** Group name */
  name: string;

  /** Group description */
  description: string;

  /** Whether the group is enabled */
  enabled: boolean;
}

/**
 * Pending member changes for the group
 */
export interface PendingMemberChanges {
  /** Users to add to the group */
  additions: UserSelectorItem[];

  /** User IDs to remove from the group */
  removals: string[];
}

/**
 * Group form state for tracking all form-related state
 */
export interface GroupFormState {
  /** Selected member IDs */
  selectedMembers: string[];

  /** Selected roles for the group */
  selectedRoles: RoleSelectionItem[];

  /** Pending member changes */
  pendingChanges: PendingMemberChanges;

  /** Whether form is saving */
  isSaving: boolean;

  /** Current step number */
  currentStep: number;
}
