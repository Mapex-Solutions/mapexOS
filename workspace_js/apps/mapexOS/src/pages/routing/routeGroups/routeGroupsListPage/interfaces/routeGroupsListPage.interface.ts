/**
 * RouteGroupsListPage Interfaces
 */

/**
 * Filter state for route groups list page
 */
export interface RouteGroupsListPageFilters {
  /** Route group name filter */
  name: string | undefined;
  /** Enabled status filter */
  enabled: boolean | undefined;
  /** Template filter */
  isTemplate: boolean | undefined;
  /** Include child organizations filter */
  includeChildren: boolean | undefined;
}

/**
 * Column visibility state for route groups list page
 */
export interface RouteGroupsListPageColumnVisibility {
  /** Organization column visibility */
  organization: boolean;
  /** Routers column visibility */
  routers: boolean;
  /** Template source column visibility */
  isTemplate: boolean;
}

/**
 * Enriched route group with additional computed fields
 */
export interface EnrichedRouteGroup {
  /** Route group ID */
  id?: string;
  /** Route group name */
  name?: string;
  /** Route group description */
  description?: string;
  /** Version */
  version?: string;
  /** Enabled status */
  enabled?: boolean;
  /** Routers array */
  routers?: any[];
  /** Is template flag */
  isTemplate?: boolean;
  /** Organization ID */
  orgId?: string;
  /** Organization name (enriched) */
  organizationName?: string;
  /** Routers count (enriched) */
  routersCount?: number;
}
