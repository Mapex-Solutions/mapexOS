/**
 * Tab configuration for UserDetailPage
 */
export interface UserDetailTab {
  /** Unique key for the tab */
  name: string;
  /** Display label */
  label: string;
  /** Icon name */
  icon: string;
  /** Optional badge count */
  badge?: number;
  /** Badge color */
  badgeColor?: string;
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
  description?: string;
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

/**
 * User detail data for display (matches API response)
 */
export interface UserDetailData {
  id: string;
  email: string;
  firstName?: string;
  lastName?: string;
  phone?: string;
  jobTitle?: string;
  avatar?: string;
  enabled: boolean;
  changePasswordNextLogin?: boolean;
  authProvider?: {
    type: 'internal' | 'google' | 'github' | 'microsoft' | 'keycloak';
    externalId?: string;
    metadata?: Record<string, any>;
  };
  orgId?: string;
  organizationName?: string;
  created?: string;
  updated?: string;
  /** Number of groups user belongs to */
  groupsCount?: number;
  /** Groups user belongs to */
  groups?: UserGroupInfo[];
  /** User memberships across organizations */
  memberships?: UserMembershipInfo[];
}

/**
 * Props for TabProfile component
 */
export interface TabProfileProps {
  /** User data to display */
  user: UserDetailData | null;
  /** Loading state */
  loading: boolean;
}

/**
 * Props for TabAccess component
 */
export interface TabAccessProps {
  /** User data to display access for */
  user: UserDetailData | null;
  /** Loading state */
  loading: boolean;
}

/**
 * Props for TabGroups component
 */
export interface TabGroupsProps {
  /** User data to display groups for */
  user: UserDetailData | null;
  /** Loading state */
  loading: boolean;
}

/**
 * Group membership item for display (legacy - kept for compatibility)
 */
export interface GroupMembershipItem {
  id: string;
  name: string;
  description?: string;
  isSystem: boolean;
  isTemplate: boolean;
  enabled: boolean;
}

/**
 * Role access item for display
 */
export interface RoleAccessItem {
  id: string;
  name: string;
  description?: string;
  permissions: string[];
  scope: string;
  isSystem: boolean;
  organizationName?: string;
}
