import { computed } from 'vue';
import { useTS } from '@utils/translation';

/**
 * Translation composable for Trigger Logs Page
 * Provides all translated strings for the trigger logs interface
 * @returns {Object} Translation object with page header, filters, columns, etc.
 */
export function useTriggerLogsPageTranslations() {
  const ts = useTS({ capitalize: true });
  const tsTitle = useTS({ titleCase: true });
  const tsRaw = useTS({ capitalize: false });

  return {
    pageHeader: {
      title: computed(() => tsTitle('pages.logs.triggerLogsPage.pageHeader.title')),
      description: computed(() => tsRaw('pages.logs.triggerLogsPage.pageHeader.description')),
    },
    listTitle: computed(() => tsTitle('pages.logs.triggerLogsPage.listTitle')),
    itemLabel: computed(() => ts('pages.logs.triggerLogsPage.itemLabel')),
    itemLabelPlural: computed(() => ts('pages.logs.triggerLogsPage.itemLabelPlural')),

    /**
     * Filter translations for Enterprise Filter Pattern
     */
    filters: {
      label: computed(() => ts('pages.logs.triggerLogsPage.filters.label')),
      searchPlaceholder: computed(() => ts('pages.logs.triggerLogsPage.filters.searchPlaceholder')),
      advancedFilters: computed(() => ts('pages.logs.triggerLogsPage.filters.advancedFilters')),
      pendingFilters: computed(() => ts('pages.logs.triggerLogsPage.filters.pendingFilters')),
      clearAll: computed(() => ts('pages.logs.triggerLogsPage.filters.clearAll')),
      allStatus: computed(() => ts('pages.logs.triggerLogsPage.filters.allStatus')),
      dateRange: computed(() => ts('pages.logs.triggerLogsPage.filters.dateRange')),
      status: computed(() => ts('pages.logs.triggerLogsPage.filters.status')),
      source: computed(() => ts('pages.logs.triggerLogsPage.filters.source')),
      triggerType: computed(() => ts('pages.logs.triggerLogsPage.filters.triggerType')),
      category: computed(() => ts('pages.logs.triggerLogsPage.filters.category')),
      triggerId: computed(() => ts('pages.logs.triggerLogsPage.filters.triggerId')),
      triggerIdPlaceholder: computed(() => ts('pages.logs.triggerLogsPage.filters.triggerIdPlaceholder')),
      includeChildren: computed(() => ts('pages.logs.triggerLogsPage.filters.includeChildren')),
      includeChildrenOrgs: computed(() => ts('pages.logs.triggerLogsPage.filters.includeChildrenOrgs')),
      options: {
        yes: computed(() => ts('pages.logs.triggerLogsPage.filters.options.yes')),
        no: computed(() => ts('pages.logs.triggerLogsPage.filters.options.no')),
        success: computed(() => ts('pages.logs.triggerLogsPage.statusOptions.success')),
        failure: computed(() => ts('pages.logs.triggerLogsPage.statusOptions.failure')),
      },
    },

    statusOptions: {
      success: computed(() => ts('pages.logs.triggerLogsPage.statusOptions.success')),
      failure: computed(() => ts('pages.logs.triggerLogsPage.statusOptions.failure')),
    },
    sourceOptions: {
      router: computed(() => ts('pages.logs.triggerLogsPage.sourceOptions.router')),
    },
    triggerTypeOptions: {
      http: computed(() => ts('pages.logs.triggerLogsPage.triggerTypeOptions.http')),
      mqtt: computed(() => ts('pages.logs.triggerLogsPage.triggerTypeOptions.mqtt')),
      rabbitmq: computed(() => ts('pages.logs.triggerLogsPage.triggerTypeOptions.rabbitmq')),
      nats: computed(() => ts('pages.logs.triggerLogsPage.triggerTypeOptions.nats')),
      websocket: computed(() => ts('pages.logs.triggerLogsPage.triggerTypeOptions.websocket')),
      email: computed(() => ts('pages.logs.triggerLogsPage.triggerTypeOptions.email')),
      teams: computed(() => ts('pages.logs.triggerLogsPage.triggerTypeOptions.teams')),
      slack: computed(() => ts('pages.logs.triggerLogsPage.triggerTypeOptions.slack')),
    },
    categoryOptions: {
      technical: computed(() => ts('pages.logs.triggerLogsPage.categoryOptions.technical')),
      communication: computed(() => ts('pages.logs.triggerLogsPage.categoryOptions.communication')),
    },
    columns: {
      triggerName: computed(() => ts('pages.logs.triggerLogsPage.columns.triggerName')),
      triggerType: computed(() => ts('pages.logs.triggerLogsPage.columns.triggerType')),
      category: computed(() => ts('pages.logs.triggerLogsPage.columns.category')),
      source: computed(() => ts('pages.logs.triggerLogsPage.columns.source')),
      status: computed(() => ts('pages.logs.triggerLogsPage.columns.status')),
      duration: computed(() => ts('pages.logs.triggerLogsPage.columns.duration')),
      timestamp: computed(() => ts('pages.logs.triggerLogsPage.columns.timestamp')),
    },
    statusBadge: {
      success: computed(() => tsRaw('pages.logs.triggerLogsPage.statusBadge.success')),
      failure: computed(() => tsRaw('pages.logs.triggerLogsPage.statusBadge.failure')),
    },
    drawer: {
      title: computed(() => tsTitle('pages.logs.triggerLogsPage.drawer.title')),
    },
    empty: {
      title: computed(() => tsTitle('pages.logs.triggerLogsPage.empty.title')),
      description: computed(() => tsRaw('pages.logs.triggerLogsPage.empty.description')),
    },
    pagination: {
      newer: computed(() => ts('pages.logs.triggerLogsPage.pagination.newer')),
      older: computed(() => ts('pages.logs.triggerLogsPage.pagination.older')),
    },
    messages: {
      loadFailed: computed(() => tsRaw('pages.logs.triggerLogsPage.messages.loadFailed')),
    },
    defaults: {
      unknown: computed(() => tsRaw('pages.logs.triggerLogsPage.defaults.unknown')),
      notAvailable: computed(() => tsRaw('pages.logs.triggerLogsPage.defaults.notAvailable')),
    },
  };
}
