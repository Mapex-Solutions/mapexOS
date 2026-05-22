/**
 * UserProfilePage Interfaces
 */

/**
 * Profile section item for StepperVertical
 * Follows StepperVerticalItem format
 */
export interface ProfileSection {
  /** Display title for the section */
  title: string;
  /** Description text for the section */
  description: string;
  /** Material icon name for the section */
  icon: string;
}

/**
 * User profile data
 */
export interface UserProfileData {
  /** User first name */
  firstName: string;
  /** User last name */
  lastName: string;
  /** User email address (read-only) */
  email: string;
  /** User phone number */
  phone: string;
  /** User job title */
  jobTitle: string;
}

/**
 * Password change data
 */
export interface PasswordData {
  /** Current password for verification */
  current: string;
  /** New password to set */
  new: string;
  /** New password confirmation */
  confirm: string;
}

/**
 * Group info from API response
 */
export interface UserGroupInfo {
  /** Group ID */
  id: string;
  /** Group name */
  name: string;
  /** Optional description */
  description?: string | undefined;
}

/**
 * Membership info from API response
 */
export interface UserMembershipInfo {
  /** Organization ID */
  orgId: string;
  /** Organization name */
  orgName: string;
  /** Organization type */
  orgType: string;
  /** Scope: local or recursive */
  scope: 'local' | 'recursive';
  /** Role names array */
  roleNames: string[];
  /** Via: "direct" or "Group: {groupName}" */
  via: string;
}
