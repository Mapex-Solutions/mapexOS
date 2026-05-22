/**
 * AssetConnectivityLogsPage Interfaces
 */

/**
 * Column visibility state for connectivity logs list page
 */
export interface AssetConnectivityLogsPageColumnVisibility {
  /** Asset name column visibility */
  asset: boolean;
  /** Asset UUID column visibility */
  assetUUID: boolean;
  /** Event type (offline/online) column visibility */
  eventType: boolean;
  /** Last seen at column visibility */
  lastSeenAt: boolean;
  /** Missed heartbeats count column visibility */
  missCount: boolean;
  /** Threshold minutes column visibility */
  thresholdMinutes: boolean;
  /** Created column visibility */
  created: boolean;
}

/**
 * Filters state for connectivity logs page.
 * Maps to AssetConnectivityQuery from backend contract.
 */
export interface AssetConnectivityLogsPageFilters {
  /** Filter events after this timestamp (ISO 8601) */
  from?: string;
  /** Filter events before this timestamp (ISO 8601) */
  to?: string;
  /** Filter by event type: 'offline' or 'online' */
  eventType?: 'offline' | 'online';
  /** Filter by asset UUID (uses backend-side thread/uuid match if exposed) */
  assetUUID?: string;
  /** Include child organizations hierarchically */
  includeChildren?: boolean;
}

/**
 * Cursor state for connectivity logs page.
 * Used for cursor-based pagination.
 */
export interface AssetConnectivityLogsPageCursor {
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
