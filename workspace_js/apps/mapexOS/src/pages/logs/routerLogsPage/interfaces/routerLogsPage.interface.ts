/**
 * Router Logs Page Interfaces
 */

/**
 * Filters for router logs page
 */
export interface RouterLogsPageFilters {
  /** Filter by thread ID (asset ID) */
  threadId?: string;

  /** Filter by asset ID */
  assetId?: string;

  /** Filter by router/RouteGroup ID */
  routerId?: string;

  /** Filter by success status */
  success?: boolean;

  /** Filter events after this timestamp */
  startTime?: string;

  /** Filter events before this timestamp */
  endTime?: string;

  /** Include children organizations */
  includeChildren?: boolean;
}

/**
 * Column visibility state for router logs page
 */
export interface RouterLogsPageColumnVisibility {
  /** Show thread ID column */
  threadId: boolean;

  /** Show name column */
  name: boolean;

  /** Show success column */
  success: boolean;

  /** Show routers count column */
  routersCount: boolean;

  /** Show published count column */
  publishedCount: boolean;

  /** Show created column */
  created: boolean;
}

/**
 * Cursor pagination state for router logs page
 */
export interface RouterLogsPageCursor {
  /** Current cursor position */
  current?: string;

  /** Cursor for next page (older items) */
  next?: string;

  /** Cursor for previous page (newer items) */
  prev?: string;

  /** Whether there are more items after current page */
  hasNext: boolean;

  /** Whether there are more items before current page */
  hasPrevious: boolean;
}
