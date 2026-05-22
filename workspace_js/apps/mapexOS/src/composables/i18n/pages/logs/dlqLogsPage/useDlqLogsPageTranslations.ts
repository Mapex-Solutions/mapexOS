import { computed } from 'vue';
import { useTS } from '@utils/translation';

/**
 * Translation composable for DLQ Logs Page
 * Provides all translated strings for the dead letter queue logs interface
 * @returns {Object} Translation object with page header, filters, detail, sidebar, etc.
 */
export function useDlqLogsPageTranslations() {
  const ts = useTS({ capitalize: true });
  const tsTitle = useTS({ titleCase: true });
  const tsRaw = useTS({ capitalize: false });

  return {
    pageHeader: {
      title: computed(() => tsTitle('pages.logs.dlqLogsPage.pageHeader.title')),
      description: computed(() => tsRaw('pages.logs.dlqLogsPage.pageHeader.description')),
    },
    listTitle: computed(() => tsTitle('pages.logs.dlqLogsPage.listTitle')),
    itemLabel: computed(() => ts('pages.logs.dlqLogsPage.itemLabel')),
    itemLabelPlural: computed(() => ts('pages.logs.dlqLogsPage.itemLabelPlural')),

    sidebar: {
      searchPlaceholder: computed(() => ts('pages.logs.dlqLogsPage.sidebar.searchPlaceholder')),
      searchErrorPlaceholder: computed(() => ts('pages.logs.dlqLogsPage.sidebar.searchErrorPlaceholder')),
      allFailures: computed(() => ts('pages.logs.dlqLogsPage.sidebar.allFailures')),
      byService: computed(() => ts('pages.logs.dlqLogsPage.sidebar.byService')),
      noServiceTypes: computed(() => tsRaw('pages.logs.dlqLogsPage.sidebar.noServiceTypes')),
      infoBanner: computed(() => tsRaw('pages.logs.dlqLogsPage.sidebar.infoBanner')),
    },

    filters: {
      label: computed(() => ts('pages.logs.dlqLogsPage.filters.label')),
      searchPlaceholder: computed(() => ts('pages.logs.dlqLogsPage.filters.searchPlaceholder')),
      advancedFilters: computed(() => ts('pages.logs.dlqLogsPage.filters.advancedFilters')),
      pendingFilters: computed(() => ts('pages.logs.dlqLogsPage.filters.pendingFilters')),
      clearAll: computed(() => ts('pages.logs.dlqLogsPage.filters.clearAll')),
      allStatus: computed(() => ts('pages.logs.dlqLogsPage.filters.allStatus')),
      serviceName: computed(() => ts('pages.logs.dlqLogsPage.filters.serviceName')),
      serviceType: computed(() => ts('pages.logs.dlqLogsPage.filters.serviceType')),
      eventType: computed(() => ts('pages.logs.dlqLogsPage.filters.eventType')),
      includeChildren: computed(() => ts('pages.logs.dlqLogsPage.filters.includeChildren')),
      includeChildrenOrgs: computed(() => ts('pages.logs.dlqLogsPage.filters.includeChildrenOrgs')),
      options: {
        yes: computed(() => ts('pages.logs.dlqLogsPage.filters.options.yes')),
        no: computed(() => ts('pages.logs.dlqLogsPage.filters.options.no')),
      },
    },

    columns: {
      serviceName: computed(() => ts('pages.logs.dlqLogsPage.columns.serviceName')),
      serviceType: computed(() => ts('pages.logs.dlqLogsPage.columns.serviceType')),
      eventType: computed(() => ts('pages.logs.dlqLogsPage.columns.eventType')),
      errorCount: computed(() => ts('pages.logs.dlqLogsPage.columns.errorCount')),
      consumerName: computed(() => ts('pages.logs.dlqLogsPage.columns.consumerName')),
      totalDeliveries: computed(() => ts('pages.logs.dlqLogsPage.columns.totalDeliveries')),
      timestamp: computed(() => ts('pages.logs.dlqLogsPage.columns.timestamp')),
    },

    detail: {
      backToList: computed(() => ts('pages.logs.dlqLogsPage.detail.backToList')),
      copyContent: computed(() => ts('pages.logs.dlqLogsPage.detail.copyContent')),
      unknownError: computed(() => tsRaw('pages.logs.dlqLogsPage.detail.unknownError')),
      errors: computed(() => ts('pages.logs.dlqLogsPage.detail.errors')),
      deliveries: computed(() => ts('pages.logs.dlqLogsPage.detail.deliveries')),
      duration: computed(() => ts('pages.logs.dlqLogsPage.detail.duration')),
      retention: computed(() => ts('pages.logs.dlqLogsPage.detail.retention')),
      searchHint: computed(() => tsRaw('pages.logs.dlqLogsPage.detail.searchHint')),
      tabs: {
        payload: computed(() => ts('pages.logs.dlqLogsPage.detail.tabs.payload')),
        headers: computed(() => ts('pages.logs.dlqLogsPage.detail.tabs.headers')),
        error: computed(() => ts('pages.logs.dlqLogsPage.detail.tabs.error')),
      },
      fallback: {
        empty: computed(() => tsRaw('pages.logs.dlqLogsPage.detail.fallback.empty')),
        noHeaders: computed(() => tsRaw('pages.logs.dlqLogsPage.detail.fallback.noHeaders')),
        noErrorMessage: computed(() => tsRaw('pages.logs.dlqLogsPage.detail.fallback.noErrorMessage')),
      },
    },

    list: {
      unknownError: computed(() => tsRaw('pages.logs.dlqLogsPage.list.unknownError')),
    },

    datePresets: {
      today: computed(() => ts('pages.logs.dlqLogsPage.datePresets.today')),
      all: computed(() => ts('pages.logs.dlqLogsPage.datePresets.all')),
    },

    time: {
      justNow: computed(() => tsRaw('pages.logs.dlqLogsPage.time.justNow')),
      minutesAgo: computed(() => tsRaw('pages.logs.dlqLogsPage.time.minutesAgo')),
      hoursAgo: computed(() => tsRaw('pages.logs.dlqLogsPage.time.hoursAgo')),
      daysAgo: computed(() => tsRaw('pages.logs.dlqLogsPage.time.daysAgo')),
    },

    empty: {
      title: computed(() => tsTitle('pages.logs.dlqLogsPage.empty.title')),
      description: computed(() => tsRaw('pages.logs.dlqLogsPage.empty.description')),
    },

    pagination: {
      newer: computed(() => ts('pages.logs.dlqLogsPage.pagination.newer')),
      older: computed(() => ts('pages.logs.dlqLogsPage.pagination.older')),
    },

    actions: {
      refresh: computed(() => ts('pages.logs.dlqLogsPage.actions.refresh')),
    },

    messages: {
      loadFailed: computed(() => tsRaw('pages.logs.dlqLogsPage.messages.loadFailed')),
      copied: computed(() => tsRaw('pages.logs.dlqLogsPage.messages.copied')),
      copyFailed: computed(() => tsRaw('pages.logs.dlqLogsPage.messages.copyFailed')),
    },

    defaults: {
      unknown: computed(() => tsRaw('pages.logs.dlqLogsPage.defaults.unknown')),
      notAvailable: computed(() => tsRaw('pages.logs.dlqLogsPage.defaults.notAvailable')),
    },
  };
}
