/**
 * Route Groups Interfaces
 *
 * TypeScript interfaces mirroring the Go contract from:
 * workspace_go/packages/contracts/services/router/routegroups/dto.go
 */

/**
 * Data Lake Data
 * Routes events to a Data Lake for storage
 */
export interface LakeHouseData {
  lakeHouseId: string;
  metadata?: Record<string, any>;
}

/**
 * Notification Data
 * Routes events to a Notification channel
 */
export interface NotificationData {
  notificationId: string;
  metadata?: Record<string, any>;
}

/**
 * Save Event Data
 * Saves event without additional routing
 */
export interface SaveEventData {
  metadata?: Record<string, any>;
}

/**
 * Workflow Data
 * Defines how events are delivered to the Workflow Service.
 *
 * Modes:
 * - "newInstance": creates a new workflow instance from a definition
 * - "signal": delivers a signal to an existing running instance
 * - "signalOrStart": tries signal, falls back to creating a new instance
 */
export interface WorkflowData {
  /** Delivery mode */
  mode: 'newInstance' | 'signal' | 'signalOrStart';

  /** Mode-specific configuration (schema defined per mode) */
  data: Record<string, any>;
}

/**
 * Match Rule
 * Defines a single conditional rule for event routing
 * Supports multiple operators for flexible event filtering
 */
export interface MatchRule {
  /** JSON path to the field (e.g., "payload.temperature", "metadata.deviceType") */
  field: string;
  /** Comparison operator */
  operator: 'eq' | 'neq' | 'gt' | 'gte' | 'lt' | 'lte' | 'in' | 'nin';
  /** Value to compare against */
  value: any;
}

/**
 * Match Config
 * Defines the matching strategy for conditional routing
 * Allows AND/OR logic for multiple rules
 */
export interface MatchConfig {
  /** "all" = AND logic (all rules must pass), "any" = OR logic (at least one rule must pass) */
  policy: 'all' | 'any';
  /** Array of matching rules */
  rules: MatchRule[];
}

/**
 * Router
 * Defines a routing destination with optional conditional logic
 */
export interface Router {
  /** Router type */
  kind: 'lake_house' | 'notification' | 'save_event' | 'workflow';
  /** Optional conditional routing rules */
  match?: MatchConfig;
  /** Data Lake configuration (required when kind = 'lake_house') */
  lakeHouse?: LakeHouseData;
  /** Notification configuration (required when kind = 'notification') */
  notification?: NotificationData;
  /** Save Event configuration (optional when kind = 'save_event') */
  saveEvent?: SaveEventData;
  /** Workflow configuration (required when kind = 'workflow') */
  workflow?: WorkflowData;
}

/**
 * Route Group Create
 * Data structure for creating a new Route Group
 */
export interface RouteGroupCreate {
  version: string;
  name: string;
  description?: string;
  enabled: boolean;
  routers?: Router[];
  isTemplate?: boolean;
  orgId?: string;
  pathKey?: string;
}

/**
 * Route Group Update
 * Data structure for updating an existing Route Group
 * All fields are optional
 */
export interface RouteGroupUpdate {
  version?: string;
  name?: string;
  description?: string;
  enabled?: boolean;
  routers?: Router[];
  isTemplate?: boolean;
  orgId?: string;
}

/**
 * Route Group Response
 * Data structure returned by the API
 */
export interface RouteGroupResponse {
  id?: string;
  version?: string;
  name?: string;
  description?: string;
  enabled?: boolean;
  isTemplate?: boolean;
  orgId?: string;
  pathKey?: string;
  customerId?: string;
  routers?: Router[];
  created?: string;
  updated?: string;
}

/**
 * Route Group Query
 * Query parameters for listing route groups
 */
export interface RouteGroupQuery {
  // Pagination
  page?: number;
  perPage?: number;

  // Sorting
  sort?: string;

  // Projection (comma-separated fields)
  projection?: string;

  // Hierarchy
  includeChildren?: boolean;

  // Module-specific filters
  name?: string;
  enabled?: boolean;
  version?: string;
}

/**
 * Paginated Response
 * Generic paginated response structure
 */
export interface PaginatedResponse<T> {
  items: T[];
  pagination: {
    page: number;
    perPage: number;
    totalItems: number;
    totalPages: number;
  };
}

/**
 * Router Kind Options
 * Available router types for UI selection
 */
export const ROUTER_KIND_OPTIONS = [
  { label: 'Data Lake', value: 'lake_house', icon: 'storage', color: 'purple-6' },
  { label: 'Notification', value: 'notification', icon: 'notifications', color: 'orange-6' },
  { label: 'Save Event', value: 'save_event', icon: 'save', color: 'green-6' },
] as const;

/**
 * Match Operator Options
 * Available comparison operators for match rules
 */
export const MATCH_OPERATOR_OPTIONS = [
  { label: 'Equal (=)', value: 'eq', description: 'Field equals value' },
  { label: 'Not Equal (!=)', value: 'neq', description: 'Field does not equal value' },
  { label: 'Greater Than (>)', value: 'gt', description: 'Field is greater than value' },
  { label: 'Greater or Equal (>=)', value: 'gte', description: 'Field is greater than or equal to value' },
  { label: 'Less Than (<)', value: 'lt', description: 'Field is less than value' },
  { label: 'Less or Equal (<=)', value: 'lte', description: 'Field is less than or equal to value' },
  { label: 'In Array', value: 'in', description: 'Field value is in the array' },
  { label: 'Not In Array', value: 'nin', description: 'Field value is not in the array' },
] as const;

/**
 * Match Policy Options
 * Available matching policies
 */
export const MATCH_POLICY_OPTIONS = [
  { label: 'All (AND)', value: 'all', description: 'All rules must match' },
  { label: 'Any (OR)', value: 'any', description: 'At least one rule must match' },
] as const;
