/**
 * AccessAuditPage Interfaces
 */

/**
 * Filter state for access audit list page
 */
export interface AccessAuditPageFilters {
  /** Filter by assignee type (user or group) */
  assigneeType: string | undefined;
  /** Filter by specific assignee ID */
  assigneeId: string | undefined;
  /** Filter by role ID */
  roleId: string | undefined;
  /** Filter by scope (local or recursive) */
  scope: string | undefined;
  /** Filter by enabled status */
  enabled: boolean | undefined;
  /** Include children organizations */
  includeChildren: boolean | undefined;
}

/**
 * Column visibility state for access audit list page
 */
export interface AccessAuditPageColumnVisibility {
  /** Show organization column */
  organization: boolean;
  /** Show roles column */
  roles: boolean;
  /** Show scope column */
  scope: boolean;
  /** Show enabled status column */
  enabled: boolean;
}

/**
 * Membership response with enriched data for display
 */
export interface EnrichedMembership {
  /** Membership ID */
  id: string;
  /** Assignee type (user or group) */
  assigneeType: string;
  /** Assignee ID */
  assigneeId: string;
  /** Assignee display name */
  assigneeName: string;
  /** Organization ID */
  orgId: string;
  /** Organization name */
  orgName: string;
  /** Organization path key */
  orgPathKey: string;
  /** Role IDs */
  roleIds: string[];
  /** Role names for display */
  roleNames: string[];
  /** Access scope */
  scope: string;
  /** Enabled status */
  enabled: boolean;
  /** Created timestamp */
  created: string;
}
