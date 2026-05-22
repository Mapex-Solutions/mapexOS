/**
 * Cursor state for trigger logs pagination
 */
export interface TriggerLogsPageCursor {
  current?: string;
  next?: string;
  prev?: string;
  hasNext: boolean;
  hasPrevious: boolean;
}

/**
 * Filter state for trigger logs page
 */
export interface TriggerLogsPageFilters {
  triggerId?: string;
  triggerType?: string;
  category?: string;
  source?: string;
  success?: boolean;
  startTime?: string;
  endTime?: string;
  includeChildren?: boolean;
}
