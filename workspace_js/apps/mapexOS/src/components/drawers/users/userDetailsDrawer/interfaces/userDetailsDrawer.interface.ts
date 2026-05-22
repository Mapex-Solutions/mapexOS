/**
 * UserDetailsDrawer Component Interface
 *
 * This drawer displays detailed information about a user by fetching
 * the full user data from the API using the user ID.
 *
 * IMPORTANT: This drawer follows the STANDARD pattern for all detail drawers:
 * - Receives only the user ID (not the full user object)
 * - Fetches complete data using apis.mapexOS.users.getById()
 * - Shows loading state while fetching
 * - Handles errors gracefully
 *
 * This ensures the drawer always displays complete and up-to-date information,
 * regardless of the projection used in the list view.
 */

/**
 * Props for UserDetailsDrawer component
 */
export interface UserDetailsDrawerProps {
  /**
   * Controls drawer visibility (v-model)
   */
  modelValue: boolean;

  /**
   * User ID to fetch and display
   * IMPORTANT: Pass only the ID, NOT the user object from the list
   */
  userId: string | null;
}

/**
 * Events emitted by UserDetailsDrawer component
 */
export interface UserDetailsDrawerEmits {
  /**
   * Emitted when drawer visibility changes
   */
  (e: 'update:modelValue', value: boolean): void;

  /**
   * Emitted when user clicks edit button
   * @param userId - The ID of the user to edit
   */
  (e: 'edit', userId: string): void;
}

/**
 * User group info for detail view
 */
export interface UserGroupInfo {
  id: string;
  name: string;
  description?: string;
}

/**
 * User membership info for detail view
 */
export interface UserMembershipInfo {
  orgId: string;
  orgName: string;
  orgType: string;
  scope: 'local' | 'recursive';
  roleNames: string[];
  via: string;
}

/**
 * Organized access data for BigTech card UI
 */
export interface OrganizedAccess {
  orgId: string;
  orgName: string;
  orgType: string;
  scope: 'local' | 'recursive';
  groups: string[];
  directRoles: string[];
}
