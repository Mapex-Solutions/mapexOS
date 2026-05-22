/**
 * AuditLogsPage Interfaces
 */

/**
 * Column visibility state for audit logs list page
 */
export interface AuditLogsPageColumnVisibility {
  /** Action column visibility */
  action: boolean;
  /** Resource column visibility */
  resource: boolean;
  /** Created column visibility */
  created: boolean;
}

/**
 * Audit log entry data structure
 */
export interface AuditLogEntry {
  /** Log entry ID */
  id: string | number;
  /** Actor (user) who performed the action */
  actor: string;
  /** Details about the action */
  details?: string;
  /** Action performed (Create, Edit, Delete) */
  action: string;
  /** Resource type */
  type: string;
  /** Resource name */
  resource: string;
  /** Status of the action */
  status: string;
  /** Created date of the action */
  created?: string;
}

/**
 * Filters for audit logs list page
 */
export interface AuditLogsPageFilters {
  /** Search by actor or resource name */
  search?: string;
  /** Filter by status (success/failure) */
  status?: string;
  /** Include child organizations */
  includeChildren?: boolean;
  /** Filter by action type */
  action?: string;
  /** Filter by resource type */
  resourceType?: string;
  /** Filter by date range */
  dateRange?: {
    from: string | null;
    to: string | null;
  };
}
