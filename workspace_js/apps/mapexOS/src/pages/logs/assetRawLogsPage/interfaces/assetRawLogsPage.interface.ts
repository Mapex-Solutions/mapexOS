/**
 * AssetRawLogsPage Interfaces
 */

/**
 * Column visibility state for raw logs list page
 */
export interface AssetRawLogsPageColumnVisibility {
  /** Data Source name column visibility */
  name: boolean;
  /** Thread ID (UUID) column visibility */
  threadId: boolean;
  /** Source column visibility */
  source: boolean;
  /** Created column visibility */
  created: boolean;
}

/**
 * Filters state for raw logs page
 * Maps to EventsRawQuery from backend contract
 */
export interface AssetRawLogsPageFilters {
  /** Filter by thread ID (data source identifier) */
  threadId?: string;
  /** Filter by source (http_gateway, mqtt_gateway, etc.) */
  source?: string;
  /** Filter by auth success status */
  success?: boolean;
  /** Filter events after this timestamp (ISO 8601) */
  startTime?: string;
  /** Filter events before this timestamp (ISO 8601) */
  endTime?: string;
  /** Include child organizations hierarchically */
  includeChildren?: boolean;
}

/**
 * Cursor state for raw logs page
 * Used for cursor-based pagination
 */
export interface AssetRawLogsPageCursor {
  /** Current cursor timestamp (for fetching next/prev page) */
  current: string | undefined;
  /** Next cursor timestamp (for fetching older items) */
  next: string | undefined;
  /** Previous cursor timestamp (for fetching newer items) */
  prev: string | undefined;
  /** Whether there are more items after current page */
  hasNext: boolean;
  /** Whether there are more items before current page */
  hasPrevious: boolean;
}
