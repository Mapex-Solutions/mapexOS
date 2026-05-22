/**
 * JsExecLogsPage Interfaces
 */

/**
 * Column visibility state for JS executor logs list page
 */
export interface JsExecLogsPageColumnVisibility {
  /** Thread ID (asset UUID) column visibility */
  threadId: boolean;
  /** Name (data source name) column visibility */
  name: boolean;
  /** Success status column visibility */
  success: boolean;
  /** Execution time column visibility */
  executionTime: boolean;
  /** Created column visibility */
  created: boolean;
}

/**
 * Filters state for JS executor logs page
 * Maps to EventsJsExecQuery from backend contract
 */
export interface JsExecLogsPageFilters {
  /** Filter by thread ID (asset UUID) */
  threadId?: string;
  /** Filter by execution success status */
  success?: boolean;
  /** Filter events after this timestamp (ISO 8601) */
  startTime?: string;
  /** Filter events before this timestamp (ISO 8601) */
  endTime?: string;
  /** Include child organizations hierarchically */
  includeChildren?: boolean;
  /** Operator for execution time filter (lt, lte, gt, gte, between) */
  execTimeOp?: string;
  /** Execution time value in milliseconds */
  execTimeValue?: number;
  /** End value for "between" operator (milliseconds) */
  execTimeValueEnd?: number;
}

/**
 * Cursor state for JS executor logs page
 * Used for cursor-based pagination
 */
export interface JsExecLogsPageCursor {
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
