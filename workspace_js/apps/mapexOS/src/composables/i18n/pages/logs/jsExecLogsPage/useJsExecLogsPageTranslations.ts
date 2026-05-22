import { computed } from 'vue';
import { useTS } from '@utils/translation';

/**
 * Translation composable for JS Exec Logs Page
 * Provides all translated strings for the JS executor logs interface
 * @returns {Object} Translation object with page header, filters, columns, etc.
 */
export function useJsExecLogsPageTranslations() {
  const ts = useTS({ capitalize: true });
  const tsTitle = useTS({ titleCase: true });
  const tsRaw = useTS({ capitalize: false });

  return {
    pageHeader: {
      title: computed(() => tsTitle('pages.logs.jsExecLogsPage.pageHeader.title')),
      description: computed(() => tsRaw('pages.logs.jsExecLogsPage.pageHeader.description')),
    },
    listTitle: computed(() => tsTitle('pages.logs.jsExecLogsPage.listTitle')),
    itemLabel: computed(() => ts('pages.logs.jsExecLogsPage.itemLabel')),
    itemLabelPlural: computed(() => ts('pages.logs.jsExecLogsPage.itemLabelPlural')),
    filters: {
      /** Filter section label */
      label: computed(() => ts('pages.logs.jsExecLogsPage.filters.label')),
      /** Search placeholder for quick search input */
      searchPlaceholder: computed(() => ts('pages.logs.jsExecLogsPage.filters.searchPlaceholder')),
      /** Advanced filters tooltip */
      advancedFilters: computed(() => ts('pages.logs.jsExecLogsPage.filters.advancedFilters')),
      /** Pending filters tooltip */
      pendingFilters: computed(() => ts('pages.logs.jsExecLogsPage.filters.pendingFilters')),
      /** Clear all button label */
      clearAll: computed(() => ts('pages.logs.jsExecLogsPage.filters.clearAll')),
      /** All status option */
      allStatus: computed(() => ts('pages.logs.jsExecLogsPage.filters.allStatus')),
      /** Include children organizations filter */
      includeChildrenOrgs: computed(() => ts('pages.logs.jsExecLogsPage.filters.includeChildrenOrgs')),
      dateRange: computed(() => ts('pages.logs.jsExecLogsPage.filters.dateRange')),
      status: computed(() => ts('pages.logs.jsExecLogsPage.filters.status')),
      uuid: computed(() => ts('pages.logs.jsExecLogsPage.filters.uuid')),
      includeChildren: computed(() => ts('pages.logs.jsExecLogsPage.filters.includeChildren')),
      execTimeOperator: computed(() => ts('pages.logs.jsExecLogsPage.filters.execTimeOperator')),
      execTime: computed(() => ts('pages.logs.jsExecLogsPage.filters.execTime')),
      execTimeEnd: computed(() => ts('pages.logs.jsExecLogsPage.filters.execTimeEnd')),
      /** Filter options */
      options: {
        success: computed(() => ts('pages.logs.jsExecLogsPage.filters.options.success')),
        failed: computed(() => ts('pages.logs.jsExecLogsPage.filters.options.failed')),
        yes: computed(() => ts('pages.logs.jsExecLogsPage.filters.options.yes')),
        no: computed(() => ts('pages.logs.jsExecLogsPage.filters.options.no')),
      },
      includeChildrenOptions: {
        yes: computed(() => ts('pages.logs.jsExecLogsPage.filters.includeChildrenOptions.yes')),
        no: computed(() => ts('pages.logs.jsExecLogsPage.filters.includeChildrenOptions.no')),
      },
    },
    columns: {
      dataSource: computed(() => ts('pages.logs.jsExecLogsPage.columns.dataSource')),
      uuid: computed(() => ts('pages.logs.jsExecLogsPage.columns.uuid')),
      status: computed(() => ts('pages.logs.jsExecLogsPage.columns.status')),
      executionTime: computed(() => ts('pages.logs.jsExecLogsPage.columns.executionTime')),
      execTime: computed(() => ts('pages.logs.jsExecLogsPage.columns.execTime')),
      timestamp: computed(() => ts('pages.logs.jsExecLogsPage.columns.timestamp')),
      retention: computed(() => ts('pages.logs.jsExecLogsPage.columns.retention')),
    },
    statusOptions: {
      success: computed(() => ts('pages.logs.jsExecLogsPage.statusOptions.success')),
      failed: computed(() => ts('pages.logs.jsExecLogsPage.statusOptions.failed')),
    },
    statusBadge: {
      success: computed(() => tsRaw('pages.logs.jsExecLogsPage.statusBadge.success')),
      failed: computed(() => tsRaw('pages.logs.jsExecLogsPage.statusBadge.failed')),
    },
    execTimeOperators: {
      lessThan: computed(() => tsRaw('pages.logs.jsExecLogsPage.execTimeOperators.lessThan')),
      lessThanOrEqual: computed(() => tsRaw('pages.logs.jsExecLogsPage.execTimeOperators.lessThanOrEqual')),
      greaterThan: computed(() => tsRaw('pages.logs.jsExecLogsPage.execTimeOperators.greaterThan')),
      greaterThanOrEqual: computed(() => tsRaw('pages.logs.jsExecLogsPage.execTimeOperators.greaterThanOrEqual')),
      between: computed(() => ts('pages.logs.jsExecLogsPage.execTimeOperators.between')),
    },
    drawer: {
      title: computed(() => tsTitle('pages.logs.jsExecLogsPage.drawer.title')),
      subtitleSuccess: computed(() => ts('pages.logs.jsExecLogsPage.drawer.subtitleSuccess')),
      subtitleFailed: computed(() => ts('pages.logs.jsExecLogsPage.drawer.subtitleFailed')),
    },
    empty: {
      title: computed(() => tsTitle('pages.logs.jsExecLogsPage.empty.title')),
      description: computed(() => tsRaw('pages.logs.jsExecLogsPage.empty.description')),
    },
    pagination: {
      newer: computed(() => ts('pages.logs.jsExecLogsPage.pagination.newer')),
      older: computed(() => ts('pages.logs.jsExecLogsPage.pagination.older')),
    },
    messages: {
      loadFailed: computed(() => tsRaw('pages.logs.jsExecLogsPage.messages.loadFailed')),
    },
    defaults: {
      unknown: computed(() => tsRaw('pages.logs.jsExecLogsPage.defaults.unknown')),
      notAvailable: computed(() => tsRaw('pages.logs.jsExecLogsPage.defaults.notAvailable')),
    },
    placeholders: {
      execTime: computed(() => tsRaw('pages.logs.jsExecLogsPage.placeholders.execTime')),
      execTimeEnd: computed(() => tsRaw('pages.logs.jsExecLogsPage.placeholders.execTimeEnd')),
    },
  };
}
