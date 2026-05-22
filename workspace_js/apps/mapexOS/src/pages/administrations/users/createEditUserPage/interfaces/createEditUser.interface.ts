/**
 * CreateEditUserPage Interfaces
 *
 * Based on contract: workspace_go/packages/contracts/services/mapexos/onboarding/dtos.go
 */

/**
 * Scope types for membership access
 */
export type ScopeType = 'local' | 'recursive';

/**
 * Access assignment type
 * - group: Add user to existing group (recommended)
 * - direct: Direct membership assignment (exception)
 * - both: Both group and direct access (edit mode supports this)
 */
export type AccessType = 'group' | 'direct' | 'both';

/**
 * Group access mode
 * - existing: Use an existing group (user just joins, inherits group's roles)
 * - new: Create a new group with specific roles
 */
export type GroupAccessMode = 'existing' | 'new';

/**
 * Existing group data (when user selects an existing group)
 * User joins the group and inherits its roles.
 * The group ALREADY has a Membership with roles defined.
 */
export interface ExistingGroupData {
  /** Group ID */
  groupId: string;

  /** Group name (for display) */
  groupName: string;
}

/**
 * New group data (when user creates a new group during onboarding)
 * Onboarding will delegate to Groups service to create:
 * - The Group entity
 * - The Membership for the group (with the specified roles)
 */
export interface NewGroupData {
  /** Group name (required, min 3, max 150) */
  name: string;

  /** Group description (optional, max 500) */
  description?: string;

  /** Role IDs to assign to the new group */
  roleIds: string[];

  /** Role names (for display) */
  roleNames: string[];
}

/**
 * Group access data (when accessType = 'group')
 * Supports two scenarios:
 * - existingGroup: User joins an existing group (inherits its roles)
 * - newGroup: Create a new group with specific roles, then user joins it
 */
export interface SelectedGroupData {
  /** Access mode: use existing or create new */
  mode: GroupAccessMode;

  /** Existing group data (when mode = 'existing') */
  existingGroup?: ExistingGroupData | undefined;

  /** New group data (when mode = 'new') */
  newGroup?: NewGroupData | undefined;
}

/**
 * Direct membership data (exception case)
 */
export interface DirectMembershipData {
  /** Organization ID */
  orgId: string;

  /** Organization name (for display) */
  orgName: string;

  /** Role IDs to assign */
  roleIds: string[];

  /** Role names (for display) */
  roleNames: string[];

  /** Access scope */
  scope: ScopeType;
}

/**
 * User form data structure
 * Mirrors CreateUserWithMemberships from the Go onboarding contract
 *
 * V1: AuthProvider removed - always internal auth.
 * Next version: auth provider will be determined by the customer's Organization.AuthConfig
 */
export interface UserFormData {
  /** User email address (required) */
  email: string;

  /** User password (required for create, optional for edit, min 8 chars) */
  password: string;

  /** Force password change on next login */
  changePasswordNextLogin: boolean;

  /** User first name (required, min 2, max 100) */
  firstName: string;

  /** User last name (required, min 2, max 100) */
  lastName: string;

  /** User phone number (optional, E.164 format) */
  phone: string;

  /** User job title (optional, max 120) */
  jobTitle: string;

  /** Whether user account is enabled */
  enabled: boolean;

  /** User avatar URL (optional) */
  avatar: string;

  /**
   * Access assignment type
   * - 'group': Add to existing group (recommended, inherits group permissions)
   * - 'direct': Direct membership (exception, harder to manage)
   * - 'both': Both group and direct access (edit mode only)
   */
  accessType: AccessType;

  /**
   * Selected group data (when accessType = 'group')
   * User will be added as member of this group
   * @deprecated Use selectedGroups instead for multi-group support
   */
  selectedGroup?: SelectedGroupData | undefined;

  /**
   * Selected groups data (supports multiple groups)
   * Array of groups the user should be member of
   * Backend uses DIFF logic: adds new, removes missing
   */
  selectedGroups?: SelectedGroupData[];

  /**
   * Direct membership data (when accessType = 'direct')
   * Creates a direct user → org membership
   * @deprecated Use directMemberships instead for multi-membership support
   */
  directMembership?: DirectMembershipData | undefined;

  /**
   * Direct memberships data (supports multiple)
   * Array of direct memberships with roles and scope
   */
  directMemberships?: DirectMembershipData[];
}

/**
 * User form state
 */
export interface UserFormState {
  /** Whether form is being saved */
  isSaving: boolean;

  /** Current step number */
  currentStep: number;
}

/**
 * Step validation errors state
 */
export interface StepValidationErrors {
  /** Step 1 has errors */
  step1: boolean;

  /** Step 2 has errors */
  step2: boolean;

  /** Step 3 has errors */
  step3: boolean;
}
