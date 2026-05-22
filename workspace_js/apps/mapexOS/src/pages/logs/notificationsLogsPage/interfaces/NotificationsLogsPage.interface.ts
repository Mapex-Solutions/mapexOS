/**
 * Filter state for notifications logs page
 */
export interface NotificationsLogsPageFilters {
  notificationName?: string | undefined;
  status?: boolean | undefined;
  notificationType?: string | undefined;
  includeChildren?: boolean | undefined;
}
