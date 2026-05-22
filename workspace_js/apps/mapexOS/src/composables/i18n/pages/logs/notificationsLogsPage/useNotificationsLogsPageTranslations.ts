import { computed } from 'vue';
import { useTS } from '@utils/translation';

/**
 * Translation composable for Notifications Logs Page
 * Provides all translated strings for the notifications logs interface
 * @returns {Object} Translation object with page header, filters, columns, etc.
 */
export function useNotificationsLogsPageTranslations() {
  const ts = useTS({ capitalize: true });
  const tsTitle = useTS({ titleCase: true });
  const tsRaw = useTS({ capitalize: false });

  return {
    pageHeader: {
      title: computed(() => tsTitle('pages.logs.notificationsLogsPage.pageHeader.title')),
      description: computed(() => tsRaw('pages.logs.notificationsLogsPage.pageHeader.description')),
    },
    listTitle: computed(() => tsTitle('pages.logs.notificationsLogsPage.listTitle')),
    itemLabel: computed(() => ts('pages.logs.notificationsLogsPage.itemLabel')),
    itemLabelPlural: computed(() => ts('pages.logs.notificationsLogsPage.itemLabelPlural')),

    /**
     * Filter translations for Enterprise Filter Pattern
     */
    filters: {
      label: computed(() => ts('pages.logs.notificationsLogsPage.filters.label')),
      searchPlaceholder: computed(() => ts('pages.logs.notificationsLogsPage.filters.searchPlaceholder')),
      advancedFilters: computed(() => ts('pages.logs.notificationsLogsPage.filters.advancedFilters')),
      pendingFilters: computed(() => ts('pages.logs.notificationsLogsPage.filters.pendingFilters')),
      clearAll: computed(() => ts('pages.logs.notificationsLogsPage.filters.clearAll')),
      allStatus: computed(() => ts('pages.logs.notificationsLogsPage.filters.allStatus')),
      name: computed(() => ts('pages.logs.notificationsLogsPage.filters.name')),
      status: computed(() => ts('pages.logs.notificationsLogsPage.filters.status')),
      notificationType: computed(() => ts('pages.logs.notificationsLogsPage.filters.notificationType')),
      includeChildren: computed(() => ts('pages.logs.notificationsLogsPage.filters.includeChildren')),
      includeChildrenOrgs: computed(() => ts('pages.logs.notificationsLogsPage.filters.includeChildrenOrgs')),
      options: {
        success: computed(() => ts('pages.logs.notificationsLogsPage.filters.options.success')),
        failure: computed(() => ts('pages.logs.notificationsLogsPage.filters.options.failure')),
        yes: computed(() => ts('pages.logs.notificationsLogsPage.filters.options.yes')),
        no: computed(() => ts('pages.logs.notificationsLogsPage.filters.options.no')),
      },
    },

    notificationTypeOptions: {
      slack: computed(() => ts('pages.logs.notificationsLogsPage.notificationTypeOptions.slack')),
      teams: computed(() => ts('pages.logs.notificationsLogsPage.notificationTypeOptions.teams')),
      email: computed(() => ts('pages.logs.notificationsLogsPage.notificationTypeOptions.email')),
      push: computed(() => ts('pages.logs.notificationsLogsPage.notificationTypeOptions.push')),
      telegram: computed(() => ts('pages.logs.notificationsLogsPage.notificationTypeOptions.telegram')),
      webhook: computed(() => ts('pages.logs.notificationsLogsPage.notificationTypeOptions.webhook')),
    },

    columns: {
      notificationName: computed(() => ts('pages.logs.notificationsLogsPage.columns.notificationName')),
      notificationType: computed(() => ts('pages.logs.notificationsLogsPage.columns.notificationType')),
      tenantId: computed(() => ts('pages.logs.notificationsLogsPage.columns.tenantId')),
      status: computed(() => ts('pages.logs.notificationsLogsPage.columns.status')),
      created: computed(() => ts('pages.logs.notificationsLogsPage.columns.created')),
    },

    menuColumns: {
      notificationType: computed(() => ts('pages.logs.notificationsLogsPage.menuColumns.notificationType')),
      tenantId: computed(() => ts('pages.logs.notificationsLogsPage.menuColumns.tenantId')),
      created: computed(() => ts('pages.logs.notificationsLogsPage.menuColumns.created')),
    },

    statusBadge: {
      success: computed(() => tsRaw('pages.logs.notificationsLogsPage.statusBadge.success')),
      failure: computed(() => tsRaw('pages.logs.notificationsLogsPage.statusBadge.failure')),
    },

    drawer: {
      title: computed(() => tsTitle('pages.logs.notificationsLogsPage.drawer.title')),
    },

    empty: {
      title: computed(() => tsTitle('pages.logs.notificationsLogsPage.empty.title')),
      description: computed(() => tsRaw('pages.logs.notificationsLogsPage.empty.description')),
    },

    pagination: {
      newer: computed(() => ts('pages.logs.notificationsLogsPage.pagination.newer')),
      older: computed(() => ts('pages.logs.notificationsLogsPage.pagination.older')),
    },

    messages: {
      loadFailed: computed(() => tsRaw('pages.logs.notificationsLogsPage.messages.loadFailed')),
    },

    defaults: {
      unknown: computed(() => tsRaw('pages.logs.notificationsLogsPage.defaults.unknown')),
      notAvailable: computed(() => tsRaw('pages.logs.notificationsLogsPage.defaults.notAvailable')),
    },
  };
}
