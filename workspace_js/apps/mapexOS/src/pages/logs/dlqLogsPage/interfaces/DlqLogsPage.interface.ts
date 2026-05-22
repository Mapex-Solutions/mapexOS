/**
 * Cursor state for DLQ logs pagination
 */
export interface DlqLogsPageCursor {
  current?: string;
  next?: string;
  prev?: string;
  hasNext: boolean;
  hasPrevious: boolean;
}

/**
 * Filter state for DLQ logs page
 */
export interface DlqLogsPageFilters {
  serviceName?: string;
  serviceType?: string;
  eventType?: string;
  startTime?: string;
  endTime?: string;
  includeChildren?: boolean;
}

/**
 * Sidebar view type for DLQ page
 */
export type DlqSidebarView = 'all' | 'service_type';

/**
 * Service type group with count for sidebar
 */
export interface DlqServiceTypeGroup {
  serviceType: string;
  icon: string;
  count: number;
}
